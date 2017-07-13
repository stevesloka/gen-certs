package main

import (
	"fmt"
	"os/exec"

	"github.com/Sirupsen/logrus"
)

func main() {

	certsDir := "/Users/slokas/godev/src/github.com/upmc-enterprises/gen-certs/certs"
	configDir := "/Users/slokas/godev/src/github.com/upmc-enterprises/gen-certs/config"

	// Generate CA Cert
	logrus.Info("Starting create CA...")
	cmdCA1 := exec.Command("cfssl", "genkey", "-initca", fmt.Sprintf("%s/ca-csr.json", configDir))
	cmdCA2 := exec.Command("cfssljson", "-bare", "certs/ca")
	_, err := pipeCommands(cmdCA1, cmdCA2)
	if err != nil {
		logrus.Error(err)
	}

	// Generate Node Cert
	logrus.Info("Starting create Node...")
	cmdNode1 := exec.Command("cfssl", "gencert", "-ca", fmt.Sprintf("%s/ca.pem", certsDir), "-ca-key", fmt.Sprintf("%s/ca-key.pem", certsDir), "-config", fmt.Sprintf("%s/ca-config.json", configDir), fmt.Sprintf("%s/req-csr.json", configDir))
	cmdNode2 := exec.Command("cfssljson", "-bare", "certs/node")
	_, err = pipeCommands(cmdNode1, cmdNode2)
	if err != nil {
		logrus.Error(err)
	}

	logrus.Info("Starting cmdConvertCA...")
	cmdConvertCA := exec.Command("openssl", "pkcs12", "-export", "-inkey", fmt.Sprintf("%s/ca-key.pem", certsDir), "-in", fmt.Sprintf("%s/ca.pem", certsDir), "-out", fmt.Sprintf("%s/ca.pkcs12", certsDir), "-password", "pass:changeit")
	out, err := cmdConvertCA.Output()
	if err != nil {
		logrus.Error(string(out))
	}

	logrus.Info("Starting cmdConvertNode...")
	cmdConvertNode := exec.Command("openssl", "pkcs12", "-export", "-inkey", fmt.Sprintf("%s/node-key.pem", certsDir), "-in", fmt.Sprintf("%s/node.pem", certsDir), "-out", fmt.Sprintf("%s/node.pkcs12", certsDir), "-password", "pass:changeit")
	out, err = cmdConvertNode.Output()
	if err != nil {
		logrus.Error(string(out))
	}

	logrus.Info("Starting cmdCAJKS...")
	cmdCAJKS := exec.Command("keytool", "-importkeystore", "-srckeystore", fmt.Sprintf("%s/ca.pkcs12", certsDir), "-srcalias", "1", "-destkeystore", fmt.Sprintf("%s/truststore.jks", certsDir),
		"-storepass", "changeit", "-srcstoretype", "pkcs12", "-srcstorepass", "changeit", "-destalias", "elasticsearch-ca")
	out, err = cmdCAJKS.Output()
	if err != nil {
		logrus.Error(string(out))
	}

	logrus.Info("Starting cmdNodeJKS...")
	cmdNodeJKS := exec.Command("keytool", "-importkeystore", "-srckeystore", fmt.Sprintf("%s/node.pkcs12", certsDir), "-srcalias", "1", "-destkeystore", fmt.Sprintf("%s/node-keystore.jks", certsDir),
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
