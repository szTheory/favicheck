package main

import (
	"bufio"
	"crypto/md5"
	"embed"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var (
	//go:embed data/database.txt
	databaseFile embed.FS
)

func main() {
	if len(os.Args) != 2 {
		printUsage()
		os.Exit(1)
	}

	pathOrUrl := os.Args[1]
	checksum := faviconChecksum(pathOrUrl)
	database := buildDatabase()

	// Find which framework matches the favicon's checksum
	if matchingFramework, ok := database[checksum]; ok {
		fmt.Println("Web framework:", matchingFramework)
	} else {
		fmt.Println("No matching web framework for this favicon")
	}
}

func printUsage() {
	fmt.Println("Usage: favicheck <filepath|url>")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  favicheck ~/Downloads/favicon.ico")
	fmt.Println("  favicheck https://static-labs.tryhackme.cloud/sites/favicon/images/favicon.ico")
}

func readFavicon(pathOrUrl string) []byte {
	var file *os.File

	// get file
	if strings.HasPrefix(pathOrUrl, "http://") || strings.HasPrefix(pathOrUrl, "https://") {
		// from URL
		file = downloadFaviconToTempfile(pathOrUrl)
		defer file.Close()
		defer os.Remove(file.Name())
	} else {
		// from filesystem
		var err error
		file, err = os.Open(pathOrUrl)
		if err != nil {
			panic("Could not open favicon file: " + pathOrUrl)
		}
		defer file.Close()
	}

	// read its contents
	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic("Could not read contents of favicon file: " + pathOrUrl)
	}

	return data
}

func downloadFaviconToTempfile(faviconUrl string) *os.File {
	// parse URL
	u, err := url.Parse(faviconUrl)
	if err != nil {
		panic("Could not parse URL: " + faviconUrl)
	}

	// download favicon from URL
	response, err := http.Get(u.String())
	if err != nil {
		panic("Error while downloading favicon")
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		panic("Error while downloading favicon: HTTP status code " + strconv.Itoa(response.StatusCode))
	}

	// create tempfile to store the favicon
	f, err := os.CreateTemp("", "example")
	if err != nil {
		panic("Error while creating tempfile")
	}

	// copy favicon to tempfile
	_, err = io.Copy(f, response.Body)
	if err != nil {
		panic(err)
	}

	return f
}

// Get the favicon file's md5 checksum
func faviconChecksum(pathOrUrl string) string {
	// get favicon data
	faviconData := readFavicon(pathOrUrl)

	// calculate its checksum
	checksumBytes := md5.Sum(faviconData)
	checksumString := fmt.Sprintf("%x", checksumBytes)

	return checksumString
}

// Build a database of favicon checksums to web framework names
func buildDatabase() map[string]string {
	// open file
	file, err := os.Open("data/database.txt")
	defer file.Close()
	if err != nil {
		panic("Could not open database.txt")
	}

	// read lines
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// create the map
	database := make(map[string]string)
	for _, line := range lines {
		split := strings.Split(line, ":")
		faviconChecksum := split[0]
		frameworkName := split[1]
		database[faviconChecksum] = frameworkName
	}

	return database
}
