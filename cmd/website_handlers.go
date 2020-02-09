package main

import (
	"fmt"
	"net/url"
	"strings"
)

func (app *application) Websites() {
	if strings.TrimSpace(app.flags.Website) != "" {
		app.load_website(app.flags.Website)
	}
	if strings.TrimSpace(app.flags.Websites) != "" {
		app.load_websites(app.flags.Websites)
	}
	if app.flags.List {
		app.list_websites()
	}
}

func (app *application) load_website(website string) {
	if !IsUrl(website) {
		app.errorLog.Println(website, "is not a valid url")
		return
	}
	h, _ := url.Parse(website)
	_, err := app.websites.Insert(h.Hostname())
	if err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) load_websites(path string) {
	websites, err := ReadFile(path)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	for _, website := range websites {
		if !IsUrl(website) {
			app.errorLog.Println(website, "is not a valid url")
			continue
		}
		h, _ := url.Parse(website)
		_, err := app.websites.Insert(h.Hostname())
		if err != nil {
			app.errorLog.Println(err)
		}
	}
}

func (app *application) list_websites() {
	w, err := app.websites.List()
	if err != nil {
		app.errorLog.Println(err)
	}
	for _, website := range w {
		fmt.Println(website.ID, website.Title)
	}
}
