package detector

import (
	"fmt"
	"net/url"
	"strings"
)

func Detect(rawURL string, cfg Config) Result {
	result := Result{URL: rawURL}

	// Validate URL
	if _, err := url.ParseRequestURI(rawURL); err != nil {
		result.Error = (&DetectionError{
			Kind:    "invalid_url",
			Message: fmt.Sprintf("invalid URL %q — make sure it starts with http:// or https://", rawURL),
		}).Error()
		return result
	}

	if cfg.UseChrome {
		frameworks, err := DetectChrome(rawURL)
		if err != nil {
			detErr := classifyError("chrome", err)
			if !cfg.UseGoquery {
				result.Error = detErr.Message
				result.Source = "chrome"
				return result
			}
			// Chrome failed but goquery also requested — fall through
		} else {
			result.Frameworks = dedupe(frameworks)
			result.Source = "chrome"
			return result
		}
	}

	if cfg.UseGoquery {
		frameworks, err := DetectGoquery(rawURL)
		if err != nil {
			detErr := classifyError("goquery", err)
			result.Error = detErr.Message
			result.Source = "goquery"
			return result
		}
		result.Frameworks = dedupe(frameworks)
		result.Source = "goquery"
		return result
	}

	return result
}

func classifyError(source string, err error) *DetectionError {
	msg := err.Error()

	switch {
	case strings.Contains(msg, "timeout") || strings.Contains(msg, "deadline exceeded"):
		return &DetectionError{
			Kind:    "timeout",
			Message: fmt.Sprintf("timed out while loading the page — the site may be too slow or unreachable"),
		}
	case strings.Contains(msg, "no such host") || strings.Contains(msg, "connection refused"):
		return &DetectionError{
			Kind:    "invalid_url",
			Message: fmt.Sprintf("could not reach the site — check the URL and your internet connection"),
		}
	case strings.Contains(msg, "exec") || strings.Contains(msg, "chrome") || strings.Contains(msg, "chromium"):
		return &DetectionError{
			Kind:    "chrome_failed",
			Message: "Chrome/Chromium not found — install it or use --goquery for static detection",
		}
	case source == "chrome":
		return &DetectionError{
			Kind:    "chrome_failed",
			Message: fmt.Sprintf("headless Chrome failed: %s", msg),
		}
	default:
		return &DetectionError{
			Kind:    source + "_failed",
			Message: fmt.Sprintf("detection failed: %s", msg),
		}
	}
}

func dedupe(in []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, v := range in {
		v = strings.TrimSpace(v)
		if v != "" && !seen[v] {
			seen[v] = true
			out = append(out, v)
		}
	}
	return out
}