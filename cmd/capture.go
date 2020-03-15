package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/pavel1337/webfingerprint/pkg/capturer"
	"github.com/pavel1337/webfingerprint/pkg/models"
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
		if app.flags.Eth == "" {
			eth, err := capturer.GetInterface()
			if err != nil {
				app.errorLog.Println(err)
				return
			}
			app.flags.Eth = eth
		}
		app.capturePcaps(app.flags.NumberOfInstances, app.flags.Headless, app.flags.ProxyString, app.flags.ProxyType)
	}
	if app.flags.ListAllPcaps {
		app.listPcaps()
	}
	if app.flags.Clean {
		app.cleanOutliers()
		// app.listSubsByProxy()
	}

}

// number of instances as input
func (app *application) capturePcaps(n int, headless bool, proxyString, proxyType string) {
	err := capturer.CheckPerms(app.flags.Eth)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	subs := []models.Sub{}
	app.db.Find(&subs)
	for _, sub := range subs {
		pcaps := []models.Pcap{}
		app.db.Where("outlier = 0").Where(models.Pcap{Proxy: proxyType, SubID: sub.ID}).Find(&pcaps)
		if len(pcaps) >= n {
			continue
		}

		for i := 1; i <= (n - len(pcaps)); i++ {
			app.infoLog.Println("Capturing:", sub.Link)
			cumul, path, err := capturer.OpenBrowser(sub.Link, app.flags.Eth, proxyString, 60, headless)
			if err != nil {
				app.errorLog.Println(err)
				continue
			}
			bcumul, err := json.Marshal(cumul)
			if err != nil {
				app.errorLog.Println(err)
				continue
			}

			pcap := models.Pcap{Path: path, SubID: sub.ID, Proxy: proxyType, BCumul: bcumul}
			app.db.Create(&pcap)
		}
	}
}

func (app *application) listPcaps() {
	subs := []models.Sub{}
	app.db.Find(&subs)
	for _, sub := range subs {
		pcaps := []models.Pcap{}
		app.db.Model(&sub).Related(&pcaps)
		for _, pcap := range pcaps {
			fmt.Println(pcap.ID, pcap.Path, sub.Link, pcap.Proxy, pcap.Outlier)
		}
	}
}

func (app *application) cleanOutliers() {
	subs := []models.Sub{}
	app.db.Find(&subs)

	for _, sub := range subs {

		proxies, err := models.DistinctProxiesBySub(sub.ID, app.db)
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		for _, proxy := range proxies {
			for {
				pcaps := []models.Pcap{}
				app.db.Where("outlier = 0").Where(models.Pcap{Proxy: proxy, SubID: sub.ID}).Find(&pcaps)
				var lcumuls []int
				for _, pcap := range pcaps {
					var cumul [50]int
					err := json.Unmarshal(pcap.BCumul, &cumul)
					if err != nil {
						app.errorLog.Println(err)
						return
					}
					lcumuls = append(lcumuls, cumul[49])
				}

				Outliers := CleanOuitliarsIQR(lcumuls)
				if len(Outliers) == 0 {
					break
				}
				for _, pcap := range pcaps {
					var cumul [50]int
					err := json.Unmarshal(pcap.BCumul, &cumul)
					if err != nil {
						app.errorLog.Println(err)
						return
					}
					for _, Outlier := range Outliers {
						if float64(cumul[49]) == Outlier {
							app.db.Model(&pcap).Update(models.Pcap{Outlier: true})
						}
					}
				}
			}
		}
	}
}
