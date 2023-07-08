package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%s/health", os.Getenv("PORT")))
	if err != nil {
		os.Exit(1)
	}
	defer resp.Body.Close()
}
