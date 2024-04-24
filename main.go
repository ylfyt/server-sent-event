package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

type Sse struct {
	w io.Writer
}
func (me *Sse) Send(data []byte) error {
	_, err := me.w.Write(data)
	return err
}

func main() {
	r := chi.NewRouter()
	r.Use(cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
	}).Handler)

	r.Get("/sse", func(w http.ResponseWriter, r *http.Request) {
		// Set headers for Server-Sent Events
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		// Send initial event
		fmt.Fprintf(w, "data: Initial event\n\n")
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming not supported", http.StatusInternalServerError)
			return
		}

		// Simulate real-time updates
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-r.Context().Done():
				fmt.Println("Stop")
				return
			case <-ticker.C:
				fmt.Fprintf(w, "data: %s\n\n", "ok")
				flusher.Flush()
			}
		}
	})

	fmt.Println("Listening on port", 8080)
	http.ListenAndServe(":8080", r)
}
