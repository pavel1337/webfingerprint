package main

import (
	"encoding/csv"
	"encoding/json"
	"os"
	"strconv"

	"github.com/pavel1337/webfingerprint/pkg/models"
)

func (app *application) Export() {
	if app.flags.ExportAll {
		app.saveTheDataSetAll()
	}
	if app.flags.ExportByWebsiteAndProxy != "" && app.flags.ProxyType != "" {
		app.saveTheDataSetByWebSiteAndProxy(app.flags.ExportByWebsiteAndProxy, app.flags.ProxyType)
	}
	if app.flags.ExportOnlyMainPages && app.flags.ProxyType != "" {
		app.saveTheDataSetOnlyMainPages(app.flags.ProxyType)
	}
}

func (app *application) saveTheDataSetAll() {
	f, err := os.Create("dataset.csv")
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	defer f.Close()

	subs := []models.Sub{}
	app.db.Find(&subs)

	var csvDataSet [][]string
	for _, sub := range subs {
		pcaps := []models.Pcap{}
		app.db.Where("outlier = 0").Where(models.Pcap{SubID: sub.ID}).Limit(app.flags.NumberOfInstances).Find(&pcaps)
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
	defer f.Close()

	website := models.Website{}
	app.db.Where(models.Website{Hostname: websiteHostname}).First(&website)

	subs := []models.Sub{}
	app.db.Model(&website).Related(&subs)

	if len(subs) == 0 {
		app.errorLog.Printf("There is no %s website in the database", websiteHostname)
		return
	}

	var csvDataSet [][]string
	for _, sub := range subs {
		pcaps := []models.Pcap{}
		app.db.Where("outlier = 0").Where(models.Pcap{Proxy: proxy, SubID: sub.ID}).Limit(app.flags.NumberOfInstances).Find(&pcaps)
		if len(pcaps) == 0 {
			app.errorLog.Printf(`There is no pcap files for "%s" link and proxy "%s"`, sub.Link, proxy)
			return
		}
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

func (app *application) saveTheDataSetOnlyMainPages(proxy string) {
	f, err := os.Create("dataset.csv")
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	defer f.Close()

	websites := []models.Website{}
	app.db.Find(&websites)

	var csvDataSet [][]string
	for _, website := range websites {
		sub := models.Sub{}
		app.db.Where(models.Sub{WebsiteID: website.ID}).First(&sub)

		pcaps := []models.Pcap{}
		app.db.Where("outlier = 0").Where(models.Pcap{Proxy: proxy, SubID: sub.ID}).Limit(app.flags.NumberOfInstances).Find(&pcaps)
		if len(pcaps) == 0 {
			app.errorLog.Printf(`There is no pcap files for "%s" link and proxy "%s"`, sub.Link, proxy)
			return
		}
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
