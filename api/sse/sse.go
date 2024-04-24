package sse

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Sse struct {
	w       io.Writer
	id      int64
	flusher http.Flusher
	notifiy chan struct{}
}
type SseData struct {
	Id   int64
	Path string
	Data any
}

func (me *Sse) write(data []byte) error {
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
func (me *Sse) Done() {
	me.notifiy <- struct{}{}
}

func (me *Sse) Closed() <-chan struct{} {
	return me.notifiy
}

func (me *Sse) Send(path string, data any) error {
	newData := SseData{
		Id:   me.id,
		Path: path,
		Data: data,
	}
	bytes, err := json.Marshal(newData)
	if err != nil {
		return err
	}
	err = me.write(bytes)
	return err
}

func New(w io.Writer, id int64) (*Sse, error) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		return nil, fmt.Errorf("streaming not supported")
	}
	sse := &Sse{
		w:       w,
		id:      id,
		flusher: flusher,
		notifiy: make(chan struct{}),
	}
	return sse, nil
}
