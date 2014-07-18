package main

import (
	"encoding/json"
	"github.com/airdispatch/go-pressure"
	"io"
)

type JSONView struct {
	Content interface{}
	cache   []byte
}

func (j *JSONView) getCache() []byte {
	if j.cache != nil {
		return j.cache
	}
	bytes, err := json.Marshal(j.Content)
	if err != nil {
		panic("JSON Encoding " + err.Error())
	}
	j.cache = bytes
	return j.cache
}

func (j *JSONView) WriteBody(w io.Writer) {
	w.Write(j.getCache())
}

func (j *JSONView) StatusCode() int               { return 200 }
func (j *JSONView) ContentLength() int            { return len(j.getCache()) }
func (j *JSONView) ContentType() string           { return "application/json" }
func (j *JSONView) Headers() pressure.ViewHeaders { return nil }
