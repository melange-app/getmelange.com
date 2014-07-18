package main

import (
	"encoding/json"
	"github.com/airdispatch/go-pressure"
  "net/http"
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


type APIView struct {
	Request *http.Request
	pressure.View
}

func (a *APIView) WriteBody(r io.Writer) {
	if a.Request.Method == "OPTIONS" {
		return
	}
	a.View.WriteBody(r)
}

func (a *APIView) StatusCode() int {
	if a.Request.Method == "OPTIONS" {
		return 200
	}
	return a.View.StatusCode()
}

func (a *APIView) ContentLength() int {
	if a.Request.Method == "OPTIONS" {
		return 0
	}
	return a.View.ContentLength()
}

func (a *APIView) Headers() pressure.ViewHeaders {
	hdrs := a.View.Headers()
	if hdrs == nil {
		hdrs = make(pressure.ViewHeaders)
	}

	hdrs["Access-Control-Allow-Origin"] = "*"
	hdrs["Access-Control-Allow-Headers"] = "Content-Type"
	return hdrs
}
