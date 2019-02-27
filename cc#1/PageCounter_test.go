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

	if n != 22 {
		t.Errorf("Unexpected word count (22), got %d", n)
	}
}
