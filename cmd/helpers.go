package main

import (
	"bufio"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/montanaflynn/stats"
)

func IsUrl(str string) bool {
	u, err := url.Parse(str)
	if err != nil {
		return false
	}
	if u.Scheme == "" || u.Host == "" {
		return false
	}
	resp, err := http.Get(str)
	if err != nil {
		return false
	}
	if resp.StatusCode != 200 {
		return false
	}
	return true
}

func CheckProxy(str string) error {
	req, err := http.NewRequest("GET", "https://duckduckgo.com", nil)
	if err != nil {
		return err
	}
	proxyUrl, err := url.Parse(str)
	if err != nil {
		return err
	}
	client := &http.Client{Timeout: time.Second * 10, Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
	_, err = client.Do(req)
	if err != nil {
		return err
	}
	return nil
}

func ReadFile(s string) ([]string, error) {
	var list []string
	file, err := os.Open(s)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		href := strings.TrimRight((strings.TrimSpace(scanner.Text())), "/")
		list = append(list, href)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func CleanOuitliarsIQR(set []int) []float64 {
	floats := []float64{}
	for _, s := range set {
		floats = append(floats, float64(s))
	}
	qutliers, _ := stats.QuartileOutliers(floats)
	return append(qutliers.Mild, qutliers.Extreme...)
}
