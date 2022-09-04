package main

import (
	"fmt"
	"os"
	"wget/pkg/downloader"
)

/*
=== Утилита wget ===

# Реализовать утилиту wget с возможностью скачивать сайты целиком

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/
func main() {
	if len(os.Args) != 2 {
		_, _ = os.Stderr.WriteString("No download link provided.\n")
		os.Exit(1)
	}
	url := os.Args[1]
	err := downloader.StartDownload(url, 0)
	if err != nil {
		fmt.Println(err)
	}
}
