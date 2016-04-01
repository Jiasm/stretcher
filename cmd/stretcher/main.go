package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/fujiwara/stretcher"
	"github.com/tcnksm/go-latest"
)

var (
	version   string
	buildDate string
)

func main() {
	var (
		showVersion  bool
		delay        float64
		maxBandWidth string
	)
	flag.BoolVar(&showVersion, "v", false, "show version")
	flag.BoolVar(&showVersion, "version", false, "show version")
	flag.Float64Var(&delay, "random-delay", 0, "sleep [0,random-delay) sec on start")
	flag.StringVar(&maxBandWidth, "max-bandwidth", "", "max bandwidth for download src archives (Bytes/sec)")
	flag.Parse()

	if showVersion {
		fmt.Println("version:", version)
		fmt.Println("build:", buildDate)
		checkLatest(version)
		return
	}
	log.Println("stretcher version:", version)
	stretcher.Version = version
	if bw, err := humanize.ParseBytes(maxBandWidth); err == nil {
		stretcher.MaxBandWidth = bw
	}
	stretcher.Init(stretcher.RandomTime(delay))
	err := stretcher.Run()
	if err != nil {
		log.Println(err)
		if os.Getenv("CONSUL_INDEX") != "" {
			// ensure exit 0 when running under `consul watch`
			return
		} else {
			os.Exit(1)
		}
	}
}

func fixVersionStr(v string) string {
	v = strings.TrimPrefix(v, "v")
	vs := strings.Split(v, "-")
	return vs[0]
}

func checkLatest(version string) {
	version = fixVersionStr(version)
	githubTag := &latest.GithubTag{
		Owner:             "fujiwara",
		Repository:        "stretcher",
		FixVersionStrFunc: fixVersionStr,
	}
	res, err := latest.Check(githubTag, version)
	if err != nil {
		fmt.Println(err)
		return
	}
	if res.Outdated {
		fmt.Printf("%s is not latest, you should upgrade to %s\n", version, res.Current)
	}
}
