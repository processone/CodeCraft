package mockdemo

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

// PageCounter is the core of the library.
type PageCounter struct {
	Client *http.Client
}

// NewPageCounter creates a structure with default parameters.
func NewPageCounter() PageCounter {
	return PageCounter{http.DefaultClient}
}

// Count return the number of times a word appear on a webpage.
func (p PageCounter) Count(uri, word string) (n int, err error) {
	resp, err := p.Client.Get(uri)
	if err != nil {
		return n, err
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return n, err
	}

	n = bytes.Count(content, []byte(word))
	return n, nil
}
