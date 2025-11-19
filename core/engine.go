package core

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"

	"SubMonsif/providers"
)

type Engine struct {
	Threads    int
	Timeout    int
	Bruteforce bool
	Recursive  bool
	Verbose    bool
}

func (e *Engine) Discover(domain string) ([]string, error) {
	var allSubdomains []string
	var mutex sync.Mutex
	var wg sync.WaitGroup

	passiveResults := make(chan string)
	wg.Add(1)
	go func() {
		defer wg.Done()
		subdomains := providers.PassiveDiscovery(domain)
		for _, sub := range subdomains {
			passiveResults <- sub
		}
	}()

	var bruteResults []string
	if e.Bruteforce {
		wg.Add(1)
		go func() {
			defer wg.Done()
			bruteResults = providers.Bruteforce(domain, e.Threads)
		}()
	}

	go func() {
		for sub := range passiveResults {
			mutex.Lock()
			if !contains(allSubdomains, sub) {
				allSubdomains = append(allSubdomains, sub)
				if e.Verbose {
					fmt.Printf("[+] Found: %s\n", sub)
				}
			}
			mutex.Unlock()
		}
	}()

	wg.Wait()
	close(passiveResults)

	mutex.Lock()
	allSubdomains = append(allSubdomains, bruteResults...)
	allSubdomains = removeDuplicates(allSubdomains)
	mutex.Unlock()

	return allSubdomains, nil
}

func (e *Engine) DiscoverFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var allResults []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		domain := strings.TrimSpace(scanner.Text())
		if domain != "" {
			results, err := e.Discover(domain)
			if err == nil {
				allResults = append(allResults, results...)
			}
		}
	}

	return removeDuplicates(allResults), nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func removeDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	var list []string
	for _, item := range slice {
		if _, value := keys[item]; !value {
			keys[item] = true
			list = append(list, item)
		}
	}
	return list
}
