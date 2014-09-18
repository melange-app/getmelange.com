package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/airdispatch/go-pressure"
	"github.com/google/go-github/github"
	"labix.org/v2/mgo"
)

var internalError = &pressure.HTTPError{
	Code: 500,
	Text: "Internal service error.",
}

type App struct {
	Id          string
	Name        string
	Description string
	Username    string
	Debug       bool `json:"-"`
	Repository  string
}

type AddApplicationController struct {
	Engine     *pressure.TemplateEngine
	Collection *mgo.Collection
}

func (a *AddApplicationController) GetResponse(p *pressure.Request, l *pressure.Logger) (pressure.View, *pressure.HTTPError) {
	if p.Request.Method == "GET" {
		gh, httpErr := clientOrError(p.Request)
		if httpErr != nil {
			return nil, httpErr
		}

		repos := make(map[string]interface{})

		myRepos, _, err := gh.Repositories.List("", &github.RepositoryListOptions{
			Type: "public",
		})
		if err != nil {
			fmt.Println("Error getting github repos.", err)
		}

		repos[""] = myRepos

		orgs, _, err := gh.Organizations.List("", nil)
		if err != nil {
			fmt.Println("Couldn't get orgs", err)
		}

		for _, v := range orgs {
			orgRepos, _, err := gh.Repositories.ListByOrg(*v.Login, nil)
			if err != nil {
				fmt.Println("Couldn't get org repos", *v.Login, err)
			}

			repos[*v.Login] = orgRepos
		}

		return a.Engine.NewTemplateView("developer/application/new.html", map[string]interface{}{
			"repos": repos,
		}), nil
	} else if p.Request.Method == "POST" {
		// id := p.Form["application_id"][0]
		url := strings.Split(p.Form["application_url"][0], "/")
		id := fmt.Sprintf("com.github.%s.%s", url[0], url[1])

		gh, httpErr := clientOrError(p.Request)
		if httpErr != nil {
			return nil, httpErr
		}

		u, _, err := gh.Users.Get("")
		if err != nil {
			fmt.Println("Error getting github user.", err)

			return nil, &pressure.HTTPError{
				Code: 500,
				Text: err.Error(),
			}
		}

		if url[0] != *u.Login {
			orgs, _, err := gh.Organizations.List("", nil)
			if err != nil {
				fmt.Println("Couldn't get orgs", err)
			}

			match := false
			for _, v := range orgs {
				if *v.Login == url[0] {
					match = true
					break
				}
			}

			if !match {
				return nil, &pressure.HTTPError{
					Code: 403,
					Text: "You cannot add other people's repos.",
				}
			}
		}

		tags, _, err := gh.Repositories.ListTags(url[0], url[1], nil)

		latest := ""
		for _, v := range tags {
			if strings.HasPrefix(*v.Name, "v") {
				latest = *v.Name
				break
			}
		}

		if latest == "" {
			return nil, &pressure.HTTPError{
				Code: 400,
				Text: "You cannot publish an application without a release.",
			}
		}

		pkg := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/package.json", url[0], url[1], latest)
		fmt.Println("Fetching package.json from", pkg)

		res, err := http.Get(pkg)
		if err != nil {
			fmt.Println("Can't get package.json file", err)
			return nil, internalError
		}
		defer res.Body.Close()

		dec := json.NewDecoder(res.Body)

		app := &App{}
		err = dec.Decode(app)
		if err != nil {
			fmt.Println("Can't decode app.", err)
			return nil, internalError
		}

		if app.Id != id {
			return nil, &pressure.HTTPError{
				Code: 400,
				Text: fmt.Sprintf("Your repo's package.json id (%s) must match the real id (%s).", app.Id, id),
			}
		}

		// Repo is Valid - let's add it.

		app.Username = *u.Login
		app.Debug = false
		app.Repository = fmt.Sprintf("http://github.com/%s/%s", url[0], url[1])

		err = a.Collection.Insert(app)
		if err != nil {
			return nil, internalError
		}

		return &RedirectView{
			Location: "/developer",
		}, nil
	}

	return nil, &pressure.HTTPError{
		Code: 401,
		Text: "Method not allowed",
	}
}
