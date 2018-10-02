package xkcd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

type WordIndex map[string]map[int]bool
type ComicIndex map[int]Comic

// FetchAndStoreOffline will read all comics and create an offline indexing of all words in transcripts
func FetchAndStoreOffline(wordIndexFile string, comicIndexFile string) error {
	comic, err := GetComic(-1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while fetching latest comic: %v\n", err)
		return err
	}

	fmt.Printf("%v\n", *comic)

	wIndex := make(WordIndex)
	cIndex := make(ComicIndex)

	for i := 1; i < comic.Num; i++ {
		comic, err := GetComic(i)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error while fetching comic Id: %d error: %v\n", i, err)
			continue
		}

		// Comic Number -> comic object mapping
		cIndex[comic.Num] = *comic

		transcriptScanner := bufio.NewScanner(strings.NewReader(comic.Transcript))
		transcriptScanner.Split(customSplit)

		for transcriptScanner.Scan() {
			token := strings.ToLower(transcriptScanner.Text())
			if _, ok := wIndex[token]; !ok {
				wIndex[token] = make(map[int]bool)
			}
			wIndex[token][comic.Num] = true
		}

		if len(comic.Alt) > 0 {
			altScanner := bufio.NewScanner(strings.NewReader(comic.Alt))
			altScanner.Split(customSplit)

			for altScanner.Scan() {
				token := strings.ToLower(altScanner.Text())
				if _, ok := wIndex[token]; !ok {
					wIndex[token] = make(map[int]bool)
				}
				wIndex[token][comic.Num] = true
			}
		}
	}

	// marshaling and saving wordIndex data
	wordData, err := json.Marshal(wIndex)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error marshaling wordIndex: %v\n", err)
		return err
	}

	if err := ioutil.WriteFile(wordIndexFile, wordData, 0777); err != nil {
		fmt.Fprintf(os.Stderr, "error writing word index data: %v\n", err)
		return err
	}

	comicData, err := json.Marshal(cIndex)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error marshaling comicIndex: %v\n", err)
		return err
	}

	if err := ioutil.WriteFile(comicIndexFile, comicData, 0777); err != nil {
		fmt.Fprintf(os.Stderr, "error writing comic index data: %v\n", err)
		return err
	}

	fmt.Printf("%v\n", cIndex)

	return nil
}

// ReadWordIndex reads the data from json file
func ReadWordIndex(wordIndexFile string) (WordIndex, error) {
	data, err := ioutil.ReadFile(wordIndexFile)
	if err != nil {
		return nil, err
	}

	wordIndex := make(WordIndex)
	if err := json.Unmarshal(data, &wordIndex); err != nil {
		return nil, err
	}

	return wordIndex, nil
}

// ReadComicIndex reads the data from json file
func ReadComicIndex(comicIndexFile string) (ComicIndex, error) {
	data, err := ioutil.ReadFile(comicIndexFile)
	if err != nil {
		return nil, err
	}

	comicIndex := make(ComicIndex)
	if err := json.Unmarshal(data, &comicIndex); err != nil {
		return nil, err
	}

	return comicIndex, nil
}

func customSplit(data []byte, atEOF bool) (advance int, token []byte, err error) {
	currentIndex, start, end := 0, 0, 0
	for currentIndex < len(data) {
		r, size := utf8.DecodeRune(data[currentIndex:])
		currentIndex += size
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			start = currentIndex - size
			break
		}
	}

	for currentIndex < len(data) {
		r, size := utf8.DecodeRune(data[currentIndex:])
		currentIndex += size
		if !(unicode.IsLetter(r) || unicode.IsDigit(r)) {
			end = currentIndex - size
			break
		}
	}

	if end > start {
		token = data[start:end]
	}

	//fmt.Printf("start: %d, end: %d, data: %v\n", start, end, data)

	return currentIndex, token, nil
}
