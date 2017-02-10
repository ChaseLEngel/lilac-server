package rss

import (
	"bytes"
	"encoding/xml"
	"net/http"
)

type Feed struct {
	XMLName xml.Name `xml:"rss"`
	Channel *Channel `xml:"channel"`
}

type Channel struct {
	XMLName     string  `xml:"channel"`
	Title       string  `xml:"title"`
	Description string  `xml:"description"`
	Link        string  `xml:"link"`
	Items       []*Item `xml:"item"`
}

type Item struct {
	XMLName     string `xml:"item"`
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
}

func Get(url string) (*Channel, error) {
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	channel, err := parse(buf.Bytes())
	if err != nil {
		return nil, err
	}
	return channel, nil
}

func parse(data []byte) (*Channel, error) {
	feed := new(Feed)
	err := xml.Unmarshal(data, &feed)
	if err != nil {
		return nil, err
	}
	return feed.Channel, nil
}
