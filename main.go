package main

import (
	"encoding/json"
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

// TODO: may be move this to seperate file
type category string

const (
	article category = "articles"
	book    category = "books"
)

type highlightRes struct {
	Count   int
	Sources []source `json:"results"`
}

type source struct {
	Readable_title string
	Category       category
	SourceUrl      *string     `json:"source_url"`
	ImgUrl         string      `json:"cover_image_url"`
	Highlights     []highlight `json:"highlights"`
}

type highlight struct {
	Text string `json:"text"`
	Url  string `json:"readwise_url"`
}

func main() {
	// TODO: make this possible to override using command line
	key, err := readKeys(keyFilePath)
	if err != nil {
		log.Fatal(err)
	}
	data, err := fetchHighlights(key)
	if err != nil {
		log.Fatal(err)
	}
	highlightRes, err := parseHightlightRes(data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(highlightRes)
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

func fetchHighlights(apiKey string) (*[]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", apiUrl+"/export", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Token "+apiKey)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch highlights failed with status: %s", resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return &body, nil
}

func parseHightlightRes(response *[]byte) (*highlightRes, error) {
	var res highlightRes
	err := json.Unmarshal(*response, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
