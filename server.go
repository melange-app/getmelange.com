package main

import (
	"github.com/airdispatch/go-pressure"
  "os"
)

func main() {
  port := os.Getenv("PORT")
  if port == "" {
    port = "9000"
  }

	// Create Server and Necessary Engines
	theServer := pressure.CreateServer(":" + port, true)

  // Register API URLS
  theServer.RegisterURL(
    pressure.NewURLRoute("^/api/updates", nil),
    pressure.NewURLRoute("^/api/trackers", nil),
    pressure.NewURLRoute("^/api/servers", nil),
    pressure.NewURLRoute("^/api/applications", nil),
  )

  // Register Application URLS
	theServer.RegisterURL(
		pressure.NewStaticFileRoute("^/public/", "app"),
		NewFileRoute("^/", "app/index.html"),
	)

	// Start the Server
	theServer.RunServer()
}
