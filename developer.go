package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	pressure "github.com/airdispatch/go-pressure"
	"github.com/google/go-github/github"
	"github.com/gorilla/sessions"
	"github.com/russross/blackfriday"
)

var store = sessions.NewCookieStore([]byte(os.Getenv("MLGSECRET")))

type tocEntity struct {
	Title string
	Link  string
}

type DeveloperRenderer struct {
	blackfriday.Renderer

	count int
	TOC   []tocEntity
}

func CreateRenderer() *DeveloperRenderer {
	return &DeveloperRenderer{
		Renderer: blackfriday.HtmlRenderer(
			blackfriday.HTML_USE_XHTML|
				blackfriday.HTML_USE_SMARTYPANTS|
				blackfriday.HTML_SMARTYPANTS_FRACTIONS|
				blackfriday.HTML_SMARTYPANTS_LATEX_DASHES,
			"", "",
		),
	}
}

func (d *DeveloperRenderer) Render(input []byte) ([]byte, []tocEntity) {
	output := blackfriday.Markdown(
		input,
		d,
		blackfriday.EXTENSION_NO_INTRA_EMPHASIS|
			blackfriday.EXTENSION_TABLES|
			blackfriday.EXTENSION_FENCED_CODE|
			blackfriday.EXTENSION_AUTOLINK|
			blackfriday.EXTENSION_STRIKETHROUGH|
			blackfriday.EXTENSION_SPACE_HEADERS|
			blackfriday.EXTENSION_HEADER_IDS,
	)

	return output, d.TOC
}

func (d *DeveloperRenderer) Header(out *bytes.Buffer, text func() bool, level int, id string) {
	if level == 1 {
		// Assign an id that we can believe in.
		d.count++
		id = fmt.Sprintf("toc_%d", d.count)
	}

	var title string
	d.Renderer.Header(out, func() bool {
		// Extract the title
		current := out.Len()
		output := text()
		title = string(out.Bytes()[current:])
		return output
	}, level, id)

	if level == 1 {
		// Populate the TOC
		if d.TOC == nil {
			d.TOC = make([]tocEntity, 0)
		}

		d.TOC = append(d.TOC, tocEntity{
			Title: title,
			Link:  fmt.Sprintf("#%s", id),
		})
	}
}

type DeveloperDocument struct {
	Engine *pressure.TemplateEngine
}

func (c *DeveloperDocument) GetResponse(p *pressure.Request, l *pressure.Logger) (pressure.View, *pressure.HTTPError) {
	// documents/getting-started-building-applications.md
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	currentDocument := p.URL["document"] + ".md"
	documentPath := path.Join(cwd, "documents", currentDocument)

	data, err := os.Open(documentPath)
	defer data.Close()
	if err != nil {
		return nil, &pressure.HTTPError{
			Code: 404,
			Text: "Can't find that document.",
		}
	}

	obj, err := ioutil.ReadAll(data)
	if err != nil {
		panic(err)
	}

	components := strings.Split(string(obj), "\n---\n")
	if len(components) < 2 {
		panic("Cannot find title")
	}

	title := components[0]
	body := []byte(components[1])

	renderer := CreateRenderer()
	output, toc := renderer.Render(body)

	return c.Engine.NewTemplateView("developer/document.html", map[string]interface{}{
		"Title":    title,
		"Body":     template.HTML(string(output)),
		"Contents": toc,
	}), nil
}

type DeveloperAPI struct {
	Engine *pressure.TemplateEngine
}

func (c *DeveloperAPI) GetResponse(p *pressure.Request, l *pressure.Logger) (pressure.View, *pressure.HTTPError) {
	fmt.Println(p.URL["document"])
	return c.Engine.NewTemplateView("developer/api.html", map[string]interface{}{}), nil
}

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
