package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	// TODO: make this possible to override using command line
	key, err := readKeys("./keys.gpg")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(key);
}

func readKeys(keyFilePath string) (string, error) {
    _, err := os.Stat(keyFilePath)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("key file %s does not exist.", keyFilePath)
	}

    cmd := exec.Command("gpg", "--decrypt", keyFilePath)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	key := string(output)
	key = strings.TrimSpace(key)
	return key, nil
}
