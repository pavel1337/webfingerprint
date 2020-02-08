package main

import (
	"bufio"
	"net/http"
	"net/url"
	"os"
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

func ReadFile(s string) ([]string, error) {
	var list []string
	file, err := os.Open(s)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		list = append(list, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return list, nil
}
