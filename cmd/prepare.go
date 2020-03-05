package main

import (
	"fmt"
	"net/url"

	"github.com/pavel1337/wasm/pkg/crawler"
	"github.com/pavel1337/wasm/pkg/models"
)

func (app *application) Prepare() {
	if !IsEmpty(app.flags.Website) {
		app.loadWebsite(app.flags.Website)
	}
	if !IsEmpty(app.flags.Websites) {
		app.loadWebsites(app.flags.Websites)
	}
	if app.flags.List {
		app.listWebsites()
	}
	if app.flags.GenerateSubs {
		app.generetaSubs(app.flags.Depth)
	}
	if app.flags.ListAllSubs {
		app.listSubs()
	}
	if !IsEmpty(app.flags.ListSubs) {
		app.listSubsByHostname(app.flags.ListSubs)
	}
	if app.flags.ListSubsWithPcap {
		app.listSubsByProxy()
	}
}

func (app *application) loadWebsite(website string) {
	if !IsUrl(website) {
		app.errorLog.Println(website, "is not a valid url")
		return
	}
	h, _ := url.Parse(website)
	app.db.Create(&models.Website{Hostname: h.Hostname()})
}

func (app *application) loadWebsites(path string) {
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
		app.db.Create(&models.Website{Hostname: h.Hostname()})
	}
}

func (app *application) listWebsites() {
	websites := []models.Website{}
	app.db.Find(&websites)
	for _, website := range websites {
		fmt.Println(website.ID, website.Hostname)
	}
}

func (app *application) generetaSubs(depth int) {
	websites := []models.Website{}
	app.db.Find(&websites)
	for _, website := range websites {
		intialurl := "https://" + website.Hostname

		subs := []models.Sub{}
		app.db.Model(&website).Related(&subs)

		if len(subs) == 0 {
			sub := models.Sub{Link: intialurl, WebsiteID: website.ID}
			app.db.Create(&sub)
		}
		if len(subs) >= depth {
			continue
		}
		app.infoLog.Println("generating sub pages for", website.Hostname, "with depth", depth)
		newsubs := crawler.Crawler(intialurl, depth)
		for _, newsub := range newsubs {
			sub := models.Sub{Link: newsub, WebsiteID: website.ID}
			app.db.Create(&sub)
		}
	}
}

func (app *application) listSubs() {
	subs := []models.Sub{}
	app.db.Find(&subs)
	for _, sub := range subs {
		fmt.Println(sub.ID, sub.Link)
	}
}

func (app *application) listSubsByHostname(hostname string) {
	website := models.Website{}
	app.db.Where(&models.Website{Hostname: hostname}).First(&website)

	subs := []models.Sub{}
	app.db.Model(&website).Related(&subs)

	for _, sub := range subs {
		fmt.Println(sub.ID, sub.Link)
	}
}

func (app *application) listSubsByProxy() {
	subs := []models.Sub{}
	app.db.Find(&subs)

	fmt.Printf("%v\t%v\t%v\t%v\n", "clean", "out", "proxy", "url")
	for _, sub := range subs {
		proxies, err := models.DistinctProxiesBySub(sub.ID, app.db)
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		for _, proxy := range proxies {

			pcapsWithOutOutlier := []models.Pcap{}
			app.db.Where("outlier = 0").Where(models.Pcap{Proxy: proxy, SubID: sub.ID}).Find(&pcapsWithOutOutlier)

			pcaps := []models.Pcap{}
			app.db.Where(models.Pcap{Proxy: proxy, SubID: sub.ID}).Find(&pcaps)
			// app.db.Model(&sub).Related(&pcaps).Where(models.Pcap{Proxy: proxy})

			fmt.Printf("%v\t%v\t%v\t%v\n", len(pcapsWithOutOutlier), (len(pcaps) - len(pcapsWithOutOutlier)), proxy, sub.Link)
		}
	}
}
