package hostDns

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

func Reset() {

	cmd := exec.Command("sh", "-c", `cat /etc/resolv.conf | grep nameserver`)
	log.Println("command generated: " + cmd.String())
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("[error]: Error: %q", err.Error())
		log.Printf("[execute_command]: Error executing command: %q", string(output))
		return
	}

	now := time.Now()
	datetime := now.Format("2006-01-02 15:04:05")
	f, err := os.Create(fmt.Sprintf("dns-config/backup/%s.txt", datetime))
    if err != nil {
        fmt.Println("Erreur lors de la création du fichier :", err)
        return
    }
	f.Write([]byte(string(output)))
    defer f.Close()
	
	cmd2 := exec.Command("sh", "-c", `networksetup -setdnsservers Wi-Fi 127.0.0.1`)
	log.Println("command generated: " + cmd2.String())
	output, err = cmd2.CombinedOutput()
	if err != nil {
		log.Printf("[error]: Error: %q", err.Error())
		log.Printf("[execute_command]: Error executing command: %q", string(output))
		return
	}

	log.Printf("dns server well configured")
}
