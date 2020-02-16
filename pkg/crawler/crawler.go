package crawler

import (
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
)

func getHref(t html.Token) (ok bool, href string) {
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = strings.TrimRight((strings.TrimSpace(a.Val)), "/")
			ok = true
		}
	}
	return
}
func getUserAgent() string {
	useragents := []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3835.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3831.6 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3818.0 Safari/537.36 Edg/77.0.189.3",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3790.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3782.0 Safari/537.36 Edg/76.0.152.0"}
	rand.Seed(time.Now().Unix())
	return useragents[rand.Intn(len(useragents))]
}

func crawl(link string) []string {
	h, _ := url.Parse(link)
	hostname := h.Hostname()
	var links []string
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return links
	}
	req.Header.Set("User-Agent", getUserAgent())
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	if err != nil {
		return links
	}
	b := resp.Body
	z := html.NewTokenizer(b)
	for {
		tt := z.Next()
		switch {
		case tt == html.ErrorToken:
			return links
		case tt == html.StartTagToken:
			t := z.Token()
			if t.Data != "a" {
				continue
			}
			ok, href := getHref(t)
			if !ok {
				continue
			}
			if strings.Index(href, "http") != 0 {
				href = link + href
			}
			if !strings.Contains(href, hostname) {
				continue
			}
			if strings.Contains(href, ".pdf") {
				continue
			}
			if strings.Contains(href, ".zip") {
				continue
			}
			if strings.Contains(href, "javascript") {
				continue
			}
			if strings.Contains(href, "#") {
				continue
			}
			if strings.Contains(href, "?") {
				continue
			}
			if strings.Contains(href[6:], ":") {
				continue
			}
			links = append(links, href)
		}
	}
	return links
}

func getRandomUrl(hrefs []string) string {
	rand.Seed(time.Now().Unix())
	for {
		returnlink := hrefs[rand.Intn(len(hrefs))]
		if returnlink != "" {
			return returnlink
		}
	}
}

func unique(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func Crawler(link string, depth int) []string {
	links := []string{}
	links = append(links, link)
	i := 0
	for {
		i++
		links = append(links, crawl(getRandomUrl(links))...)
		links = unique(links)
		if len(links) > depth {
			break
		}
		if i > 10 {
			return links[1:]
		}
	}
	returnlinks := []string{}
	for i, l := range links {
		if i == 0 {
			continue
		}
		returnlinks = append(returnlinks, l)
		if len(returnlinks) >= depth {
			break
		}
	}
	return returnlinks
}
