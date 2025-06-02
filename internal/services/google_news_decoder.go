package services

import (
	"encoding/base64"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// GoogleNewsDecoder handles decoding of Google News encoded URLs
type GoogleNewsDecoder struct {
	client *http.Client
}

// NewGoogleNewsDecoder creates a new Google News URL decoder
func NewGoogleNewsDecoder() *GoogleNewsDecoder {
	return &GoogleNewsDecoder{
		client: &http.Client{
			Timeout: 10 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				// Don't follow redirects, we want to capture the redirect URL
				return http.ErrUseLastResponse
			},
		},
	}
}

// DecodeGoogleNewsURL attempts to decode a Google News URL to get the original article URL
func (gnd *GoogleNewsDecoder) DecodeGoogleNewsURL(encodedURL string) (string, error) {
	// Check if this is actually a Google News URL
	if !strings.Contains(encodedURL, "news.google.com") {
		return encodedURL, nil
	}

	// Try different decoding methods
	if decoded := gnd.tryBase64Decode(encodedURL); decoded != "" {
		return decoded, nil
	}

	if decoded := gnd.tryHTTPRedirect(encodedURL); decoded != "" {
		return decoded, nil
	}

	// If all decoding methods fail, return the original URL
	log.Printf("Could not decode Google News URL: %s", encodedURL)
	return encodedURL, nil
}

// tryBase64Decode attempts to decode the URL using base64 decoding method
func (gnd *GoogleNewsDecoder) tryBase64Decode(encodedURL string) string {
	// Extract the encoded part from URLs like:
	// https://news.google.com/rss/articles/CBMi...?oc=5
	urlParts := strings.Split(encodedURL, "/articles/")
	if len(urlParts) < 2 {
		return ""
	}

	// Get the encoded part (everything after /articles/ and before ?)
	encodedPart := strings.Split(urlParts[1], "?")[0]

	// Try to decode the base64 encoded part
	decoded, err := base64.URLEncoding.DecodeString(encodedPart)
	if err != nil {
		// Try standard base64 decoding
		decoded, err = base64.StdEncoding.DecodeString(encodedPart)
		if err != nil {
			return ""
		}
	}

	// The decoded data contains the original URL
	// Look for URL patterns in the decoded data
	urlPattern := regexp.MustCompile(`https?://[^\s<>"]+`)
	matches := urlPattern.FindAllString(string(decoded), -1)

	for _, match := range matches {
		// Skip Google URLs
		if !strings.Contains(match, "google.com") && strings.Contains(match, "http") {
			return strings.TrimRight(match, `"'\s`)
		}
	}

	return ""
}

// tryHTTPRedirect attempts to follow the redirect to get the original URL
func (gnd *GoogleNewsDecoder) tryHTTPRedirect(encodedURL string) string {
	resp, err := gnd.client.Head(encodedURL)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	// Check if we got a redirect response
	if resp.StatusCode >= 300 && resp.StatusCode < 400 {
		location := resp.Header.Get("Location")
		if location != "" && !strings.Contains(location, "google.com") {
			return location
		}
	}

	return ""
}

// BatchDecodeURLs decodes multiple Google News URLs concurrently
func (gnd *GoogleNewsDecoder) BatchDecodeURLs(urls []string) map[string]string {
	results := make(map[string]string)

	// For now, decode sequentially to avoid overwhelming Google's servers
	// In production, you might want to implement rate limiting and concurrent processing
	for _, url := range urls {
		if decoded, err := gnd.DecodeGoogleNewsURL(url); err == nil {
			results[url] = decoded
		} else {
			results[url] = url // fallback to original URL
		}

		// Small delay to be respectful to Google's servers
		time.Sleep(100 * time.Millisecond)
	}

	return results
}

// IsGoogleNewsURL checks if a URL is a Google News encoded URL
func IsGoogleNewsURL(url string) bool {
	return strings.Contains(url, "news.google.com") &&
		   (strings.Contains(url, "/articles/") || strings.Contains(url, "/read/"))
}

// ExtractSourceFromGoogleNewsURL attempts to extract the source domain from a Google News URL
func ExtractSourceFromGoogleNewsURL(encodedURL string) string {
	// This is a simplified version - in practice, you'd need more sophisticated parsing
	decoder := NewGoogleNewsDecoder()
	if decodedURL, err := decoder.DecodeGoogleNewsURL(encodedURL); err == nil {
		if parsedURL, err := url.Parse(decodedURL); err == nil {
			return parsedURL.Host
		}
	}
	return ""
}