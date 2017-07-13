package main

import (
	"fmt"
	"os/exec"

	"github.com/Sirupsen/logrus"
)

func main() {

	currentDir := "/Users/slokas/godev/src/github.com/upmc-enterprises/gen-certs/certs"

	// Generate CA Cert
	logrus.Info("Starting create CA...")
	cmdCA1 := exec.Command("cfssl", "genkey", "-initca", "config/ca-csr.json")
	cmdCA2 := exec.Command("cfssljson", "-bare", "certs/ca")
	_, err := pipeCommands(cmdCA1, cmdCA2)
	if err != nil {
		logrus.Error(err)
	}

	// Generate Node Cert
	logrus.Info("Starting create Node...")
	cmdNode1 := exec.Command("cfssl", "gencert", "-ca", "certs/ca.pem", "-ca-key", "certs/ca-key.pem", "-config", "config/ca-config.json", "config/req-csr.json")
	cmdNode2 := exec.Command("cfssljson", "-bare", "certs/node")
	_, err = pipeCommands(cmdNode1, cmdNode2)
	if err != nil {
		logrus.Error(err)
	}

	logrus.Info("Starting cmdConvertCA...")
	cmdConvertCA := exec.Command("openssl", "pkcs12", "-export", "-inkey", "certs/ca-key.pem", "-in", "certs/ca.pem", "-out", "certs/ca.pkcs12", "-password", "pass:changeit")
	out, err := cmdConvertCA.Output()
	if err != nil {
		logrus.Error(string(out))
	}

	logrus.Info("Starting cmdConvertNode...")
	cmdConvertNode := exec.Command("openssl", "pkcs12", "-export", "-inkey", "certs/node-key.pem", "-in", "certs/node.pem", "-out", "certs/node.pkcs12", "-password", "pass:changeit")
	out, err = cmdConvertNode.Output()
	if err != nil {
		logrus.Error(string(out))
	}

	logrus.Info("Starting cmdCAJKS...")
	cmdCAJKS := exec.Command("keytool", "-importkeystore", "-srckeystore", fmt.Sprintf("%s/ca.pkcs12", currentDir), "-srcalias", "1", "-destkeystore", fmt.Sprintf("%s/truststore.jks", currentDir),
		"-storepass", "changeit", "-srcstoretype", "pkcs12", "-srcstorepass", "changeit", "-destalias", "elasticsearch-ca")
	out, err = cmdCAJKS.Output()
	if err != nil {
		logrus.Error(string(out))
	}

	logrus.Info("Starting cmdNodeJKS...")
	cmdNodeJKS := exec.Command("keytool", "-importkeystore", "-srckeystore", fmt.Sprintf("%s/node.pkcs12", currentDir), "-srcalias", "1", "-destkeystore", fmt.Sprintf("%s/node-keystore.jks", currentDir),
		"-storepass", "changeit", "-srcstoretype", "pkcs12", "-srcstorepass", "changeit", "-destalias", "elasticsearch-node")
	out, err = cmdNodeJKS.Output()
	if err != nil {
		logrus.Error(string(out))
	}
}

// https://gist.github.com/dagoof/1477401
func pipeCommands(commands ...*exec.Cmd) ([]byte, error) {
	for i, command := range commands[:len(commands)-1] {
		out, err := command.StdoutPipe()
		if err != nil {
			return nil, err
		}
		command.Start()
		commands[i+1].Stdin = out
	}
	final, err := commands[len(commands)-1].Output()
	if err != nil {
		return nil, err
	}
	return final, nil
}
