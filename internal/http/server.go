package http

import (
	"fmt"
	"net/http"
)

func StartHealth(port string) {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})
	fmt.Printf("Health server listening on :%s\n", port)
	http.ListenAndServe(":"+port, nil)
}
