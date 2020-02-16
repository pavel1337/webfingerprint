package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/pavel1337/webfingerprint/pkg/capturer"
	"github.com/pavel1337/webfingerprint/pkg/visualizer"
	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
)

func (app *application) Pcaps() {
	if app.flags.CapturePcaps {
		if app.flags.ProxyString != "" && app.flags.ProxyType == "none" {
			app.errorLog.Println(errors.New("please do not confuse database and define -proxy-type flag!"))
			return
		}
		if app.flags.ProxyString != "" {
			err := CheckProxy(app.flags.ProxyString)
			if err != nil {
				app.errorLog.Println(err)
				return
			}
		}
		app.capturePcaps(app.flags.NumberOfInstances, app.flags.Headless, app.flags.ProxyString, app.flags.ProxyType)
	}
	if app.flags.ListAllPcaps {
		app.listPcaps()
	}
	if app.flags.Clean {
		app.cleanOutliers()
		app.listSubsByProxy()
	}
	if strings.TrimSpace(app.flags.VisualizeByUrlAndProxy) != "" && strings.TrimSpace(app.flags.ProxyType) != "" {
		app.visualizeByUrlAndProxy(app.flags.VisualizeByUrlAndProxy, app.flags.ProxyType, app.flags.Path)
	}
	if strings.TrimSpace(app.flags.VisualizeByUrl) != "" {
		app.visualizeByUrl(app.flags.VisualizeByUrl, app.flags.Path)
	}
	if strings.TrimSpace(app.flags.VisualizeByWebsiteAndProxy) != "" && strings.TrimSpace(app.flags.ProxyType) != "" {
		app.visualizeByWebsiteAndProxy(app.flags.VisualizeByWebsiteAndProxy, app.flags.ProxyType, app.flags.Path)
	}

}

// number of instances as input
func (app *application) capturePcaps(n int, headless bool, proxyString, proxyType string) {
	err := capturer.CheckPerms(app.flags.Eth)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	s, err := app.subs.List()
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	for _, sub := range s {
		p, err := app.pcaps.GetBySubIdAndProxyAndNotOutlier(sub.ID, proxyType)
		if err != nil {
			app.errorLog.Println(err)
			continue
		}
		if len(p) >= n {
			continue
		}
		for i := 1; i <= (n - len(p)); i++ {
			app.infoLog.Println("Capturing:", sub.Title)
			cumul, path, err := capturer.OpenBrowser(sub.Title, app.flags.Eth, proxyString, 60, headless)
			if err != nil {
				app.errorLog.Println(err)
				continue
			}
			_, err = app.pcaps.Insert(path, sub.ID, proxyType, cumul)
			if err != nil {
				app.errorLog.Println(err)
				return
			}
		}
	}
}

func (app *application) listPcaps() {
	pcaps, err := app.pcaps.List()
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	for _, pcap := range pcaps {
		sub, err := app.subs.GetById(pcap.SubId)
		if err != nil {
			app.errorLog.Println(err)
			continue
		}
		fmt.Println(pcap.ID, pcap.Path, sub.Title, pcap.Proxy, pcap.Outlier)
	}
}

func (app *application) visualizeByUrlAndProxy(url, proxy, path string) {
	sub, err := app.subs.GetByUrl(url)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	pcaps, err := app.pcaps.GetBySubIdAndProxy(sub.ID, proxy)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	firstOutlier := false
	firstNotOutlier := false
	arser := []chart.Series{}
	for _, p := range pcaps {
		if !p.Outlier {
			if !firstNotOutlier {
				arser = append(arser, visualizer.CraftTheSeriesWithLegend(p.Cumul, drawing.ColorBlue, "not outlier"))
				firstNotOutlier = true
			} else {
				arser = append(arser, visualizer.CraftTheSeries(p.Cumul, drawing.ColorBlue))
			}
		} else {
			if !firstOutlier {
				arser = append(arser, visualizer.CraftTheSeriesWithLegend(p.Cumul, drawing.ColorRed, "outlier"))
				firstOutlier = true
			} else {
				arser = append(arser, visualizer.CraftTheSeries(p.Cumul, drawing.ColorRed))
			}
		}
	}
	if path == "" {
		path = strconv.Itoa(sub.ID) + "-" + proxy + ".png"
	}
	err = visualizer.SaveTheGraph(arser, path)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	fmt.Println("file saved to:", path)
}

func (app *application) visualizeByWebsiteAndProxy(website, proxy, path string) {
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

	arser := []chart.Series{}

	for _, sub := range subs {
		pcaps, err := app.pcaps.GetBySubIdAndProxyAndNotOutlier(sub.ID, proxy)
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		color := visualizer.GetRandomColor()
		for i, pcap := range pcaps {
			if i == 0 {
				arser = append(arser, visualizer.CraftTheSeriesWithLegend(pcap.Cumul, color, ("subpage: "+sub.Title)))
				continue
			}
			arser = append(arser, visualizer.CraftTheSeries(pcap.Cumul, color))
		}
	}
	if path == "" {
		path = w.Title + "-" + proxy + ".png"
	}
	err = visualizer.SaveTheGraph(arser, path)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	fmt.Println("file saved to:", path)
}

func (app *application) visualizeByUrl(url, path string) {
	sub, err := app.subs.GetByUrl(url)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	proxies, err := app.pcaps.ListProxiesBySubid(sub.ID)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	arser := []chart.Series{}
	for _, proxy := range proxies {
		color := visualizer.GetRandomColor()
		pcaps, err := app.pcaps.GetBySubIdAndProxyAndNotOutlier(sub.ID, proxy)
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		for i, pcap := range pcaps {
			if i == 0 {
				arser = append(arser, visualizer.CraftTheSeriesWithLegend(pcap.Cumul, color, ("proxy: "+proxy)))
				continue
			}
			arser = append(arser, visualizer.CraftTheSeries(pcap.Cumul, color))
		}
	}
	if path == "" {
		path = strconv.Itoa(sub.ID) + ".png"
	}
	err = visualizer.SaveTheGraph(arser, path)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	fmt.Println("file saved to:", path)
}

func (app *application) cleanOutliers() {
	subs, err := app.subs.List()
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	proxies, err := app.pcaps.ListProxies()
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	for _, sub := range subs {
		for _, proxy := range proxies {
			for {
				pcaps, err := app.pcaps.GetBySubIdAndProxyAndNotOutlier(sub.ID, proxy)
				if err != nil {
					app.errorLog.Println(err)
					return
				}
				var lcumuls []int
				for _, pcap := range pcaps {
					lcumuls = append(lcumuls, pcap.Cumul[49])
				}
				Outliers := CleanOuitliarsIQR(lcumuls)
				if len(Outliers) == 0 {
					break
				}
				for _, pcap := range pcaps {
					for _, Outlier := range Outliers {
						if float64(pcap.Cumul[49]) == Outlier {
							err = app.pcaps.SetOutlierById(pcap.ID)
							if err != nil {
								app.errorLog.Println(err)
								return
							}
						}
					}
				}
			}
		}
	}
}
