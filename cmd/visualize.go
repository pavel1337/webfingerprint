package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/pavel1337/webfingerprint/pkg/models"
	"github.com/pavel1337/webfingerprint/pkg/visualizer"
	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
)

func (app *application) Visualize() {
	if app.flags.VisualizeByUrlAndProxy != "" && app.flags.ProxyType != "" {
		app.visualizeByUrlAndProxy(app.flags.VisualizeByUrlAndProxy, app.flags.ProxyType, app.flags.Path)
	}
	if app.flags.VisualizeByUrl != "" {
		app.visualizeByUrl(app.flags.VisualizeByUrl, app.flags.Path)
	}
	if app.flags.VisualizeByWebsiteAndProxy != "" && app.flags.ProxyType != "" {
		app.visualizeByWebsiteAndProxy(app.flags.VisualizeByWebsiteAndProxy, app.flags.ProxyType, app.flags.Path)
	}
	if app.flags.VisualizeOnlyMainPages && app.flags.ProxyType != "" {
		app.visualizeOnlyMainPages(app.flags.ProxyType, app.flags.Path)
	}
}

func (app *application) visualizeByUrlAndProxy(url, proxy, path string) {
	sub := models.Sub{}
	app.db.Where(models.Sub{Link: url}).First(&sub)

	pcaps := []models.Pcap{}
	app.db.Where(models.Pcap{Proxy: proxy, SubID: sub.ID}).Find(&pcaps).Limit(app.flags.NumberOfInstances)

	firstOutlier := false
	firstNotOutlier := false
	arser := []chart.Series{}
	for _, p := range pcaps {
		var cumul [50]int
		err := json.Unmarshal(p.BCumul, &cumul)
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		if !p.Outlier {
			if !firstNotOutlier {
				arser = append(arser, visualizer.CraftTheSeriesWithLegend(cumul, drawing.ColorBlue, "not outlier"))
				firstNotOutlier = true
			} else {
				arser = append(arser, visualizer.CraftTheSeries(cumul, drawing.ColorBlue))
			}
		} else {
			if !firstOutlier {
				arser = append(arser, visualizer.CraftTheSeriesWithLegend(cumul, drawing.ColorRed, "outlier"))
				firstOutlier = true
			} else {
				arser = append(arser, visualizer.CraftTheSeries(cumul, drawing.ColorRed))
			}
		}
	}
	if path == "" {
		path = fmt.Sprintf("charts/%d-%s.png", sub.ID, proxy)
	}
	err := visualizer.SaveTheGraph(arser, path)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	fmt.Println("file saved to:", path)
}

func (app *application) visualizeByWebsiteAndProxy(website, proxy, path string) {
	w := models.Website{}
	app.db.Where(models.Website{Hostname: website}).First(&w)

	subs := []models.Sub{}
	app.db.Model(&w).Related(&subs)
	if len(subs) == 0 {
		app.errorLog.Println(errors.New("No subpages for the website in the database"))
		return
	}

	arser := []chart.Series{}

	for _, sub := range subs {
		pcaps := []models.Pcap{}
		app.db.Where("outlier = 0").Where(models.Pcap{Proxy: proxy, SubID: sub.ID}).Find(&pcaps).Limit(app.flags.NumberOfInstances)
		if len(pcaps) == 0 {
			app.errorLog.Println(errors.New("No pcaps for the subpage in the database"))
			return
		}
		color := visualizer.GetRandomColor()
		for i, pcap := range pcaps {
			var cumul [50]int
			err := json.Unmarshal(pcap.BCumul, &cumul)
			if err != nil {
				app.errorLog.Println(err)
				return
			}
			if i == 0 {
				arser = append(arser, visualizer.CraftTheSeriesWithLegend(cumul, color, ("subpage: "+sub.Link)))
				continue
			}
			arser = append(arser, visualizer.CraftTheSeries(cumul, color))
		}
	}
	if path == "" {
		path = fmt.Sprintf("charts/%s-%s.png", w.Hostname, proxy)
	}
	err := visualizer.SaveTheGraph(arser, path)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	fmt.Println("file saved to:", path)
}

func (app *application) visualizeByUrl(url, path string) {
	sub := models.Sub{}
	app.db.Where(models.Sub{Link: url}).First(&sub)
	proxies, err := models.DistinctProxiesBySub(sub.ID, app.db)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	var arser []chart.Series
	for _, proxy := range proxies {
		color := visualizer.GetRandomColor()
		pcaps := []models.Pcap{}

		app.db.Where("outlier = 0").Where(models.Pcap{Proxy: proxy, SubID: sub.ID}).Find(&pcaps).Limit(app.flags.NumberOfInstances)

		for i, pcap := range pcaps {
			var cumul [50]int
			err := json.Unmarshal(pcap.BCumul, &cumul)
			if err != nil {
				app.errorLog.Println(err)
				return
			}
			if i == 0 {
				arser = append(arser, visualizer.CraftTheSeriesWithLegend(cumul, color, ("proxy: "+proxy)))
				continue
			}
			arser = append(arser, visualizer.CraftTheSeries(cumul, color))
		}
	}
	if path == "" {
		path = fmt.Sprintf("charts/%v.png", sub.ID)
	}
	err = visualizer.SaveTheGraph(arser, path)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	fmt.Println("file saved to:", path)
}

func (app *application) visualizeOnlyMainPages(proxy, path string) {
	websites := []models.Website{}
	app.db.Find(&websites)

	var arser []chart.Series
	for _, website := range websites {
		sub := models.Sub{}
		app.db.Where(models.Sub{WebsiteID: website.ID}).First(&sub)
		color := visualizer.GetRandomColor()

		pcaps := []models.Pcap{}
		app.db.Where("outlier = 0").Where(models.Pcap{Proxy: proxy, SubID: sub.ID}).Find(&pcaps).Limit(app.flags.NumberOfInstances)
		if len(pcaps) == 0 {
			app.errorLog.Printf(`There is no pcap files for "%s" link and proxy "%s"`, sub.Link, proxy)
			return
		}
		for i, pcap := range pcaps {
			var cumul [50]int
			err := json.Unmarshal(pcap.BCumul, &cumul)
			if err != nil {
				app.errorLog.Println(err)
				return
			}
			if i == 0 {
				arser = append(arser, visualizer.CraftTheSeriesWithLegend(cumul, color, ("website: "+website.Hostname)))
				continue
			}
			arser = append(arser, visualizer.CraftTheSeries(cumul, color))
		}
	}
	if path == "" {
		path = "charts/only_main.png"
	}

	err := visualizer.SaveTheGraph(arser, path)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	fmt.Println("file saved to:", path)
}
