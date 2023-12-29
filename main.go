package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

const (
	defaultKeyPath = "./keys.gpg"
	apiUrl      = "https://readwise.io/api/v2"
	defaultReadwiseDir = "./readwise"
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
	Title string `json:"readable_title"`
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
	args := os.Args
	rootDir := defaultReadwiseDir
	keyfile := defaultKeyPath
	if hasFlag(args, "--h") || hasFlag(args, "-help") {
		fmt.Println("Usage readwiseSync [path_to_valut, path_to_key_file]")
		return
	}
	if len(args) == 3 {
		rootDir = args[1]
		keyfile = args[2]
	}
	key, err := readKeys(keyfile)
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
	for _, source := range highlightRes.Sources  {
		if source.Highlights == nil || len(source.Highlights) == 0 {
			continue
		}
		err = writeSource(source, rootDir)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func hasFlag(args []string, flag string) bool {
	for _, each := range args {
		if each == flag {
			return true
		}
	}
	return false
}

func orgTemplate() (*template.Template, error) {
	orgTemplate := `#+title: {{.Title}}
#+category: {{.Category}}
* Highlights
{{range .Highlights}}- {{.Text}} [[{{.Url}}][(ref)]]
{{end}}
* Source:
+ {{with .SourceUrl}}[[{{.}}][url]]{{else}}N/A{{end}}
`
	return template.New("orgTemplate").Parse(orgTemplate)
}

func writeSource(source source, rootDir string) error {
	template, err := orgTemplate()
	if err != nil {
		return err
	}
	file, err := createOrgFile(rootDir, source)
	if err != nil {
		return err
	}
	defer file.Close()
	return writeWithTemplate(template, source, file)
}

func createOrgFile(rootDir string, src source) (*os.File, error) {
	outputDirectory := filepath.Join(rootDir, string(src.Category))
	err := os.MkdirAll(outputDirectory, os.ModePerm)
	if err != nil {
		return nil, err
	}
	outputFilePath := filepath.Join(outputDirectory, sanitizeFileName(src.Title, "org"))
	return os.Create(outputFilePath)
}

func writeWithTemplate(template *template.Template, src source, writer io.Writer) error {
	return template.Execute(writer, src)
}

func sanitizeFileName(title, ext string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9._-]`)
	sanitized := re.ReplaceAllString(title, "_")

	if len(sanitized) > 255 {
		sanitized = sanitized[:255]
	}

	return sanitized + "." + ext
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
