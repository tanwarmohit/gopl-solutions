package main

import (
	"fmt"
	"os"

	"github.com/gopl-solutions/chapter4/ex4.13"

	"github.com/gopl-solutions/chapter4/ex4.12"
)

func xkcdDriver() {
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

func omdbDriver() {
	apikey := os.Args[1]
	title := os.Args[2]
	movie, err := omdb.GetMovie(apikey, title)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error fetching info for: %s.. %v\n", title, err)
		return
	}

	fmt.Printf("Movie: %v\n", *movie)

	if len(movie.Error) > 0 {
		fmt.Fprintf(os.Stderr, "error getting movie: %s\n", movie.Error)
		return
	}

	filename, err := omdb.DownloadPoster(movie)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error downloading poster for: %s.. %v\n", title, err)
	}

	fmt.Printf("Poster downloaded at: %s\n", filename)
}

func main() {
	if len(os.Args) == 3 {
		omdbDriver()
	}
}
