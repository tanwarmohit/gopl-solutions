package xkcd

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	// LatestComicURL get the latest comic.
	LatestComicURL = "https://xkcd.com/info.0.json"

	// RandomComicURL gets the selective comic
	RandomComicURL = "https://xkcd.com/%d/info.0.json"
)

// GetComic will retrieve comic
func GetComic(num int) (*Comic, error) {
	fetchURL := LatestComicURL
	if num > 0 {
		fetchURL = fmt.Sprintf(RandomComicURL, num)
	}

	fmt.Printf("fetching :%s\n", fetchURL)

	resp, err := http.Get(fetchURL)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch comic: %s due to: %s", fetchURL, resp.Status)
	}

	comic := new(Comic)
	if err := json.NewDecoder(resp.Body).Decode(comic); err != nil {
		return nil, err
	}

	return comic, nil
}
