package main

import (
	pressure "github.com/airdispatch/go-pressure"
	"labix.org/v2/mgo"
)

// Squirrel.Mac Update Response:
//
// {
//     "url": "http://mycompany.com/myapp/releases/myrelease",
//     "name": "My Release Name",
//     "notes": "Theses are some release notes innit",
//     "pub_date": "2013-09-18T12:29:53+01:00",
// }

type UpdateController struct{}

func (u *UpdateController) GetResponse(
	p *pressure.Request,
	l *pressure.Logger,
) (pressure.View, *pressure.HTTPError) {
	return nil, &pressure.HTTPError{
		Code: 204,
		Text: "",
	}
}

type Provider struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	Url           string `json:"url"`
	Fingerprint   string `json:"fingerprint"`
	EncryptionKey string `bson:"encryption_key" json:"encryption_key"`
	Proof         string `json:"proof"`
	Alias         string `json:"alias"`
	Users         int    `json:"users"`
}

type TrackerController struct {
	Collection *mgo.Collection
}

func (u *TrackerController) GetResponse(
	p *pressure.Request,
	l *pressure.Logger,
) (pressure.View, *pressure.HTTPError) {
	var trackers []*Provider
	err := u.Collection.Find(nil).Sort("-users").All(&trackers)
	if err != nil {
		return nil, &pressure.HTTPError{500, "500: Error loading trackers"}
	}

	return &APIView{
		Method: p.Method,
		View: &JSONView{
			Content: trackers,
		},
	}, nil
}

type ServerController struct {
	Collection *mgo.Collection
}

func (u *ServerController) GetResponse(
	p *pressure.Request,
	l *pressure.Logger,
) (pressure.View, *pressure.HTTPError) {
	var servers []*Provider
	err := u.Collection.Find(nil).Sort("-users").All(&servers)
	if err != nil {
		return nil, &pressure.HTTPError{500, "500: Error loading servers"}
	}

	return &APIView{
		Method: p.Method,
		View: &JSONView{
			Content: servers,
		},
	}, nil
}

type ApplicationController struct{}

func (u *ApplicationController) GetResponse(
	p *pressure.Request,
	l *pressure.Logger,
) (pressure.View, *pressure.HTTPError) {
	return nil, &pressure.HTTPError{
		Code: 204,
		Text: "",
	}
}
