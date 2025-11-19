package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"httpx.mrmonsif/core"
)

var (
	threads     = flag.Int("t", 100, "Threads")
	timeout     = flag.Int("timeout", 10, "Timeout in seconds")
	output      = flag.String("o", "", "Output file")
	verbose     = flag.Bool("v", false, "Verbose mode")
	techDetect  = flag.Bool("tech", true, "Technology detection")
	securityScan= flag.Bool("security", true, "Security headers scan")
	ports       = flag.String("p", "80,443,8080,8443", "Ports to scan")
	followRedirects = flag.Bool("fr", true, "Follow redirects")
)

func banner() {
	fmt.Println(`
	╔══════════════════════════════════════════╗
	║            httpx.mrmonsif                ║
	║             HTTP SCANNER                 ║
	║          Created by: MrMonsif            ║
	║     https://github.com/monsifhmouri      ║
	╚══════════════════════════════════════════╝
	`)
}

func main() {
	flag.Parse()
	banner()

	scanner := &core.Scanner{
		Threads:        *threads,
		Timeout:        *timeout,
		TechDetection:  *techDetect,
		SecurityScan:   *securityScan,
		Ports:          *ports,
		FollowRedirects: *followRedirects,
	}

	results, err := scanner.Scan(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	// Output results
	core.Output(results, *output, *verbose)
}