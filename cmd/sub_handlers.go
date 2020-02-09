package main

import (
	"fmt"
	"strings"

	"github.com/pavel1337/webfingerprint/pkg/crawler"
)

func (app *application) Subs() {
	if app.flags.GenerateSubs {
		app.generetaSubs(app.flags.Depth)
	}
	if app.flags.ListAllSubs {
		app.listSubs()
	}
	if strings.TrimSpace(app.flags.ListSubs) != "" {
		app.listSubsByTitle(app.flags.ListSubs)
	}
}

func (app *application) generetaSubs(depth int) {
	w, err := app.websites.List()
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	for _, website := range w {
		intialurl := "https://" + website.Title
		subs, err := app.subs.GetByWebsiteId(website.ID)
		if err != nil {
			app.errorLog.Println(err)
			continue
		}
		if len(subs) == 0 {
			_, err := app.subs.Insert(intialurl, website.ID)
			if err != nil {
				app.errorLog.Println(err)
			}
		}
		if len(subs) >= depth {
			continue
		}
		app.infoLog.Println("generating sub pages for", website.Title, "with depth", depth)
		newsubs := crawler.Crawler(intialurl, depth)
		for _, newsub := range newsubs {
			_, err := app.subs.Insert(newsub, website.ID)
			if err != nil {
				app.errorLog.Println(err)
			}
		}
	}
}

func (app *application) listSubsByTitle(website string) {
	w, err := app.websites.GetByTitle(website)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	subs, err := app.subs.GetByWebsiteId(w.ID)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	for _, sub := range subs {
		fmt.Println(sub.ID, sub.Title)
	}
}

func (app *application) listSubs() {
	subs, err := app.subs.List()
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	for _, sub := range subs {
		fmt.Println(sub.ID, sub.Title)
	}
}
