package detector

type Config struct {
	UseChrome  bool
	UseGoquery bool
}

type Result struct {
	URL        string   `json:"url"`
	Frameworks []string `json:"frameworks"`
	Source     string   `json:"source"`
	Error      string   `json:"error,omitempty"`
}

type DetectionError struct {
	Kind    string // "invalid_url", "chrome_failed", "goquery_failed", "timeout"
	Message string
}

func (e *DetectionError) Error() string {
	return e.Message
}