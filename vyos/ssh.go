package vyos

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
	"strings"
	"syscall"
)

func FetchConfig(user, addr string) ([]byte, error) {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PasswordCallback(func() (string, error) {
				fmt.Print("Password: ")
				bytePassword, err := term.ReadPassword(int(syscall.Stdin))
				if err != nil {
					return "", err
				}
				fmt.Println()
				return strings.TrimSpace(string(bytePassword)), nil
			}),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %v", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run("/opt/vyatta/sbin/vyatta-cfg-cmd-wrapper show"); err != nil {
		return nil, fmt.Errorf("failed to run: %v", err)
	}

	return b.Bytes(), nil
}
