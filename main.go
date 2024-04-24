package main

import (
	"fmt"
	"go_sse/api/sse"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

var counter = 0

type Subs map[int64]*sse.Sse

func (me *Subs) Send(path string, data any) error {
	for _, sse := range *me {
		err := sse.Send("/timer", data)
		if err != nil {
			return err
		}
	}
	return nil
}

func increment(subs Subs) {
	counter++
	if counter >= 10 {
		counter = 0
		for _, sse := range subs {
			sse.Done()
		}
		return
	}
	subs.Send("/counter", counter)
}

func main() {
	subs := make(Subs)

	r := chi.NewRouter()
	r.Use(cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
	}).Handler)

	r.Post("/inc", func(w http.ResponseWriter, r *http.Request) {
		increment(subs)
	})

	r.Get("/sse", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		id := time.Now().UnixMilli()
		sse, err := sse.New(w, id)
		if err != nil {
			fmt.Println("Failed to initiate SSE", id, err)
			http.Error(w, "Failed to initiate SSE", http.StatusInternalServerError)
			return
		}

		fmt.Println("New SSE Connection", id)
		subs[id] = sse
		sse.Send("/counter", counter)

		defer func() {
			fmt.Println("Done", id)
			delete(subs, id)
		}()
		for {
			select {
			case <-r.Context().Done():
				return
			case <-sse.Closed():
				return
			}
		}
	})

	fmt.Println("Listening on port", 8080)
	http.ListenAndServe(":8080", r)
}
