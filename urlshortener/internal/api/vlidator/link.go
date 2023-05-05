package vlidator

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// ValidLink - check that the link we're creating a shortlink for is a absolute URL path
func ValidLink(link string) error {
	r, err := regexp.Compile("^(http|https)://")
	if err != nil {
		return err
	}
	link = strings.TrimSpace(link)
	// log.Printf("checking for valid link: %s", link)
	// Check if string matches the regex
	if r.MatchString(link) {
		return nil
	}
	return errors.New(fmt.Sprintf("invalid link: %s", link))
}
