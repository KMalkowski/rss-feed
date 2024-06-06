package feedfetcher

import (
	"encoding/xml"
	"fmt"
	"net/http"
)

func FetchFeed(url string) (RssFeed, error) {
	resp, err := http.Get(url)
	if err != nil {
		return RssFeed{}, fmt.Errorf("could not fetche the rss feed")
	}

	xmlBody := RssFeed{}
	err = xml.NewDecoder(resp.Body).Decode(&xmlBody)
	if err != nil {
		return RssFeed{}, fmt.Errorf("could not read the rss reposnse body")
	}

	return xmlBody, nil
}
