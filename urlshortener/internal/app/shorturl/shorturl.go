package shorturl

import (
	"math/rand"
)

const letterBytes = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type ShortUrl string

func New() *ShortUrl {
	shortUrlLen := 7
	shortUrl := randStringBytes(shortUrlLen)
	return (*ShortUrl)(&shortUrl)
}

func Parse(s string) *ShortUrl {
	// здесь можно сделать проверку, что в строке содержаться символы из letterBytes
	return (*ShortUrl)(&s)
}

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
