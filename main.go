package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"jsframework-detector/detector"
)

// ANSI colors
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorBold   = "\033[1m"
	colorDim    = "\033[2m"
)

func main() {
	useGoquery := flag.Bool("goquery", false, "Use goquery static HTML detection instead of headless Chrome")
	useJSON    := flag.Bool("json", false, "Output results as JSON")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%sUsage:%s jsdetect [flags] <url>\n\n", colorBold, colorReset)
		fmt.Fprintf(os.Stderr, "%sFlags:%s\n", colorBold, colorReset)
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n%sExamples:%s\n", colorBold, colorReset)
		fmt.Fprintf(os.Stderr, "  jsdetect https://vercel.com\n")
		fmt.Fprintf(os.Stderr, "  jsdetect --goquery https://vercel.com\n")
		fmt.Fprintf(os.Stderr, "  jsdetect --json https://vercel.com\n")
	}

	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		
	}

	rawURL := flag.Arg(0)

	cfg := detector.Config{
		UseChrome:  !*useGoquery,
		UseGoquery: *useGoquery,
	}

	// Show scanning message (skip if JSON)
	if !*useJSON {
		source := "Chrome"
		if *useGoquery {
			source = "goquery"
		}
		fmt.Printf("\n%s⚡ Scanning%s %s%s%s %s[via %s]%s\n\n",
			colorBold, colorReset,
			colorCyan, rawURL, colorReset,
			colorDim, source, colorReset,
		)
	}

	result := detector.Detect(rawURL, cfg)

	if *useJSON {
		printJSON(result)
		return
	}

	printPretty(result)
}

func printJSON(result detector.Result) {
	out, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s✗ Failed to encode result: %v%s\n", colorRed, err, colorReset)
		
	}
	fmt.Println(string(out))
}

func printPretty(result detector.Result) {
	// Error state
	if result.Error != "" {
		fmt.Printf("%s✗ Error:%s %s\n\n", colorRed+colorBold, colorReset, result.Error)
		
	}

	// No frameworks found
	if len(result.Frameworks) == 0 {
		fmt.Printf("%s● No frameworks detected%s %s[source: %s]%s\n\n",
			colorYellow, colorReset,
			colorDim, result.Source, colorReset,
		)
		return
	}

	// Results header
	fmt.Printf("%s✓ Detected %d framework(s)%s %s[source: %s]%s\n\n",
		colorGreen+colorBold, len(result.Frameworks), colorReset,
		colorDim, result.Source, colorReset,
	)

	for _, fw := range result.Frameworks {
		fmt.Printf("  %s→%s %s%s%s\n", colorCyan, colorReset, colorBold, fw, colorReset)
	}

	fmt.Println()
}