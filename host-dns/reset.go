package hostDns

import (
	"dns-focus/utils"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

const BACKUP_DIR = "host-dns/backup"

func SetDnsForMac() {

	cmd := exec.Command("sh", "-c", `cat /etc/resolv.conf | grep nameserver`)
	log.Println("command generated: " + cmd.String())
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("[error]: Error: %q", err.Error())
		log.Printf("[execute_command]: Error executing command: %q", string(output))
		return
	}

	now := time.Now()
	datetime := now.Format("2006-01-02T15:04:05")
	f, err := os.Create(fmt.Sprintf("%s/%s.txt", BACKUP_DIR, datetime))
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

func ResetDnsForMac() {

	file, err := getlastBackupFile()
	if err != nil {
		log.Printf("[Error] %s", err.Error())
	}

	lines, err := utils.ReadFile(BACKUP_DIR + "/" +file.Name())
	log.Println(string(lines))
	if err != nil {
		log.Printf("[error]: Error: %q", err.Error())
	}
	ipAddresses := ""
	for _,line := range strings.Split(string(strings.TrimSpace(string(lines))), "\n") {
		adress := strings.Split(line, " ")
		ipAddresses += adress[1]
	}

	cmd := exec.Command("sh", "-c", `networksetup -setdnsservers Wi-Fi ` + ipAddresses)
	log.Println("command generated: " + cmd.String())
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("[error]: Error: %q", err.Error())
		log.Printf("[execute_command]: Error executing command: %q", string(output))
		return
	}

	log.Printf("dns server ips restored %s", ipAddresses)
}

func getlastBackupFile() (fs.FileInfo, error) {

	files, err := ioutil.ReadDir(BACKUP_DIR)
	if err != nil {
		fmt.Println("Erreur lors de la lecture du répertoire:", err)
		os.Exit(1)
	}

	var latestFile os.FileInfo
	var latestTime time.Time

	for _, file := range files {
		if !file.IsDir() {
			fileTime, err := time.Parse("2006-01-02T15:04:05.txt", file.Name())
			if err == nil {
				if fileTime.After(latestTime) {
					latestFile = file
					latestTime = fileTime
				}
			}
		}
	}

	if latestFile != nil {
		fmt.Println("Most recent file:", latestFile.Name())
		return latestFile, nil
	} else {
		fmt.Println("No file with a valid datetime name found.")
		return nil, errors.New("No file with a valid datetime name found.")
	}
}
