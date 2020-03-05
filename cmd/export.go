package main

import (
	"encoding/csv"
	"encoding/json"
	"os"
	"strconv"

	"github.com/pavel1337/wasm/pkg/models"
)

func (app *application) Export() {
	if app.flags.ExportAll {
		app.saveTheDataSetAll()
	}
	if app.flags.ExportByWebsiteAndProxy != "" && app.flags.ProxyType != "" {
		app.saveTheDataSetByWebSiteAndProxy(app.flags.ExportByWebsiteAndProxy, app.flags.ProxyType)
	}
}

func (app *application) saveTheDataSetAll() {
	f, err := os.Create("dataset.csv")
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	subs := []models.Sub{}
	app.db.Find(&subs)

	var csvDataSet [][]string
	for _, sub := range subs {
		pcaps := []models.Pcap{}
		app.db.Model(&sub).Related(&pcaps)
		for _, pcap := range pcaps {
			var cumul [50]int
			err := json.Unmarshal(pcap.BCumul, &cumul)
			if err != nil {
				app.errorLog.Println(err)
				return
			}

			var csvData []string
			for _, i := range cumul {
				csvData = append(csvData, strconv.Itoa(i))
			}
			csvData = append(csvData, (sub.Link + "_" + pcap.Proxy))
			csvDataSet = append(csvDataSet, csvData)
		}

	}

	w := csv.NewWriter(f)
	w.WriteAll(csvDataSet) // calls Flush internally
	err = f.Close()
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	if err := w.Error(); err != nil {
		app.errorLog.Println(err)
		return
	}
}

func (app *application) saveTheDataSetByWebSiteAndProxy(websiteHostname, proxy string) {
	f, err := os.Create("dataset.csv")
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	website := models.Website{}
	app.db.Where(models.Website{Hostname: websiteHostname}).First(&website)

	subs := []models.Sub{}
	app.db.Model(&website).Related(&subs)

	var csvDataSet [][]string
	for _, sub := range subs {
		pcaps := []models.Pcap{}
		app.db.Where("outlier = 0").Where(models.Pcap{Proxy: proxy, SubID: sub.ID}).Find(&pcaps)

		for _, pcap := range pcaps {
			var cumul [50]int
			err := json.Unmarshal(pcap.BCumul, &cumul)
			if err != nil {
				app.errorLog.Println(err)
				return
			}
			var csvData []string
			for _, i := range cumul {
				csvData = append(csvData, strconv.Itoa(i))
			}
			csvData = append(csvData, (sub.Link + "_" + pcap.Proxy))
			csvDataSet = append(csvDataSet, csvData)
		}
	}

	w := csv.NewWriter(f)
	w.WriteAll(csvDataSet) // calls Flush internally
	err = f.Close()
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	if err := w.Error(); err != nil {
		app.errorLog.Println(err)
		return
	}
}
