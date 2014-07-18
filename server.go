package main

import (
	"fmt"
	"github.com/airdispatch/go-pressure"
	"labix.org/v2/mgo"
	"os"
)

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

	// Create Server and Necessary Engines
	theServer := pressure.CreateServer(":"+port, true)

	session, err := ConnectToDB()
	if err != nil {
		panic(err)
	}

	db := session.DB("get-melange-site")

	// Register API URLS
	theServer.RegisterURL(
		pressure.NewURLRoute("^/api/updates", &UpdateController{}),
		pressure.NewURLRoute("^/api/trackers", &TrackerController{
			Collection: db.C("trackers"),
		}),
		pressure.NewURLRoute("^/api/servers", &ServerController{
			Collection: db.C("servers"),
		}),
		pressure.NewURLRoute("^/api/applications", &ApplicationController{}),
	)

	// Register Application URLS
	theServer.RegisterURL(
		pressure.NewStaticFileRoute("^/public/", "app"),
		NewFileRoute("^/", "app/index.html"),
	)

	// Start the Server
	theServer.RunServer()
}
