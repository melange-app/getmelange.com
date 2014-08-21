package main

import (
	"encoding/json"
	"net/http"

	"github.com/airdispatch/go-pressure"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
)

type RedirectView struct {
	Location string
}

func (r *RedirectView) WriteBody(w http.ResponseWriter) {}

func (r *RedirectView) StatusCode() int {
	return 302
}
func (r *RedirectView) ContentLength() int  { return 0 }
func (r *RedirectView) ContentType() string { return "" }
func (r *RedirectView) Headers() pressure.ViewHeaders {
	a := make(pressure.ViewHeaders)
	a["Location"] = r.Location
	return a
}

type SessionView struct {
	Request *http.Request
	Session *sessions.Session
	pressure.View
}

func (r *SessionView) AddCookies(w http.ResponseWriter) {
	defer context.Clear(r.Request)

	r.Session.Save(r.Request, w)
}

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

func (j *JSONView) WriteBody(w http.ResponseWriter) {
	w.Write(j.getCache())
}

func (j *JSONView) StatusCode() int               { return 200 }
func (j *JSONView) ContentLength() int            { return len(j.getCache()) }
func (j *JSONView) ContentType() string           { return "application/json" }
func (j *JSONView) Headers() pressure.ViewHeaders { return nil }

type APIView struct {
	Method string
	pressure.View
}

func (a *APIView) WriteBody(r http.ResponseWriter) {
	if a.Method == "OPTIONS" {
		return
	}
	a.View.WriteBody(r)
}

func (a *APIView) StatusCode() int {
	if a.Method == "OPTIONS" {
		return 200
	}
	return a.View.StatusCode()
}

func (a *APIView) ContentLength() int {
	if a.Method == "OPTIONS" {
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
