package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/sjwhitworth/golearn/base"
	"github.com/sjwhitworth/golearn/ensemble"
	"github.com/sjwhitworth/golearn/evaluation"
	"github.com/sjwhitworth/golearn/knn"
)

func (app *application) Learn() {
	if app.flags.LearnRandomForest {
		app.learnRandomForest(app.flags.DatasetPath)
	}
	if app.flags.LearnKNN {
		app.learnKNN(app.flags.DatasetPath)
	}
}

func (app *application) learnRandomForest(path string) {
	if path == "" {
		path = "dataset.csv"
	}

	var cls base.Classifier

	iris, err := base.ParseCSVToInstances(path, false)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	rand.Seed(int64(time.Now().Nanosecond()))

	cls = ensemble.NewRandomForest(100, 4)
	cfs, err := evaluation.GenerateCrossFoldValidationConfusionMatrices(iris, cls, 10)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	for _, cf := range cfs {
		fmt.Println(evaluation.GetSummary(cf))
	}

	mean, variance := evaluation.GetCrossValidatedMetric(cfs, evaluation.GetAccuracy)
	stdev := math.Sqrt(variance)

	fmt.Printf("Mean: %.2f\t Stdev: (+/- %.2f)\n", mean, stdev)
}

func (app *application) learnKNN(path string) {
	if path == "" {
		path = "dataset.csv"
	}

	var cls base.Classifier

	iris, err := base.ParseCSVToInstances(path, false)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	rand.Seed(int64(time.Now().Nanosecond()))

	cls = knn.NewKnnClassifier("euclidean", "linear", 20)
	cfs, err := evaluation.GenerateCrossFoldValidationConfusionMatrices(iris, cls, 10)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	for _, cf := range cfs {
		fmt.Println(evaluation.GetSummary(cf))
	}

	mean, variance := evaluation.GetCrossValidatedMetric(cfs, evaluation.GetAccuracy)
	stdev := math.Sqrt(variance)

	fmt.Printf("%.2f\t(+/- %.2f)\n", mean, stdev)
}
