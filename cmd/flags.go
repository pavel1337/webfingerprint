package main

import (
	"flag"
	"fmt"
	"os"
)

type Flags struct {
	Mode              string
	Dsn               string
	Website           string
	Websites          string
	List              bool
	GenerateSubs      bool
	Depth             int
	ListSubs          string
	ListAllSubs       bool
	Eth               string
	NumberOfInstances int
	CapturePcaps      bool
	Headless          bool
	ListAllPcaps      bool
	ProxyString       string
	ProxyType         string
	LocalIp           string
	ExtractFeatures   bool
}

var Usage = func() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
	fmt.Println("   -dsn\n", "      MySQL database string", `(default "web:password@tcp(172.17.0.2:3306)/webfingerprint?parseTime=true")`)
	fmt.Println("   -mode\n", `      mode of operation (available: "prepare", "capture")`)

	fmt.Println("\noptions for prepare mode:")
	fmt.Println("   -url\n", "      insert url to database")
	fmt.Println("   -urls-file\n", "      insert urls from file to database")
	fmt.Println("   -list-websites\n", `      list urls from database`)
	fmt.Println("   -generate-subs\n", `      generate sub-pages`)
	fmt.Println("   -depth\n", `      depth of sub-pages`, `(default "4")`)
	fmt.Println("   -list-subs-by-website\n", `      list sub-pages for website from database`)
	fmt.Println("   -list-all-subpages\n", `      list urls from database`)

	fmt.Println("\noptions for capture mode:")
	fmt.Println("   -dev\n", `      interface to listen on while capture`, `(default "eth0")`)
	fmt.Println("   -num\n", `      number of instances to capture`, `(default "10")`)
	fmt.Println("   -capture-pcaps\n", `      capture traffic for all subpages in database`)
	fmt.Println("   -headless\n", `      capture in headless mode`)
	fmt.Println("   -proxy-string\n", `      proxy string to use`, `(example "http://127.0.0.1:4001")`)
	fmt.Println("   -proxy-type\n", `      proxy type for the database`, `(default "none")`)
	fmt.Println("   -list-all-pcaps\n", `      list all pcaps from database`)

	fmt.Println("\noptions for visualize mode:")
	fmt.Println("   -local-ip\n", `      local ip used for capture`, `(default will be taken automatically)`)
	fmt.Println("   -extract-features\n", `      extract features out of pcaps in database`)
	// fmt.Println("   -by-website\n", `      visualize traffic for one website`, `(example "www.google.de")`)
	// fmt.Println("   -by-url\n", `      visualize traffic for one url`, `(example "https://www.dmca.com/dashboard")`)

	fmt.Println("\nExamples:")
	fmt.Printf("    %v -mode prepare -url https://google.de -proxy http://127.0.0.1:4001 -i eth0\n", os.Args[0])
	fmt.Printf("    %v -mode capture -u https://google.de -proxy http://127.0.0.1:4001 -i eth0\n", os.Args[0])
	fmt.Printf("    %v -mode visualize -path captured_traffic -D\n", os.Args[0])
}

func ParseFlags() Flags {
	var f Flags
	flag.StringVar(&f.Mode, "mode", "", `mode of operation ("prepare", "capture" and "visualize" available)`)
	flag.StringVar(&f.Dsn, "dsn", "web:password@tcp(172.17.0.2:3306)/webfingerprint?parseTime=true", "MySQL database string")
	flag.StringVar(&f.Website, "url", "", "insert url to database")
	flag.StringVar(&f.Websites, "urls-file", "", "insert urls from file to database")
	flag.BoolVar(&f.List, "list-websites", false, `list urls from database`)
	flag.BoolVar(&f.GenerateSubs, "generate-subs", false, `generate sub-pages`)
	flag.IntVar(&f.Depth, "depth", 4, `depth of sub-pages`)
	flag.StringVar(&f.ListSubs, "list-subs-by-website", "", `list sub-pages for website from database`)
	flag.BoolVar(&f.ListAllSubs, "list-all-subpages", false, `list urls from database`)
	flag.StringVar(&f.Eth, "dev", "eth0", `interface to listen on while capture`)
	flag.IntVar(&f.NumberOfInstances, "num", 10, `number of instances to capture`)
	flag.BoolVar(&f.CapturePcaps, "capture-pcaps", false, `capture traffic for all subpages in database`)
	flag.BoolVar(&f.Headless, "headless", false, `capture in headless mode`)
	flag.StringVar(&f.ProxyString, "proxy-string", "", `proxy string to use`)
	flag.StringVar(&f.ProxyType, "proxy-type", "none", `proxy type for the database`)
	flag.BoolVar(&f.ListAllPcaps, "list-all-pcaps", false, `list all pcaps from database`)
	flag.StringVar(&f.LocalIp, "local-ip", "", `local ip used for capture`)
	flag.BoolVar(&f.ExtractFeatures, "extract-features", false, `extract features out of pcaps in database`)

	flag.Parse()
	return f
}
