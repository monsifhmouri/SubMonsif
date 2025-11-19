package security

import (
	"net/http"
)

// AnalyzeHeaders 
func AnalyzeHeaders(headers http.Header) map[string]string {
	securityInfo := make(map[string]string)

	// 
	checks := []string{
		"Strict-Transport-Security",
		"X-Frame-Options",
		"X-XSS-Protection",
		"X-Content-Type-Options",
		"Content-Security-Policy",
		"Referrer-Policy",
		"Permissions-Policy",
	}

	for _, h := range checks {
		if value := headers.Get(h); value != "" {
			securityInfo[h] = value
		} else {
			securityInfo[h] = "Missing"
		}
	}

	return securityInfo
}
