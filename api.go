package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"airdispat.ch/tracker"
	pressure "github.com/airdispatch/go-pressure"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

// Squirrel.Mac Update Response:
//
// {
//     "url": "http://mycompany.com/myapp/releases/myrelease",
//     "name": "My Release Name",
//     "notes": "Theses are some release notes innit",
//     "pub_date": "2013-09-18T12:29:53+01:00",
// }

type Controller404 struct{}

func (u *Controller404) GetResponse(
	p *pressure.Request,
	l *pressure.Logger,
) (pressure.View, *pressure.HTTPError) {
	return nil, pressure.HTTPErrorNotFound
}

type Version string

func (v Version) After(a Version) bool {
	vComp := strings.Split(string(v), ".")
	aComp := strings.Split(string(a), ".")
	if len(vComp) != 3 || len(aComp) != 3 {
		fmt.Println("Not valid versions", v, a)
		return false
	}

	majorV, errv := strconv.Atoi(vComp[0])
	majorA, erra := strconv.Atoi(aComp[0])

	if errv != nil || erra != nil {
		return false
	}
	if majorV != majorA {
		return (majorV > majorA)
	}

	minorV, errv := strconv.Atoi(vComp[1])
	minorA, erra := strconv.Atoi(aComp[1])

	if errv != nil || erra != nil {
		return false
	}
	if minorV != minorA {
		return (minorV > minorA)
	}

	buildV, errv := strconv.Atoi(vComp[2])
	buildA, erra := strconv.Atoi(aComp[2])

	if errv != nil || erra != nil {
		return false
	}
	if buildV != buildA {
		return (buildV > buildA)
	}

	return false
}

type Update struct {
	Version   string            `json:"version"`
	Changelog string            `json:"changelog"`
	Download  string            `json:"download" bson:"-"`
	Supports  map[string]string `json:"-"`
}

type UpdateController struct {
	Collection *mgo.Collection
}

func (u *UpdateController) GetResponse(
	p *pressure.Request,
	l *pressure.Logger,
) (pressure.View, *pressure.HTTPError) {
	var updates []*Update
	err := u.Collection.Find(nil).Sort("-version").All(&updates)
	if err != nil {
		return nil, &pressure.HTTPError{500, "500: Error loading updates"}
	}

	if len(updates) == 0 ||
		!Version(updates[0].Version).After(Version(p.URL["version"])) ||
		updates[0].Supports[p.URL["platform"]] == "" {
		return nil, &pressure.HTTPError{
			Code: 422,
			Text: "No updates at this time.",
		}
	}

	updates[0].Download = updates[0].Supports[p.URL["platform"]]

	return &APIView{
		Method: p.Method,
		View: &JSONView{
			Content: updates[0],
		},
	}, nil
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

func IsDebug(req *http.Request) bool {
	q := req.URL.RawQuery

	values, err := url.ParseQuery(q)
	if err != nil {
		return false
	}

	obj := values.Get("debug")
	return obj != ""
}

type TrackerController struct {
	Collection *mgo.Collection
}

func (u *TrackerController) GetResponse(
	p *pressure.Request,
	l *pressure.Logger,
) (pressure.View, *pressure.HTTPError) {
	var trackers []*Provider
	err := u.Collection.Find(map[string]interface{}{
		"debug": IsDebug(p.Request),
	}).Sort("-users").All(&trackers)
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
	err := u.Collection.Find(map[string]interface{}{
		"debug": IsDebug(p.Request),
	}).Sort("-users").All(&servers)
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

type ApplicationController struct {
	Collection *mgo.Collection
}

func (u *ApplicationController) GetResponse(
	p *pressure.Request,
	l *pressure.Logger,
) (pressure.View, *pressure.HTTPError) {
	// Load Query Params
	var find map[string]interface{}
	query, err := url.ParseQuery(p.Request.URL.RawQuery)
	if err != nil {
		fmt.Println("Error parsing query", err)
		return nil, internalError
	}

	// Get Query from URL
	if query.Get("q") != "" {
		find = map[string]interface{}{
			"name": bson.RegEx{
				Pattern: fmt.Sprintf(`.*(%s).*`, query),
			},
		}
	}

	// Find Apps
	apps := make([]*App, 0)
	err = u.Collection.Find(find).All(&apps)
	if err != nil {
		return nil, &pressure.HTTPError{500, "500: Error loading servers"}
	}

	return &APIView{
		Method: p.Method,
		View: &JSONView{
			Content: apps,
		},
	}, nil
}

type ResolverController struct{}

func (u *ResolverController) GetResponse(
	p *pressure.Request,
	l *pressure.Logger,
) (pressure.View, *pressure.HTTPError) {
	return &pressure.BasicView{
		Status: 200,
		Text: tracker.GetTrackingServerLocationFromURL(
			p.URL["url"],
		),
	}, nil
}
