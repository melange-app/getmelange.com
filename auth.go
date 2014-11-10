package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/airdispatch/go-pressure"
	"github.com/golang/oauth2"
	"github.com/google/go-github/github"
	"github.com/gorilla/sessions"
)

var config = &oauth2.Options{
	ClientID:     GithubId,
	ClientSecret: GithubSecret,
	Scopes:       []string{"user:email"},
	RedirectURL:  "http://www.getmelange.com/developer/login",
}

func GithubFromTransport(t *oauth2.Transport) *github.Client {
	return github.NewClient(&http.Client{
		Transport: t,
	})
}

func GetConfig() (*oauth2.Flow, error) {
	return oauth2.New(
		oauth2.Client(GithubId, GithubSecret),
		oauth2.Endpoint(
			"https://github.com/login/oauth/autorize",
			"https://github.com/login/oauth/access_token",
		),
		oauth2.RedirectURL("http://www.getmelange.com/developer/login"),
		oauth2.Scope("user:email"),
	)
}

func CreateClient(location string) (*github.Client, *oauth2.Transport) {
	config.RedirectURL = fmt.Sprintf("%s/developer/login", location)

	c, err := GetConfig()
	if err != nil {
		fmt.Println("Error getting config", err)
	}

	t := c.NewTransport()

	return github.NewClient(&http.Client{
		Transport: t,
	}), t
}

func PutTokenInSession(t *oauth2.Token, s *sessions.Session) {
	s.Values["access_token"] = t.AccessToken
	s.Values["refresh_token"] = t.RefreshToken
	s.Values["expiry"] = t.Expiry.Unix()
}

func GetTokenFromSession(s *sessions.Session) *oauth2.Token {
	return &oauth2.Token{
		AccessToken:  s.Values["access_token"].(string),
		RefreshToken: s.Values["refresh_token"].(string),
		Expiry:       time.Unix(s.Values["expiry"].(int64), 0),
	}
}

type LogoutController struct{}

func (c *LogoutController) GetResponse(p *pressure.Request, l *pressure.Logger) (pressure.View, *pressure.HTTPError) {
	session, err := store.Get(p.Request, "tokens")
	if err != nil {
		fmt.Println("Error getting session", err)
	}

	session.Values["authenticated"] = false
	session.Values["access_token"] = ""
	session.Values["refresh_token"] = ""
	session.Values["expiry"] = 0

	return &SessionView{
		Request: p.Request,
		Session: session,
		View: &RedirectView{
			Location: "/developer",
		},
	}, nil
}

type LoginController struct{}

func (c *LoginController) GetResponse(p *pressure.Request, l *pressure.Logger) (pressure.View, *pressure.HTTPError) {
	config, err := GetConfig()
	if err != nil {
		fmt.Println("Couldn't get config", err)
		return nil, &pressure.HTTPError{
			Code: 500,
			Text: err.Error(),
		}
	}

	t, err := config.NewTransportFromCode(p.Form["code"][0])
	if err != nil {
		fmt.Println("Couldn't exchange token", err)
		return nil, &pressure.HTTPError{
			Code: 500,
			Text: err.Error(),
		}
	}

	session, err := store.Get(p.Request, "tokens")
	if err != nil {
		fmt.Println("Error getting session", err)
		return nil, &pressure.HTTPError{
			Code: 500,
			Text: err.Error(),
		}
	}

	PutTokenInSession(t.Token(), session)
	session.Values["authenticated"] = true

	return &SessionView{
		Request: p.Request,
		Session: session,
		View: &RedirectView{
			Location: "/developer",
		},
	}, nil
}
