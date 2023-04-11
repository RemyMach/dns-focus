package resolver

import (
	"bufio"
	"crypto/rand"
	"dns-focus/src/common"
	"dns-focus/src/config"
	"fmt"
	"log"
	"math/big"
	"net"

	"dns-focus/src/utils"

	"golang.org/x/net/dns/dnsmessage"
)

var ROOT_SERVERS = [8]string{"198.41.0.4","199.9.14.201","192.33.4.12","199.7.91.13","192.203.230.10","192.5.5.241","192.112.36.4","198.97.190.53"}

func HandlePacket(pc net.PacketConn, addr net.Addr, buf []byte, dnsConfig *config.DnsConfig, serverMode common.ServerMode) {

	ipBlocked, err := handleBlockDomains(pc, addr, buf, dnsConfig)
	if err != nil {
		log.Println(err.Error())
		return		
	}
	if ipBlocked {
		return
	}

	if serverMode == common.Proxy {
		handleDNSRequestToProxy(pc, buf, addr)
	} else if serverMode == common.Server {
		if err := handlePacket(pc, addr, buf); err != nil {
			fmt.Printf("handlePacket error [%s]: %s\n", addr.String(), err)
		}
	}
}

func handlePacket(pc net.PacketConn, addr net.Addr, buf []byte) error {
	p := dnsmessage.Parser{}
	header, err := p.Start(buf)
	if err != nil {
		return err
	}
	questions, err := p.AllQuestions()
	if err != nil {
		return err
	}
	response, err := dnsQuery(getRootServers(), questions[0], 3)
	if err != nil {
		return err
	}
	// set response object here
	response.Header.ID = header.ID
	response.Questions = questions
	responseBuffer, err := response.Pack()
	if err != nil {
		return err
	}

	_, err = pc.WriteTo(responseBuffer, addr)
	if err != nil {
		return err
	}

	return nil
}

func dnsQuery(servers []net.IP, question dnsmessage.Question, numberRetry int) (*dnsmessage.Message, error) {
	
	for i := 0; i < numberRetry; i++ {
		dnsAnswer, header, err := outgoingDnsQuery(servers, question)
		if err != nil {
			return nil, err
		}
		parsedAnswers, err := dnsAnswer.AllAnswers()
		if err != nil {
			return nil, err
		}

		if header.Authoritative {
			return &dnsmessage.Message{
				Header:  dnsmessage.Header{
					Response: true, 
					RecursionAvailable: true,
				},
				Answers: parsedAnswers,
			}, nil
		}
		authorities, err := dnsAnswer.AllAuthorities()
		if err != nil {
			return nil, err
		}

		if len(authorities) == 0 {
			return &dnsmessage.Message{
				Header: dnsmessage.Header{RCode: dnsmessage.RCodeNameError},
			}, nil
		}

		nameservers := make([]string, len(authorities))
		for k, authority := range authorities {
			if authority.Header.Type == dnsmessage.TypeNS {
				nameservers[k] = authority.Body.(*dnsmessage.NSResource).NS.String()
			}
		}

		additionals, err := dnsAnswer.AllAdditionals()
		if err != nil {
			return nil, err
		}
		newResolverServersFound := false
		servers = []net.IP{}
		for _, additional := range additionals {
			if additional.Header.Type == dnsmessage.TypeA {
				for _, nameserver := range nameservers {
					if additional.Header.Name.String() == nameserver {
						newResolverServersFound = true
						servers = append(servers, additional.Body.(*dnsmessage.AResource).A[:])
					}
				}
			}
		}

		if !newResolverServersFound {
			for _, nameserver := range nameservers {
				if !newResolverServersFound {
					response, err := dnsQuery(getRootServers(), dnsmessage.Question{Name: dnsmessage.MustNewName(nameserver), Type: dnsmessage.TypeA, Class: dnsmessage.ClassINET}, 3)
					if err != nil {
						fmt.Printf("warning: lookup of nameserver %s failed: %err\n", nameserver, err)
					} else {
						newResolverServersFound = true
						for _, answer := range response.Answers {
							if answer.Header.Type == dnsmessage.TypeA {
								servers = append(servers, answer.Body.(*dnsmessage.AResource).A[:])
							}
						}
					}
				}
			}
		}
	}
	return &dnsmessage.Message{
		Header: dnsmessage.Header{RCode: dnsmessage.RCodeServerFailure},
	}, nil
}

