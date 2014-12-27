package main

import (
	"fmt"
	"net/http"
	"time"

	pressure "github.com/airdispatch/go-pressure"
	"github.com/google/go-github/github"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
)

var config = &oauth2.Config{
	ClientID:     GithubId,
	ClientSecret: GithubSecret,
	Scopes:       []string{"user:email"},
	RedirectURL:  "http://www.getmelange.com/developer/login",
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://github.com/login/oauth/autorize",
		TokenURL: "https://github.com/login/oauth/access_token",
	},
}

func GithubFromTransport(t *http.Client) *github.Client {
	return github.NewClient(t)
}

func CreateClient(location string, token *oauth2.Token) (*github.Client, *http.Client) {
	config.RedirectURL = fmt.Sprintf("%s/developer/login", location)

	transport := config.Client(nil, token)

	return github.NewClient(
		transport,
	), transport
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
	t, err := config.Exchange(nil, p.Form["code"][0])
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

	PutTokenInSession(t, session)
	session.Values["authenticated"] = true

	return &SessionView{
		Request: p.Request,
		Session: session,
		View: &RedirectView{
			Location: "/developer",
		},
	}, nil
}
