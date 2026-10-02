package hasher

import (
	"io"
	"os"

	"github.com/zeebo/xxh3"
)

type fileInfo struct {
	path string
	size int64
	hash [16]byte // XXH3-128 output
}

func HashFile(path string) ([16]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return [16]byte{}, err
	}
	defer f.Close()

	h := xxh3.New()
	if _, err := io.Copy(h, f); err != nil {
		return [16]byte{}, err
	}
	sum := h.Sum128()
	return sum.Bytes(), nil
}
