package baseClass

import (
	"encoding/json"
	"log"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
)

type Dialer struct {
	client *ssh.Client
}

func (v *Dialer) Dial(address string) (net.Conn, error) {
	return v.client.Dial("tcp", address)
}

func SSHConn() *ssh.Client {
	data, err := os.ReadFile("config.json")
	if err != nil {
		log.Println("Failed to open config:", err)
		return nil
	}
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Println("Failed to open config:", err)
		return nil
	}

	if !config.UseSSH {
		return nil
	}
	key, err := os.ReadFile(config.SSH.PrivateKey)
	if err != nil {
		log.Println("Failed to open config:", err)
		return nil
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Println("Failed to open config:", err)
		return nil
	}
	ssh_config := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 忽略主机密钥
	}
	// 作为客户端连接SSH服务器
	client, err := ssh.Dial("tcp", config.SSH.Host+":"+config.SSH.Port, ssh_config)
	if err != nil {
		log.Println("Failed to dial: ", err)
		return nil
	}
	return client
}
