package main

import (
	"fmt"
	"os"
	"path/filepath"

	pressure "github.com/airdispatch/go-pressure"
	"labix.org/v2/mgo"
)

var GithubId = os.Getenv("GITHUB_ID")
var GithubSecret = os.Getenv("GITHUB_SECRET")

func ConnectToDB() (*mgo.Session, error) {
	// Load From Env
	user := os.Getenv("MONGO_USERNAME")
	pass := os.Getenv("MONGO_PASSWORD")
	url := os.Getenv("MONGO_URL")
	connect := fmt.Sprintf("mongodb://%s:%s@%s", user, pass, url)
	return mgo.Dial(connect)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}

	temp_wd, _ := os.Getwd()
	template_dir := filepath.Join(temp_wd, "templates")

	// Create Server and Necessary Engines
	theServer := pressure.CreateServer(":"+port, true)

	session, err := ConnectToDB()
	if err != nil {
		panic(err)
	}

	db := session.DB("get-melange-site")

	t := theServer.CreateTemplateEngine(template_dir, "base.html")

	// Register Golang Packages
	theServer.RegisterURL(
		// Application Packages
		pressure.NewURLRoute("^/app$", &GolangFetchController{"app", "melange-core"}),
		pressure.NewURLRoute("^/app/controllers$", &GolangFetchController{"app", "melange-backend"}),
		pressure.NewURLRoute("^/app/framework$", &GolangFetchController{"app", "melange-backend"}),
		pressure.NewURLRoute("^/app/models$", &GolangFetchController{"app", "melange-backend"}),
		pressure.NewURLRoute("^/app/packaging$", &GolangFetchController{"app", "melange-backend"}),
		// DAP Packages
		pressure.NewURLRoute("^/dap$", &GolangFetchController{"dap", "dap"}),
		pressure.NewURLRoute("^/dap/wire$", &GolangFetchController{"dap", "dap"}),
		// Dispatcher Packages
		pressure.NewURLRoute("^/dispatcher$", &GolangFetchController{"dispatcher", "dispatcher"}),
		pressure.NewURLRoute("^/dispatcher/dispatcher$", &GolangFetchController{"dispatcher", "dispatcher"}),
		// Router Packages
		pressure.NewURLRoute("^/router$", &GolangFetchController{"router", "melange-router"}),
		// Server Packages
		// pressure.NewURLRoute("^/server$", &GolangFetchController{"server", "melange-backend"}),
		// Tracker Packages
		pressure.NewURLRoute("^/tracker$", &GolangFetchController{"tracker", "tracker"}),
		pressure.NewURLRoute("^/tracker/tracker$", &GolangFetchController{"tracker", "tracker"}),
		// Updater Packages
		pressure.NewURLRoute("^/updater$", &GolangFetchController{"updater", "melange-updater"}),
		pressure.NewURLRoute("^/updater/updater$", &GolangFetchController{"updater", "melange-updater"}),
	)

	// Register API URLS
	theServer.RegisterURL(
		pressure.NewURLRoute("^/favicon.ico", &Controller404{}),
		pressure.NewURLRoute("^/api/trackers", &TrackerController{
			Collection: db.C("trackers"),
		}),
		pressure.NewURLRoute("^/api/servers", &ServerController{
			Collection: db.C("servers"),
		}),
		pressure.NewURLRoute(`^/api/updates/(?P<version>[\w\.]*)/(?P<platform>\w*)`, &UpdateController{
			Collection: db.C("updates"),
		}),
		pressure.NewURLRoute("^/api/applications", &ApplicationController{
			Collection: db.C("apps"),
		}),
	)

	// Register Application URLS
	theServer.RegisterURL(
		pressure.NewStaticFileRoute("^/public/", "static"),

		// Developer URLS
		pressure.NewURLRoute("^/developer/login", &LoginController{}),
		pressure.NewURLRoute("^/developer/logout", &LogoutController{}),

		pressure.NewURLRoute("^/developer/view/(?P<document>.*)", &DeveloperDocument{
			Engine: t,
		}),

		pressure.NewURLRoute("^/developer/api/(?P<document>.*)", &DeveloperAPI{
			Engine: t,
		}),

		pressure.NewURLRoute("^/developer/publish/app", &AddApplicationController{
			Collection: db.C("apps"),
			Engine:     t,
		}),

		pressure.NewURLRoute("^/developer/publish/tracker", &DeveloperIndex{
			Engine: t,
		}),
		pressure.NewURLRoute("^/developer/publish/server", &DeveloperIndex{
			Engine: t,
		}),

		pressure.NewURLRoute("^/developer", &DeveloperIndex{
			Engine: t,
		}),

		// Regular URLS
		pressure.NewURLRoute("^/providers", &ProviderController{
			Trackers: db.C("trackers"),
			Servers:  db.C("servers"),
			Engine:   t,
		}),

		pressure.NewURLRoute("^/apps", &AppController{
			Collection: db.C("apps"),
			Engine:     t,
		}),

		pressure.NewURLRoute("^/", &BasicTemplate{
			Engine: t,
			Name:   "index.html",
		}),
		// NewFileRoute("^/", "app/index.html"),
	)

	// Start the Server
	theServer.RunServer()
}

type BasicTemplate struct {
	Engine *pressure.TemplateEngine
	Name   string
}

func (b *BasicTemplate) GetResponse(*pressure.Request, *pressure.Logger) (pressure.View, *pressure.HTTPError) {
	return b.Engine.NewTemplateView(b.Name, nil), nil
}

type GolangFetchController struct {
	prefixName  string
	packageName string
}

func (c *GolangFetchController) GetResponse(p *pressure.Request, l *pressure.Logger) (pressure.View, *pressure.HTTPError) {
	return pressure.NewHTMLView(
		`<html>
			<head>
				<meta name="go-import" content="getmelange.com/` + c.prefixName + ` git https://github.com/melange-app/` + c.packageName + `">
			</head>
		</html>`), nil
}
