package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// Metadata holds the metadata information for a URL
type Metadata struct {
	URL        string
	LinkCount  int
	ImageCount int
	LastFetch  time.Time
}

func fetchHTML(urlStr string) (*goquery.Document, error) {
	response, err := http.Get(urlStr)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("failed to fetch: %s", response.Status)
	}

	return goquery.NewDocumentFromReader(response.Body)
}

func saveHTML(urlStr string, doc *goquery.Document) error {
	u, err := url.Parse(urlStr)
	if err != nil {
		return err
	}

	fileName := strings.Replace(u.Hostname()+u.Path, "/", "_", -1) + ".html"
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, strings.NewReader(doc.Selection.Text()))
	if err != nil {
		return err
	}

	return nil
}

func fetchAndSave(urlStr string) error {
	doc, err := fetchHTML(urlStr)
	if err != nil {
		return err
	}

	err = saveHTML(urlStr, doc)
	if err != nil {
		return err
	}

	return nil
}

func getMetadata(url string) (Metadata, error) {
	var metadata Metadata

	fileInfo, err := os.Stat(strings.Replace(url, "https://", "", 1) + ".html")
	if err != nil {
		return metadata, err
	}

	file, err := os.Open(strings.Replace(url, "https://", "", 1) + ".html")
	if err != nil {
		return metadata, err
	}
	defer file.Close()

	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		return metadata, err
	}

	imgCount := 0
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		imgCount++
	})

	linkCount := 0
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		linkCount++
	})

	metadata.URL = url
	metadata.ImageCount = imgCount
	metadata.LinkCount = linkCount
	metadata.LastFetch = fileInfo.ModTime()

	return metadata, nil
}

func main() {
	var metadataFlag bool
	flag.BoolVar(&metadataFlag, "metadata", false, "Flag to retrieve metadata")
	flag.Parse()

	urls := flag.Args()

	if metadataFlag {
		for _, url := range urls {
			metadata, err := getMetadata(url)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Printf("Metadata for %s:\n", metadata.URL)
			fmt.Printf("Link_NUM: %d\n", metadata.LinkCount)
			fmt.Printf("Image_NUM: %d\n", metadata.ImageCount)
			fmt.Printf("Last_fetch: %s\n", metadata.LastFetch)
		}
	} else {
		for _, url := range urls {
			err := fetchAndSave(url)
			if err != nil {
				fmt.Printf("Error fetching %s: %v\n", url, err)
			} else {
				fmt.Printf("Fetched and saved: %s\n", url)
			}
		}
	}
}
