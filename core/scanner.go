package core

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"httpx.mrmonsif/pkg/technologies"
	"httpx.mrmonsif/pkg/security"
)

type Result struct {
	URL           string
	StatusCode    int
	Title         string
	Server        string
	ContentType   string
	ContentLength int64
	Headers       map[string]string
	Technologies  []string
	SecurityInfo  map[string]string
	ResponseTime  time.Duration
	IP            string
	Port          int
}

type Scanner struct {
	Threads         int
	Timeout         int
	TechDetection   bool
	SecurityScan    bool
	Ports           string
	FollowRedirects bool
}

func (s *Scanner) Scan(input *os.File) ([]Result, error) {
	var results []Result
	resultChan := make(chan Result)
	var wg sync.WaitGroup

	// Process results
	go func() {
		for result := range resultChan {
			results = append(results, result)
		}
	}()

	scanner := bufio.NewScanner(input)
	limiter := make(chan struct{}, s.Threads)

	for scanner.Scan() {
		domain := strings.TrimSpace(scanner.Text())
		domain = strings.TrimPrefix(domain, "https://")
		domain = strings.TrimPrefix(domain, "http://")
		domain = strings.TrimSuffix(domain, "/")

		if domain == "" {
			continue
		}

		ports := s.parsePorts()
		for _, port := range ports {
			wg.Add(1)
			limiter <- struct{}{}

			go func(d string, p int) {
				defer wg.Done()
				defer func() { <-limiter }()

				result := s.scanHost(d, p)
				if result.StatusCode > 0 {
					resultChan <- result
				}
			}(domain, port)
		}
	}

	wg.Wait()
	close(resultChan)
	return results, nil
}

func (s *Scanner) scanHost(domain string, port int) Result {
	start := time.Now()

	protocol := "https"
	if port == 80 {
		protocol = "http"
	}

	url := fmt.Sprintf("%s://%s:%d", protocol, domain, port)

	client := &http.Client{
		Timeout: time.Duration(s.Timeout) * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if s.FollowRedirects {
				return nil
			}
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Get(url)
	if err != nil {
		return Result{}
	}
	defer resp.Body.Close()

	result := Result{
		URL:           url,
		StatusCode:    resp.StatusCode,
		Server:        resp.Header.Get("Server"),
		ContentType:   resp.Header.Get("Content-Type"),
		ContentLength: resp.ContentLength,
		Headers:       make(map[string]string),
		ResponseTime:  time.Since(start),
		Port:          port,
	}

	for k, v := range resp.Header {
		if len(v) > 0 {
			result.Headers[k] = v[0]
		}
	}

	if s.TechDetection {
		result.Technologies = technologies.Detect(resp)
	}

	if s.SecurityScan {
		result.SecurityInfo = security.AnalyzeHeaders(resp.Header)
	}

	if result.SecurityInfo == nil {
		result.SecurityInfo = make(map[string]string)
	}
	serverLower := strings.ToLower(result.Server)
	switch {
	case strings.Contains(serverLower, "cloudflare"):
		result.SecurityInfo["WAF"] = "Cloudflare"
	case strings.Contains(serverLower, "akamai"):
		result.SecurityInfo["WAF"] = "Akamai"
	case strings.Contains(serverLower, "imperva"):
		result.SecurityInfo["WAF"] = "Imperva"
	}

	suspicious := ""
	if result.StatusCode == 403 || result.StatusCode == 503 {
		suspicious = "[🔥 Suspicious] "
	}

	statusColor := color.New(color.FgWhite)
	switch {
	case result.StatusCode >= 200 && result.StatusCode < 300:
		statusColor = color.New(color.FgGreen)
	case result.StatusCode >= 300 && result.StatusCode < 400:
		statusColor = color.New(color.FgBlue)
	case result.StatusCode == 403:
		statusColor = color.New(color.FgRed)
	case result.StatusCode >= 400 && result.StatusCode < 500:
		statusColor = color.New(color.FgYellow)
	case result.StatusCode >= 500:
		statusColor = color.New(color.FgHiRed)
	}

	statusColor.Printf("%s[%d] %s [%s] [Tech: %v] [Time: %v]\n",
		suspicious,
		result.StatusCode,
		result.URL,
		result.Server,
		result.Technologies,
		result.ResponseTime)

	return result
}

func (s *Scanner) parsePorts() []int {
	var ports []int
	for _, p := range strings.Split(s.Ports, ",") {
		if port, err := strconv.Atoi(strings.TrimSpace(p)); err == nil {
			ports = append(ports, port)
		}
	}
	return ports
}