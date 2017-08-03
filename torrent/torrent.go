package torrent

import (
	bencode "github.com/jackpal/bencode-go"
	"os"
)

type File struct {
	Path   []string
	Length int
}

type Info struct {
	Name        string
	PieceLength int
	Length      int
	Pieces      string
	Files       []File
}

type Torrent struct {
	Announce string
	Info     Info
}

func Parse(path string) (Torrent, error) {
	var t Torrent

	r, err := os.Open(path)
	if err != nil {
		return t, err
	}

	if err := bencode.Unmarshal(r, &t); err != nil {
		return t, err
	}
	return t, nil
}
