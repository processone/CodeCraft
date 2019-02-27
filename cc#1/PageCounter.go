package mockdemo

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

type PageCounter struct {
	Client *http.Client
}

func NewPageCounter() PageCounter {
	return PageCounter{http.DefaultClient}
}

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
