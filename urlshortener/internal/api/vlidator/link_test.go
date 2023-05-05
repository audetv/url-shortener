package vlidator

import "testing"

func TestValidLink(t *testing.T) {
	link := "http://localhost"

	err := ValidLink(link)
	if err != nil {
		t.Errorf("ValidLink(%v) = %d; want nill", link, err)
	}
}

func TestInvalidLink(t *testing.T) {
	link := "localhost"

	err := ValidLink(link)
	if err == nil {
		t.Errorf("InvalidLink(%v) = %d; want Error", link, err)
	}
}
