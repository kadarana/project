package api

import (
	"encoding/json"
	"net/http"
	"project/internal/cache"
	"strings"
)

func NewServer(c *cache.Cache) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/order/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/order/")

		order, ok := c.Get(id)
		if !ok {
			http.Error(w, "order not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(order)
	})

	return mux
}
