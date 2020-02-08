package main

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	Mode        string
	Dsn         string
	Website     string
	Websites    string
	List        bool
	Generate    bool
	Depth       int
	ListSubs    string
	ListAllSubs bool
}

var Usage = func() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Println("\nExamples:")
	fmt.Printf("    %v -mode prepare -url https://google.de -proxy http://127.0.0.1:4001 -i eth0\n", os.Args[0])
	fmt.Printf("    %v -mode capture -u https://google.de -proxy http://127.0.0.1:4001 -i eth0\n", os.Args[0])
	fmt.Printf("    %v -mode visualize -path captured_traffic -D\n", os.Args[0])
}

func ParseConfig() Config {
	var c Config
	flag.StringVar(&c.Mode, "mode", "", `mode of operation ("prepare", "capture" and "visualize" available)`)
	flag.StringVar(&c.Dsn, "dsn", "web:password@tcp(172.17.0.2:3306)/webfingerprint?parseTime=true", "MySQL database string")
	flag.StringVar(&c.Website, "url", "", "insert url to database")
	flag.StringVar(&c.Websites, "urls-file", "", "insert urls from file to database")
	flag.BoolVar(&c.List, "list-websites", false, `list urls from database`)
	flag.BoolVar(&c.Generate, "generate-subs", false, `generate sub-pages`)
	flag.IntVar(&c.Depth, "depth", 4, `depth of sub-pages`)
	flag.StringVar(&c.ListSubs, "list-subs-by-website", "", `list sub-pages for website from database`)
	flag.BoolVar(&c.ListAllSubs, "list-all-subpages", false, `list urls from database`)
	flag.Parse()
	return c
}
