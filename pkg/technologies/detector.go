package technologies

import (
	"net/http"
	"strings"
)

func Detect(resp *http.Response) []string {
	var tech []string

	headers := map[string]string{}
	for k, v := range resp.Header {
		if len(v) > 0 {
			headers[strings.ToLower(k)] = strings.ToLower(v[0])
		}
	}

	// Server detection
	if server := headers["server"]; server != "" {
		switch {
		case strings.Contains(server, "nginx"):
			tech = append(tech, "Nginx")
		case strings.Contains(server, "apache"):
			tech = append(tech, "Apache")
		case strings.Contains(server, "iis"):
			tech = append(tech, "IIS")
		}
	}

	// Framework detection
	if poweredBy := headers["x-powered-by"]; poweredBy != "" {
		switch {
		case strings.Contains(poweredBy, "php"):
			tech = append(tech, "PHP")
		case strings.Contains(poweredBy, "asp.net"):
			tech = append(tech, "ASP.NET")
		case strings.Contains(poweredBy, "express"):
			tech = append(tech, "Express.js")
		}
	}

	// CMS detection via headers
	if headers["x-generator"] != "" {
		if strings.Contains(headers["x-generator"], "wordpress") {
			tech = append(tech, "WordPress")
		}
	}

	return tech
}