package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

const (
	keyFilePath = "./keys.gpg"
	apiUrl      = "https://readwise.io/api/v2"
)

func main() {
	// TODO: make this possible to override using command line
	key, err := readKeys(keyFilePath)
	if err != nil {
		log.Fatal(err)
	}
	res, err := fetchHighlights(key);
	if err != nil {
		log.Fatal(err);
	}
	fmt.Print(res)
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

func fetchHighlights(apiKey string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", apiUrl + "/export", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Token "+apiKey)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("fetch highlights failed with status: %s", resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body);
	if err != nil {
		return "", err
	}
	return string(body), nil
}
