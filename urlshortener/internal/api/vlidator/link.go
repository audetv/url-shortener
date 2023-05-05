package vlidator

import (
	"fmt"
	"regexp"
	"strings"
)

// ValidLink - check that the link we're creating a shortlink for is an absolute URL path
func ValidLink(link string) error {
	r := regexp.MustCompile("^(http|https)://")

	link = strings.TrimSpace(link)
	// log.Printf("checking for valid link: %s", link)
	// Check if string matches the regex
	if r.MatchString(link) {
		return nil
	}
	return fmt.Errorf("invalid link: %s", link)
}
