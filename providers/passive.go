package providers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func PassiveDiscovery(domain string) []string {
	var allSubdomains []string
	allSubdomains = append(allSubdomains, getFromCRTsh(domain)...)
	allSubdomains = append(allSubdomains, getFromHackerTarget(domain)...)
	allSubdomains = append(allSubdomains, getFromVirusTotal(domain)...)
	allSubdomains = append(allSubdomains, getFromOTX(domain)...)
	allSubdomains = append(allSubdomains, getFromShodan(domain)...)
	return removeDuplicates(allSubdomains)
}

func getFromCRTsh(domain string) []string {
	url := fmt.Sprintf("https://crt.sh/?q=%%25.%s&output=json", domain)
	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return []string{}
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var results []map[string]interface{}
	json.Unmarshal(body, &results)

	subMap := make(map[string]struct{})
	for _, entry := range results {
		if nameVal, ok := entry["name_value"].(string); ok {
			for _, sub := range strings.Split(nameVal, "\n") {
				s := strings.TrimSpace(sub)
				if s != "" && !strings.Contains(s, "*") {
					subMap[s] = struct{}{}
				}
			}
		}
	}
	var out []string
	for sub := range subMap {
		out = append(out, sub)
	}
	return out
}

func getFromHackerTarget(domain string) []string {
	url := fmt.Sprintf("https://api.hackertarget.com/hostsearch/?q=%s", domain)
	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return []string{}
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	lines := strings.Split(string(body), "\n")
	var subs []string
	for _, line := range lines {
		parts := strings.Split(line, ",")
		if len(parts) > 0 {
			s := strings.TrimSpace(parts[0])
			if s != "" {
				subs = append(subs, s)
			}
		}
	}
	return subs
}

func getFromVirusTotal(domain string) []string {
	apiKey := "............."
	url := fmt.Sprintf("https://www.virustotal.com/api/v3/domains/%s/subdomains?limit=40", domain)
	client := &http.Client{Timeout: 20 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []string{}
	}
	req.Header.Set("x-apikey", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return []string{}
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var parsed struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	json.Unmarshal(body, &parsed)

	var subs []string
	for _, item := range parsed.Data {
		s := strings.TrimSpace(item.ID)
		if s != "" {
			subs = append(subs, s)
		}
	}
	return subs
}

func getFromOTX(domain string) []string {
	apiKey := ".............."
	url := fmt.Sprintf("https://otx.alienvault.com/api/v1/indicators/domain/%s/passive_dns", domain)
	client := &http.Client{Timeout: 20 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []string{}
	}
	req.Header.Set("X-OTX-API-KEY", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return []string{}
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var parsed struct {
		PassiveDNS []struct {
			Hostname string `json:"hostname"`
		} `json:"passive_dns"`
	}
	json.Unmarshal(body, &parsed)

	var subs []string
	for _, item := range parsed.PassiveDNS {
		s := strings.TrimSpace(item.Hostname)
		if s != "" {
			subs = append(subs, s)
		}
	}
	return subs
}

func getFromShodan(domain string) []string {
	apiKey := "................"
	url := fmt.Sprintf("https://api.shodan.io/dns/domain/%s?key=%s", domain, apiKey)
	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return []string{}
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var parsed struct {
		Subdomains []string `json:"subdomains"`
	}
	json.Unmarshal(body, &parsed)

	var subs []string
	for _, sub := range parsed.Subdomains {
		s := fmt.Sprintf("%s.%s", sub, domain)
		subs = append(subs, s)
	}
	return subs
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
