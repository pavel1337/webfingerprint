package main

import (
	"errors"
	"fmt"

	"github.com/pavel1337/webfingerprint/pkg/capturer"
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
	e := 0
	for _, sub := range s {
		if e > 10 {
			return
		}
		p, err := app.pcaps.GetBySubIdAndProxy(sub.ID, proxyType)
		if err != nil {
			e++
			app.errorLog.Println(err)
			continue
		}
		if len(p) >= n {
			continue
		}
		for i := 1; i <= (n - len(p)); i++ {
			path, err := capturer.OpenBrowser(sub.Title, app.flags.Eth, proxyString, 60, headless)
			if err != nil {
				e++
				app.errorLog.Println(err)
				continue
			}
			_, err = app.pcaps.Insert(path, sub.ID, proxyType)
			if err != nil {
				e++
				app.errorLog.Println(err)
				continue
			}
			e--
		}
		e--
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
		fmt.Println(pcap.ID, pcap.Path, sub.Title, pcap.Proxy)
	}
}
