package visualizer

import (
	"math/rand"
	"os"
	"time"

	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
)

func GetRandomColor() drawing.Color {
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(256)
	g := rand.Intn(256)
	b := rand.Intn(256)
	return drawing.Color{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: 255,
	}
}

func CraftTheSeries(v [50]int, color drawing.Color) chart.Series {
	keys := []float64{}
	value := []float64{}
	for k, v := range v {
		keys = append(keys, float64(k))
		value = append(value, float64(v))
	}
	ser := chart.ContinuousSeries{
		Style: chart.Style{
			StrokeColor: color,
		},
		XValues: keys,
		YValues: value,
	}
	return ser
}

func CraftTheSeriesWithLegend(v [50]int, color drawing.Color, name string) chart.Series {
	keys := []float64{}
	value := []float64{}
	for k, v := range v {
		keys = append(keys, float64(k))
		value = append(value, float64(v))
	}
	ser := chart.ContinuousSeries{
		Name: name,
		Style: chart.Style{
			StrokeColor: color,
		},
		XValues: keys,
		YValues: value,
	}
	return ser
}

func SaveTheGraph(arser []chart.Series, path string) error {
	graph := chart.Chart{
		YAxis: chart.YAxis{
			Name: "Cumulative packet length, Bytes",
		},
		XAxis: chart.XAxis{
			Name: "Samples count",
			Ticks: []chart.Tick{
				{Value: 0.0, Label: "0"},
				{Value: 5.0, Label: "5"},
				{Value: 10.0, Label: "10"},
				{Value: 15.0, Label: "15"},
				{Value: 20.0, Label: "20"},
				{Value: 25.0, Label: "25"},
				{Value: 30.0, Label: "30"},
				{Value: 35.0, Label: "35"},
				{Value: 40.0, Label: "40"},
				{Value: 44.5, Label: "45"},
				{Value: 49.0, Label: "50"},
			},
		},
		Series: arser,
	}
	graph.Elements = []chart.Renderable{
		chart.Legend(&graph),
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	graph.Render(chart.PNG, f)
	return nil
}
