package utils

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

const pubKey = "/.ssh/id_rsa.pub"
const fileToEnc = "~/side-projects/data.txt"

// EncryptFile creates a (temporary) encrypted file for the recipient
func EncryptFile(file *os.File, recipient string) (string, error) {
	return encryptWithGPGBinary(file, recipient)
}

// DecryptFile creates a (temporary) encrypted file for the recipient
func DecryptFile(file io.Reader) (string, error) {
	reader, err := decryptWithGPGBinary(file)
	if err != nil {
		log.Errorf("Error decrypting: %s", err.Error())
		return "", err
	}
	cipherBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Errorf("Error reading bytes: %s", err.Error())
		return "", err
	}
	cipher := string(cipherBytes)
	return cipher, nil
}

func encryptWithGPGBinary(file *os.File, recipient string) (string, error) {
	args := []string{
		"--yes",
		"--encrypt",
		"-a",
		"-r",
		recipient,
	}
	cmd := exec.Command("gpg", args...)
	var stdout, stderr bytes.Buffer
	// cmd.Stdin = bytes.NewReader(dataKey)
	cmd.Stdin = file
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Errorf("Error encrypting message with gpg binary: %s", err.Error())
		return "", err
	}
	return string(stdout.Bytes()), nil
}

func decryptWithGPGBinary(file io.Reader) (io.Reader, error) {
	args := []string{
		"--use-agent",
		"-d",
	}
	cmd := exec.Command("gpg", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdin = file
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Errorf("Error decrypting message with gpg binary: %s", err.Error())
		return nil, err
	}
	return bytes.NewReader(stdout.Bytes()), nil
}
