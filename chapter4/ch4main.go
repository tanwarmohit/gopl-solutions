package main

import (
	"fmt"
	"os"

	"github.com/gopl-solutions/chapter4/ex4.12"
)

func main() {
	wordFile := "word.json"
	comicFile := "comic.json"

	wordFilePresent, comicFilePresent := true, true

	if _, err := os.Stat(wordFile); os.IsNotExist(err) {
		fmt.Printf("wordFile: %s not exists\n", wordFile)
		wordFilePresent = false
	}

	if _, err := os.Stat(comicFile); os.IsNotExist(err) {
		fmt.Printf("comicFile: %s not exists\n", comicFile)
		comicFilePresent = false
	}

	if wordFilePresent && comicFilePresent {
		wordData, err := xkcd.ReadWordIndex(wordFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "word file read failed: %v\n", err)
			os.Exit(1)
		}

		comicData, err := xkcd.ReadComicIndex(comicFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "comic file read failed: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("words :%d, comics: %d\n", len(wordData), len(comicData))

	} else {
		if err := xkcd.FetchAndStoreOffline(wordFile, comicFile); err != nil {
			fmt.Fprintf(os.Stderr, "error fetching and saving comic data: %v\n", err)
			os.Exit(1)
		}
	}

}
