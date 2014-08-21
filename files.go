package main

import (
	"github.com/airdispatch/go-pressure"
  "fmt"
	"os"
)

type FileRoute struct {
  File string
  pressure.URLRoute
}

func NewFileRoute(pattern string, file string) *FileRoute {
  u := pressure.NewURLRoute(pattern, nil)
  return &FileRoute {
    File: file,
    URLRoute: u,
  }
}

func (f *FileRoute) GetMatch(path string) (*pressure.RouteMatch, bool) {
  _, b := f.URLRoute.GetMatch(path)
  if !b {
    fmt.Println(b)
    return nil, false
  }

  fi, err := os.Open(f.File)
  if err != nil {
    fmt.Println(err)
    return nil, false
  }

  return &pressure.RouteMatch {
    Path: path,
    Controller: pressure.StaticFileViewController{
      File: fi,
      Name: f.File,
    },
  }, true
}
