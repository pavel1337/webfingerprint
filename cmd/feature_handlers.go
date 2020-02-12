package main

import (
	"fmt"

	"github.com/pavel1337/webfingerprint/pkg/extractor"
)

func (app *application) Feautures() {
	// app.extractFeautures()
	app.listFeatures()
}

func (app *application) extractFeautures() {
	ps, err := app.pcaps.List()
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	if app.flags.LocalIp == "" {
		app.flags.LocalIp, err = ExternalIP()
		if err != nil {
			app.errorLog.Println(err)
			return
		}
	}

	for _, p := range ps {
		featureset, err := extractor.Extract(p.Path, app.flags.LocalIp)
		if err != nil {
			app.errorLog.Println(err)
			continue
		}
		_, err = app.features.Insert(p.ID, featureset)
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		app.infoLog.Println(featureset)
	}
}

func (app *application) listFeatures() {
	features, err := app.features.List()
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	for _, feature := range features {
		fmt.Println(feature.FeatureSetJson.Cumul)
	}
}
