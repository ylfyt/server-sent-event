package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

type Sse struct {
	w       io.Writer
	id      int64
	flusher http.Flusher
	notifiy chan struct{}
}

func (me *Sse) Send(data []byte) error {
	defer func() {
		me.flusher.Flush()
	}()
	_, err := me.w.Write([]byte("data: "))
	if err != nil {
		return err
	}
	_, err = me.w.Write(data)
	if err != nil {
		return err
	}
	_, err = me.w.Write([]byte("\n\n"))
	if err != nil {
		return err
	}
	return nil
}
func (me *Sse) Done() <-chan struct{} {
	return me.notifiy
}

func timerHandler(subs map[int64]*Sse) {
	for {
		now := time.Now().UnixMilli()
		for _, sse := range subs {
			sse.Send([]byte(strconv.Itoa(int(now))))
		}
		time.Sleep(1 * time.Second)
	}
}

func main() {
	subs := make(map[int64]*Sse)

	r := chi.NewRouter()
	r.Use(cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
	}).Handler)

	r.Get("/sse", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming not supported", http.StatusInternalServerError)
			return
		}

		id := time.Now().UnixMilli()
		sse := &Sse{
			w:       w,
			id:      id,
			flusher: flusher,
			notifiy: make(chan struct{}),
		}
		subs[id] = sse

		fmt.Println("New SSE Connection", id)

		defer func() {
			fmt.Println("Done", id)
			delete(subs, id)
		}()
		for {
			select {
			case <-r.Context().Done():
				return
			case <-sse.Done():
				return
			}
		}
	})
	go timerHandler(subs)

	fmt.Println("Listening on port", 8080)
	http.ListenAndServe(":8080", r)
}
