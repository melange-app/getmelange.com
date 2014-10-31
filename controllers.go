package main

import (
	pressure "github.com/airdispatch/go-pressure"
	"labix.org/v2/mgo"
)

type ProviderController struct {
	Trackers *mgo.Collection
	Servers  *mgo.Collection
	Engine   *pressure.TemplateEngine
}

func (c *ProviderController) GetResponse(p *pressure.Request, l *pressure.Logger) (pressure.View, *pressure.HTTPError) {
	var trackers []*Provider
	err := c.Trackers.Find(map[string]interface{}{
		"debug": IsDebug(p.Request),
	}).Sort("-users").Limit(6).All(&trackers)
	if err != nil {
		return nil, &pressure.HTTPError{500, "500: Error loading trackers"}
	}

	var servers []*Provider
	err = c.Servers.Find(map[string]interface{}{
		"debug": IsDebug(p.Request),
	}).Sort("-users").Limit(6).All(&servers)
	if err != nil {
		return nil, &pressure.HTTPError{500, "500: Error loading trackers"}
	}

	return c.Engine.NewTemplateView("pages/providers.html", map[string]interface{}{
		"trackers": trackers,
		"servers":  servers,
	}), nil
}

type AppController struct {
	Collection *mgo.Collection
	Engine     *pressure.TemplateEngine
}

func (c *AppController) GetResponse(p *pressure.Request, l *pressure.Logger) (pressure.View, *pressure.HTTPError) {
	var apps []*App
	err := c.Collection.Find(map[string]interface{}{
		"debug": IsDebug(p.Request),
	}).Sort("-users").Limit(6).All(&apps)
	if err != nil {
		return nil, &pressure.HTTPError{500, "500: Error loading apps"}
	}

	return c.Engine.NewTemplateView("pages/apps.html", map[string]interface{}{
		"apps":     apps,
		"mostUsed": apps,
	}), nil
}
