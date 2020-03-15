package main

import (
	"flag"
	"fmt"
	"os"
)

type Flags struct {
	Mode                       string
	Dsn                        string
	Website                    string
	Websites                   string
	List                       bool
	GenerateSubs               bool
	Depth                      int
	ListSubs                   string
	ListSubsWithPcap           bool
	ListAllSubs                bool
	Eth                        string
	NumberOfInstances          int
	CapturePcaps               bool
	Headless                   bool
	ListAllPcaps               bool
	ProxyString                string
	ProxyType                  string
	Clean                      bool
	VisualizeByUrl             string
	VisualizeByUrlAndProxy     string
	Path                       string
	VisualizeByWebsiteAndProxy string
	ExportAll                  bool
	ExportByWebsiteAndProxy    string
	ExportOnlyMainPages        bool
	VisualizeOnlyMainPages     bool
	LearnRandomForest          bool
	LearnKNN                   bool
	DatasetPath                string
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
	fmt.Println("   -list-subs-with-pcap-info\n", `      list sub-pages with pcap info`)
	fmt.Println("   -list-all-subpages\n", `      list urls from database`)

	fmt.Println("\noptions for capture mode:")
	fmt.Println("   -dev\n", `      interface to listen on while capture`, `(autodiscover by default)`)
	fmt.Println("   -num\n", `      number of instances to capture`, `(default "10")`)
	fmt.Println("   -capture-pcaps\n", `      capture traffic for all subpages in database`)
	fmt.Println("   -headless\n", `      capture in headless mode`)
	fmt.Println("   -proxy-string\n", `      proxy string to use`, `(example "http://127.0.0.1:4001")`)
	fmt.Println("   -proxy-type\n", `      proxy type for the database`, `(default "none")`)
	fmt.Println("   -list-all-pcaps\n", `      list all pcaps from database`)

	fmt.Println("\noptions for visualize mode:")
	fmt.Println("   -clean\n", `      clean outliers`)
	fmt.Println("   -by-url\n", `      visualize traffic for one url`, `(example "https://www.dmca.com/dashboard")`)
	fmt.Println("   -by-url-proxy\n", `      visualize traffic for one url`, `(example "https://www.dmca.com/dashboard")`)
	fmt.Println("   -by-website-proxy\n", `      visualize traffic for one website`, `(example "www.dmca.com")`)
	fmt.Println("   -path\n", `      path to the plot`, `(example "plot.png")`)
	fmt.Println("   -visualize-only-main\n", `      visualize only main pages for given proxy`)

	fmt.Println("\noptions for export mode:")
	fmt.Println("   -export-all\n", `      export all as csv`)
	fmt.Println("   -export-by-website-proxy\n", `      export by website and proxy as csv`)
	fmt.Println("   -export-only-main\n", `      export only main pages for given proxy as csv`)

	fmt.Println("\noptions for learn mode:")
	fmt.Println("   -dataset-path\n", `      path to the dataset (defaults to dataset.csv) `)
	fmt.Println("   -learn-random-forest\n", `      use random forest for ML`)
	fmt.Println("   -learn-knn\n", `      use knn for ML`)
}

func ParseFlags() Flags {
	var f Flags
	flag.StringVar(&f.Mode, "mode", "", `mode of operation ("prepare", "capture" and "visualize" available)`)
	flag.StringVar(&f.Dsn, "dsn", "web:password@tcp(172.17.0.2:3306)/gorm_test?charset=utf8&parseTime=True&loc=Local", "MySQL database string")
	flag.StringVar(&f.Website, "url", "", "insert url to database")
	flag.StringVar(&f.Websites, "urls-file", "", "insert urls from file to database")
	flag.BoolVar(&f.List, "list-websites", false, `list urls from database`)
	flag.BoolVar(&f.GenerateSubs, "generate-subs", false, `generate sub-pages`)
	flag.IntVar(&f.Depth, "depth", 4, `depth of sub-pages`)
	flag.StringVar(&f.ListSubs, "list-subs-by-website", "", `list sub-pages for website from database`)
	flag.BoolVar(&f.ListSubsWithPcap, "list-subs-with-pcap-info", false, `list sub-pages with pcap info`)
	flag.BoolVar(&f.ListAllSubs, "list-all-subpages", false, `list urls from database`)
	flag.StringVar(&f.Eth, "dev", "", `interface to listen on while capture`)
	flag.IntVar(&f.NumberOfInstances, "num", 20, `number of instances to capture`)
	flag.BoolVar(&f.CapturePcaps, "capture-pcaps", false, `capture traffic for all subpages in database`)
	flag.BoolVar(&f.Headless, "headless", false, `capture in headless mode`)
	flag.StringVar(&f.ProxyString, "proxy-string", "", `proxy string to use`)
	flag.StringVar(&f.ProxyType, "proxy-type", "none", `proxy type for the database`)
	flag.BoolVar(&f.ListAllPcaps, "list-all-pcaps", false, `list all pcaps from database`)
	flag.BoolVar(&f.Clean, "clean", false, `clean outliers`)
	flag.StringVar(&f.VisualizeByUrl, "by-url", "", `visualize traffic for one url (example "https://www.dmca.com/dashboard")`)
	flag.StringVar(&f.VisualizeByUrlAndProxy, "by-url-proxy", "", `visualize traffic for one url (example "https://www.dmca.com/dashboard")`)
	flag.StringVar(&f.VisualizeByWebsiteAndProxy, "by-website-proxy", "", `visualize traffic for one website (example "www.dmca.com")`)
	flag.StringVar(&f.Path, "path", "", `path to the plot (example "plot.png")`)
	flag.BoolVar(&f.ExportAll, "export-all", false, `export all as csv`)
	flag.StringVar(&f.ExportByWebsiteAndProxy, "export-by-website-proxy", "", `export by website and proxy as csv`)
	flag.BoolVar(&f.ExportOnlyMainPages, "export-only-main", false, `export only main pages for given proxy as csv`)
	flag.BoolVar(&f.VisualizeOnlyMainPages, "visualize-only-main", false, `visualize only main pages for given proxy`)
	flag.StringVar(&f.DatasetPath, "dataset-path", "", `path to the dataset (defaults to dataset.csv)`)
	flag.BoolVar(&f.LearnRandomForest, "learn-random-forest", false, `use random forest for ML`)
	flag.BoolVar(&f.LearnKNN, "learn-knn", false, `use knn for ML`)

	flag.Parse()
	return f
}
