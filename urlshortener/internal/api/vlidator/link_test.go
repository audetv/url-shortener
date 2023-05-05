package vlidator

import (
	"fmt"
	"testing"
)

func TestValidLink(t *testing.T) {
	link := "http://localhost"

	err := ValidLink(link)
	if err != nil {
		t.Errorf("ValidLink(%v) = %d; want nill", link, err)
	}
	if err == nil {
		t.Logf("Success !")
	}
}

func TestInvalidLink(t *testing.T) {
	link := "localhost"
	expect := fmt.Sprintf("invalid link: %s", link)

	err := ValidLink(link)

	if err.Error() == expect {
		t.Logf("Success !")
	} else {
		t.Errorf("InvalidLink(%v) = %d; want Error", link, err)
	}
}
