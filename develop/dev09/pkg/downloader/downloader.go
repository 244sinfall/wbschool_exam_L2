package downloader

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

const MaxDepth = 1

func SeparateLinkToSlice(url string) ([]string, error) {
	regex, _ := regexp.Compile("([A-Z.a-z0-9-]+)")
	link := regex.FindAllString(url, -1)
	if len(link) < 2 {
		return nil, errors.New("incorrect link")
	}
	return link, nil
}

func ValidatePath(path *string) error {
	fs, err := os.Open(*path)
	fStats, _ := fs.Stat()
	if err != nil || fStats.IsDir() {
		*path += "\\index.html"
		fs, err = os.Open(*path)
		fStats, _ = fs.Stat()
		if err != nil || fStats.IsDir() {
			return errors.New("unable to open file path:" + *path + err.Error())
		}
	}
	_ = fs.Close()
	if fStats.Size() > 0 {
		return errors.New("file already exist and non empty")
	}
	return nil
}

func StartDownload(url string, depth int) error {
	linkSliced, err := SeparateLinkToSlice(url)
	if err != nil {
		return err
	}
	var protocol, baseUrl, filePath string
	pureLink := linkSliced[1:]
	protocol = linkSliced[0] + "://"
	baseUrl = pureLink[0]
	ValidateLink(pureLink)
	filePath = strings.Join(pureLink, "\\")
	err = ValidatePath(&filePath)
	if err != nil {
		return err
	}
	fs, _ := os.OpenFile(filePath, os.O_RDWR, 0666)
	fmt.Printf("Downloading %v...\n", url)
	res, err := http.Get(url)
	if err != nil {
		return errors.New("unable to download file:" + url + err.Error())
	}
	_, err = io.Copy(fs, res.Body)
	if err != nil {
		return err
	}
	receivedFile, _ := os.ReadFile(fs.Name())
	_ = res.Body.Close()
	_ = fs.Close()
	linkParser, _ := regexp.Compile(`([a-zA-Z]+:/)?/([A-Za-z0-9?=&/-])?/?([A-Za-z./?=&0-9-]+\.[A-Za-z.?=&0-9-]+)`)
	newLinks := linkParser.FindAllString(string(receivedFile), -1)
	if depth < MaxDepth {
		for _, l := range newLinks {
			var err error
			var newDepth = depth + 1
			if strings.HasPrefix(l, "/") {
				err = StartDownload(protocol+baseUrl+l, newDepth)
			} else {
				err = StartDownload(l, newDepth)
			}
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}

func ValidateLink(urlParts []string) {
	baseFolder, _ := os.Getwd()
	for i, v := range urlParts {
		if i == len(urlParts)-1 {
			if strings.Contains(v, ".") && len(urlParts) != 1 {
				if _, err := os.Stat(v); errors.Is(err, os.ErrNotExist) {
					f, _ := os.Create(v)
					_ = f.Close()
				}
				break
			}
		}
		err := os.Mkdir(v, 0666)
		err = os.Chdir(v)
		if err != nil {
			fmt.Println("unable to manage folder:", v)
		}
	}
	if !strings.Contains(urlParts[len(urlParts)-1], ".") || len(urlParts) == 1 {
		if _, err := os.Stat("index.html"); errors.Is(err, os.ErrNotExist) {
			f, _ := os.Create("index.html")
			_ = f.Close()
		}
	}
	_ = os.Chdir(baseFolder)
}
