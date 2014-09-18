package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/airdispatch/go-pressure"
	"github.com/google/go-github/github"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte(os.Getenv("MLGSECRET")))

type DeveloperIndex struct {
	Engine *pressure.TemplateEngine
}

func (c *DeveloperIndex) GetResponse(p *pressure.Request, l *pressure.Logger) (pressure.View, *pressure.HTTPError) {
	session, err := store.Get(p.Request, "tokens")
	if err != nil {
		fmt.Println("Error getting session", err)
	}

	obj, ok := session.Values["authenticated"]
	auth := ok && obj.(bool)

	var u *github.User
	if auth {
		validToken := GetTokenFromSession(session)

		config, _ := GetConfig()
		t := config.NewTransport()
		t.SetToken(validToken)

		gh := GithubFromTransport(t)

		u, _, err = gh.Users.Get("")
		if err != nil {
			fmt.Println("Error getting github user.", err)

			return nil, &pressure.HTTPError{
				Code: 500,
				Text: err.Error(),
			}
		}
	}

	return c.Engine.NewTemplateView("developer/index.html", map[string]interface{}{
		"clientId":      GithubId,
		"authenticated": auth,
		"user":          u,
	}), nil
}

func clientOrError(req *http.Request) (*github.Client, *pressure.HTTPError) {
	session, err := store.Get(req, "tokens")
	if err != nil {
		fmt.Println("Error getting session", err)
	}

	obj, ok := session.Values["authenticated"]
	auth := ok && obj.(bool)

	if !auth {
		return nil, &pressure.HTTPError{
			Code: 403,
			Text: "Must be logged in.",
		}
	}

	validToken := GetTokenFromSession(session)

	config, _ := GetConfig()
	t := config.NewTransport()
	t.SetToken(validToken)

	return GithubFromTransport(t), nil
}
