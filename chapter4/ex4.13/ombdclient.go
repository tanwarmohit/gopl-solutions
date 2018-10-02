package omdb

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/mozillazg/go-slugify"
)

// Movie data structure for storing omdb movie attributes
type Movie struct {
	Title    string
	Year     string
	Poster   string
	Response string
	Error    string
}

const omdbURL = "https://omdbapi.com/?apikey=%s&t=%s"

// GetMovie downloads the moview data from omdb endpoint
func GetMovie(apiKey string, title string) (*Movie, error) {
	fetchURL := fmt.Sprintf(omdbURL, apiKey, title)
	resp, err := http.Get(fetchURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error fetching movie %v\n", err)
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http replied with status: %s", resp.Status)
	}

	movie := new(Movie)

	if err := json.NewDecoder(resp.Body).Decode(movie); err != nil {
		return nil, err
	}

	return movie, nil
}

// DownloadPoster downloads the poster from http endpopint
func DownloadPoster(movie *Movie) (string, error) {
	fmt.Printf("downloading poster from: %s\n", movie.Poster)
	resp, err := http.Get(movie.Poster)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error downloading file :%s", resp.Status)
	}

	extension := filepath.Ext(movie.Poster)
	title := slugify.Slugify(movie.Title)
	filename := fmt.Sprintf("%s_[%s]%s", title, movie.Year, extension)

	fmt.Printf("saving %s...\n", filename)
	out, err := os.Create(filename)
	if err != nil {
		return "", err
	}

	defer out.Close()

	writer := bufio.NewWriter(out)
	size, err := writer.ReadFrom(resp.Body)
	if err != nil {
		return "", err
	}

	fmt.Printf("Write successfully: %d bytes\n", size)
	return filename, nil
}
