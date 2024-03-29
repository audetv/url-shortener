package shorturl

import (
	"math/rand"
)

const letterBytes = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type ShortUrl string

// New создает строку из набора символов letterBytes с заданной длиной len
func New(len int) *ShortUrl {
	shortUrlLen := len
	shortUrl := randStringBytes(shortUrlLen)
	return (*ShortUrl)(&shortUrl)
}

func Parse(s string) *ShortUrl {
	// TODO здесь можно сделать проверку, что в строке содержаться символы из letterBytes
	return (*ShortUrl)(&s)
}

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
