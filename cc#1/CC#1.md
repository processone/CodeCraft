autoscale: true

# > CodeCraft.go
.....
.....
.....       Write tests with HTTPMock

---

# Introduction

Hello Gophers!

Welcome to our new serie of videos on programming languages.

This time, I will explain how to mock HTTP requests in Go,
  The goal is to be able to test code that triggers HTTP calls very easily.

Let's get started!

---

## Requirements

I have built a web crawler and page analysis tool for data portability.
and needed a tool making easy to write my HTTP tests.

To summarize my requirements, I need a tool: 

- to write stable and reproducible tests: the goal to work consistently no matter what changes on the target page or API.
- to run fast and work offline: be independent from network conditions.
- to be easy to set up 

I didn't find any suitable tool,
so I decided to write my own HTPP Mock tool to address my own needs.

---

## Principles

To address those requirements, I have designed HTTPMock library on two tools:

1. A recorder command-line application to perform real HTTP request and store the behaviour in a scenario file.
2. An instrumented Go HTTP client that can load a scenario and trigger predefined responses in your own code.

---

## Installing the recorder

is very simple

You can install the recorder using `go get` command-line tool:

```bash
go get -u gosrc.io/httpmock/httprec
```

---

## Recording a scenario

is also very simple.

```bash
mkdir fixtures
httprec add fixtures/test1 -u https://www.process-one.net/
httprec add fixtures/test1 -u https://fluux.io/
```

You can store as many requests as you need in your scenario file.

---

## Understanding the scenario files

```
ls -la fixtures/
total 368
drwxr-xr-x  6 mremond  staff     192 Feb 27 18:21 .
drwxr-xr-x  8 mremond  staff     256 Feb 27 18:20 ..
-rw-r--r--  1 mremond  staff   67675 Feb 27 18:20 test1-1-2.html
-rw-r--r--  1 mremond  staff  109578 Feb 27 18:21 test1-2-1.html
-rw-r--r--  1 mremond  staff    3600 Feb 27 18:21 test1.json
-rw-r--r--  1 mremond  staff      47 Feb 27 18:21 test1.url
...
```

The result consists of three types of file:

- the scenario file, which is store in json format
- the url file which store all the entrypoints on separate line. This can be used to regenerate the scenario in case you need to update it after a change of content or API (or if the scenario file format changes in the future)
- the body files, which store the HTTP reply bodies. In this case it is HTML, but it could be images or any type of data.

---

## Using the scenario in a test

The next step is to use the recorded scenario file in one of your test.

In this example, I assume I would like to test a library that returns the number of time a given word is present on an HTTP page.

You cannot really test on the real page, because it would be slow, fragile and would prevent you from developing while offline.

With our recorded scenario it will be easy to do so.

---

## The example library

The library is designed to allow configuring the Go HTTP client.
This is always a good practice to let users of your lib customize the behaviour to suit their needs (like timeouts, tls support, etc).

```go
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
```

---

## Example application using the library

Here is an example application using the lib. When you run it, it will display `Count: 22` (as of today):

```go
package main

import (
	"fmt"
	"mockdemo"
)

func main() {
	counter := mockdemo.NewPageCounter()
	if n, err := counter.Count("https://www.process-one.net", "ProcessOne"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Count:", n)
	}
}
```

---

## The test code (1/3)

Now, we have everything to write the test code.

The test code just have to:
- Set up the HTTPMock and load the scenario
- Set the library to use our client
- Run the test and check the result

---

## The test code (2/3)

```go
package mockdemo_test

import (
	"mockdemo"
	"testing"

	"gosrc.io/httpmock"
)

func TestProcessOne(t *testing.T) {
	// Setup HTTP Mock
	mock := httpmock.NewMock("fixtures/")
	fixtureName := "test1"
	if err := mock.LoadScenario(fixtureName); err != nil {
		t.Errorf("Cannot load fixture scenario %s: %s", fixtureName, err)
		return
	}

	// Setup library
	counter := mockdemo.NewPageCounter()
	counter.Client = mock.Client // Tell the library to use our HTTP client

	// Run test
	n, err := counter.Count("https://www.process-one.net/", "ProcessOne")
	if err != nil {
		t.Error("Unexpected error:", err)
		return
	}

	if n != 23 {
		t.Errorf("Unexpected word count (23), got %d", n)
	}
}
```

---

## The test code (3/3)

The test code run fast and passes.

You can even try running the code with no network and it will still work fine.

---

## Conclusion

As you have seen, adding new HTTP tests to your code thanks to HTTPMock library can be done
in a matter of minutes.

I hope you have found this video useful and will give the library a try.


