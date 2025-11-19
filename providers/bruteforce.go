package providers

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

func Bruteforce(domain string, threads int) []string {
	var subdomains []string
	var mutex sync.Mutex
	var wg sync.WaitGroup

	wordlist := loadWordlist()
	jobs := make(chan string, len(wordlist))

	for i := 0; i < threads; i++ {
		wg.Add(1)
		go worker(jobs, domain, &subdomains, &mutex, &wg)
	}

	for _, word := range wordlist {
		jobs <- word
	}
	close(jobs)

	wg.Wait()
	return subdomains
}

func worker(jobs <-chan string, domain string, results *[]string, mutex *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	for word := range jobs {
		subdomain := fmt.Sprintf("%s.%s", word, domain)
		if checkSubdomain(subdomain) {
			mutex.Lock()
			*results = append(*results, subdomain)
			fmt.Printf("[Brute] Found: %s\n", subdomain)
			mutex.Unlock()
		}
	}
}

func checkSubdomain(subdomain string) bool {
	_, err := net.LookupHost(subdomain)
	return err == nil
}

func loadWordlist() []string {
	defaultWords := []string{
		"www", "api", "mail", "ftp", "cpanel", "webmail", "admin", "blog",
		"shop", "dev", "test", "staging", "ns1", "ns2", "cdn", "assets",
		"static", "media", "img", "images", "video", "download", "portal",
		"secure", "login", "dashboard", "app", "apps", "mobile", "m",
		"support", "help", "docs", "wiki", "forum", "community", "news",
		"events", "calendar", "files", "share", "cloud", "storage", "backup",
		"db", "database", "sql", "mysql", "oracle", "redis", "cache", "proxy",
		"vpn", "remote", "ssh", "sftp", "git", "svn", "jenkins", "ci", "cd",
		"docker", "k8s", "kubernetes", "monitor", "metrics", "grafana",
		"prometheus", "alert", "alerts", "status", "health", "ping",
	}

	file, err := os.Open("wordlists/subdomains.txt")
	if err != nil {
		return defaultWords
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words = append(words, strings.TrimSpace(scanner.Text()))
	}

	return words
}