func outgoingDnsQuery(servers []net.IP, question dnsmessage.Question) (*dnsmessage.Parser, *dnsmessage.Header, error) {
	max := ^uint16(0)
	idGenerated, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return nil, nil, err
	}
	message := dnsmessage.Message{
		Header: dnsmessage.Header{
			ID:       uint16(idGenerated.Int64()),
			Response: false,
			OpCode:   dnsmessage.OpCode(0),
		},
		Questions: []dnsmessage.Question{question},
	}
	buf, err := message.Pack()
	if err != nil {
		return nil, nil, err
	}
	var conn net.Conn
	for _, server := range servers {
		conn, err = net.Dial("udp", server.String()+":53")
		if err == nil {
			break
		}
	}
	if conn == nil {
		return nil, nil, fmt.Errorf("failed to make connection to servers: %s", err)
	}

	_, err = conn.Write(buf)
	if err != nil {
		return nil, nil, err
	}

	answer := make([]byte, 512)
	n, err := bufio.NewReader(conn).Read(answer)
	if err != nil {
		return nil, nil, err
	}

	conn.Close()

	var p dnsmessage.Parser

	header, err := p.Start(answer[:n])
	if err != nil {
		return nil, nil, fmt.Errorf("parser start error: %s", err)
	}

	questions, err := p.AllQuestions()
	if err != nil {
		return nil, nil, err
	}
	if len(questions) != len(message.Questions) {
		return nil, nil, fmt.Errorf("answer packet doesn't have the same amount of questions")
	}
	err = p.SkipAllQuestions()
	if err != nil {
		return nil, nil, err
	}

	return &p, &header, nil
}

func getRootServers() []net.IP {
	rootServers := []net.IP{}
	for _, rootServer := range ROOT_SERVERS {
		rootServers = append(rootServers, net.ParseIP(rootServer))
	}
	return rootServers
}

func RespondToBlockIp(pc net.PacketConn, addr net.Addr, buf []byte) {
	for {

		var msg dnsmessage.Message
		if err := msg.Unpack(buf); err != nil {
			fmt.Printf("Erreur when Unpack message : %v\n", err)
			continue
		}

		responseHeader := dnsmessage.Header{
			ID:                 msg.Header.ID,
			Response:           true,
			OpCode:             msg.Header.OpCode,
			Authoritative:      true,
			RecursionAvailable: msg.Header.RecursionDesired,
			RecursionDesired:   true,
			RCode:              dnsmessage.RCodeSuccess,
		}

		var answers []dnsmessage.Resource
		for _, question := range msg.Questions {
			answer := dnsmessage.Resource{
				Header: dnsmessage.ResourceHeader{
					Name:   question.Name,
					Type:   dnsmessage.TypeA,
					Class:  dnsmessage.ClassINET,
					TTL:    300,
				},
				Body: &dnsmessage.AResource{
					A: [4]byte{127, 0, 0, 1},
				},
			}
			answers = append(answers, answer)
		}

		response := dnsmessage.Message{
			Header:   responseHeader,
			Questions: msg.Questions,
			Answers:   answers,
		}

		responseBytes, err := response.Pack()
		if err != nil {
			fmt.Printf("Error during the construction of the message : %v\n", err)
			continue
		}

		pc.WriteTo(responseBytes, addr)
	}
}

func handleBlockDomains(pc net.PacketConn, addr net.Addr, buf []byte, dnsConfig *config.DnsConfig) (bool, error) {
	var msg dnsmessage.Message
	if err := msg.Unpack(buf); err != nil {
		fmt.Printf("Error when unpacking the message : %v\n", err)
		return false, err
	}

	for _, domain := range dnsConfig.DomainsBlocked {
		matchRegex := utils.MatchRegex(`^.*` + domain + `\.$`, msg.Questions[0].Name.String())
		if matchRegex {
			fmt.Printf("----------------------------------------------\n")
			fmt.Printf("Block Ip for [%s]: %s\n", addr.String(), domain)
			fmt.Printf("----------------------------------------------\n")
			RespondToBlockIp(pc, addr, buf)
			
			return true, nil
		}
	}

	return false, nil
}


func handleDNSRequestToProxy(pc net.PacketConn, buf []byte, addr net.Addr) error {

	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		log.Printf("Error when connecting to the Google DNS server: %v", err)
		return err
	}
	defer conn.Close()

	_, err = conn.Write(buf)
	if err != nil {
		log.Printf("Error when sending the request to the Google DNS server: %v", err)
		return err
	}

	responseData := make([]byte, 512)
	n, err := conn.Read(responseData)
	if err != nil {
		log.Printf("Error when receiving the response from the Google DNS server: %v", err)
		return err
	}

	_, err = pc.WriteTo(responseData[:n], addr)
	if err != nil {
		log.Printf("Error when writing the response to the client: %v", err)
		return err
	}

	return nil
}
