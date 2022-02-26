package main

import (
	"bufio"
	"crypto/md5"
	"embed"
	"errors"
	"fmt"
	"io"
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
	checksum, err := faviconChecksum(pathOrUrl)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
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

func readFavicon(pathOrUrl string) ([]byte, error) {
	var file *os.File

	// get file
	if strings.HasPrefix(pathOrUrl, "http://") || strings.HasPrefix(pathOrUrl, "https://") {
		// from URL
		if !strings.HasSuffix(pathOrUrl, ".ico") {
			return nil, errors.New("The URL is not a favicon")
		}

		var err error
		file, err = downloadFaviconToTempfile(pathOrUrl)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		defer os.Remove(file.Name())
	} else {
		// from filesystem
		if !strings.HasSuffix(pathOrUrl, ".ico") {
			return nil, errors.New("The file is not a favicon")
		}

		var err error
		file, err = os.Open(pathOrUrl)
		if err != nil {
			return nil, errors.New("Could not open favicon file: " + pathOrUrl)
		}
		defer file.Close()
	}

	// read its contents
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, errors.New("Could not read contents of favicon file: " + pathOrUrl)
	}

	return data, nil
}

func downloadFaviconToTempfile(faviconUrl string) (*os.File, error) {
	// parse URL
	u, err := url.Parse(faviconUrl)
	if err != nil {
		return nil, errors.New("Could not parse URL: " + faviconUrl)
	}

	// download favicon from URL
	response, err := http.Get(u.String())
	if err != nil {
		return nil, errors.New("Error while downloading favicon")
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return nil, errors.New("Error while downloading favicon: HTTP status code " + strconv.Itoa(response.StatusCode))
	}

	// create tempfile to store the favicon
	tempFile, err := os.CreateTemp("", "favicheck*.ico")
	if err != nil {
		return nil, errors.New("Error while creating tempfile")
	}

	// copy favicon to tempfile
	_, err = io.Copy(tempFile, response.Body)
	if err != nil {
		return nil, errors.New("Error while copying download to tempfile")
	}

	// seek back to beginning after copy
	tempFile.Seek(0, io.SeekStart)

	return tempFile, nil
}

// Get the favicon file's md5 checksum
func faviconChecksum(pathOrUrl string) (string, error) {
	// get favicon data
	faviconData, err := readFavicon(pathOrUrl)
	if err != nil {
		return "", err
	}

	// calculate its checksum
	checksumBytes := md5.Sum(faviconData)
	checksumString := fmt.Sprintf("%x", checksumBytes)

	return checksumString, nil
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
