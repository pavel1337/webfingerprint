package main

import (
	"bufio"
	"errors"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
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

func ExternalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?")
}
