package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"SubMonsif/core"
)

var (
	domain      = flag.String("d", "", "Target domain")
	domainsFile = flag.String("dl", "", "File containing domains")
	threads     = flag.Int("t", 100, "Threads")
	timeout     = flag.Int("timeout", 10, "Timeout in seconds")
	output      = flag.String("o", "", "Output file")
	verbose     = flag.Bool("v", false, "Verbose mode")
	bruteforce  = flag.Bool("brute", true, "Enable bruteforce")
	recursive   = flag.Bool("recursive", false, "Recursive subdomain discovery")
)

func banner() {
	fmt.Println(`
	╔══════════════════════════════════════════╗
	║               SubMonsif                  ║
	║            SUBDOMAIN ENGINE              ║
	║           Created by: MrMonsif           ║
	║     https://github.com/monsifhmouri      ║
	╚══════════════════════════════════════════╝
	`)
}

func main() {
	flag.Parse()
	banner()

	if *domain == "" && *domainsFile == "" {
		fmt.Println("Error: Please specify a domain or domains file")
		fmt.Println("Usage: SubMonsif -d example.com")
		fmt.Println("       SubMonsif -dl domains.txt")
		os.Exit(1)
	}

	engine := &core.Engine{
		Threads:    *threads,
		Timeout:    *timeout,
		Bruteforce: *bruteforce,
		Recursive:  *recursive,
		Verbose:    *verbose,
	}

	var results []string
	var err error

	if *domain != "" {
		results, err = engine.Discover(*domain)
	} else {
		results, err = engine.DiscoverFromFile(*domainsFile)
	}

	if err != nil {
		log.Fatal(err)
	}

	core.SaveResults(results, *output)
	fmt.Printf("\n[+] Found %d subdomains\n", len(results))
}
