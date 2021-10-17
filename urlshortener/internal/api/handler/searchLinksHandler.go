package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (rt *Router) SearchLinks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	q := r.URL.Query().Get("q")

	ch, err := rt.links.SearchLinks(r.Context(), q)
	if err != nil {
		http.Error(w, "error when searching", http.StatusInternalServerError)
		return
	}

	enc := json.NewEncoder(w)
	first := true
	fmt.Printf("[")
	defer fmt.Printf("]\r\n")

	for {
		select {
		case <-r.Context().Done():
			return
		case l, ok := <-ch:
			if !ok {
				return
			}
			if first {
				first = false
			} else {
				fmt.Fprintf(w, ",")
			}
			_ = enc.Encode(
				Link{
					Short:  l.Short,
					Origin: l.Origin,
				},
			)
			w.(http.Flusher).Flush()
		}
	}
}
