package detector

import (
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func DetectGoquery(url string) ([]string, error) {
	html, err := fetchHTML(url)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	var f []string

	// React
	if doc.Find("[data-reactroot], [data-reactid]").Length() > 0 {
		f = append(f, "React.js")
	}

	// Next.js
	if doc.Find("script#__NEXT_DATA__").Length() > 0 {
		f = append(f, "Next.js")
	}

	// Gatsby
	if doc.Find("#___gatsby").Length() > 0 {
		f = append(f, "Gatsby.js")
	}

	// AngularJS (v1)
	if doc.Find(".ng-binding, [ng-app], [data-ng-app], [ng-controller], [data-ng-controller], [ng-repeat], [data-ng-repeat]").Length() > 0 ||
		doc.Find("script[src*='angular.js'], script[src*='angular.min.js']").Length() > 0 {
		f = append(f, "Angular.js")
	}

	// Angular (v2+)
	if doc.Find("[ng-version]").Length() > 0 {
		f = append(f, "Angular")
	}

	// Vue
	if doc.Find("[data-v-app], [v-cloak]").Length() > 0 {
		f = append(f, "Vue.js")
	}

	// Svelte / SvelteKit
	if doc.Find("[data-svelte-h], sveltekit-endpoint, sveltekit-app").Length() > 0 {
		f = append(f, "Svelte.js/SvelteKit")
	}

	// fq.js
	if doc.Find("script[src*='fq.js'], #fq-root").Length() > 0 {
		f = append(f, "fq.js")
	}

	// Backbone / Ember / Meteor — script src fallback
	if doc.Find("script[src*='backbone.js'], script[src*='backbone.min.js']").Length() > 0 {
		f = append(f, "Backbone.js")
	}
	if doc.Find("script[src*='ember.js'], script[src*='ember.min.js']").Length() > 0 {
		f = append(f, "Ember.js")
	}
	if doc.Find("script[src*='meteor.js']").Length() > 0 {
		f = append(f, "Meteor.js")
	}
	if doc.Find("script[src*='jquery.js'], script[src*='jquery.min.js']").Length() > 0 {
		f = append(f, "jQuery.js")
	}
	if doc.Find("script[src*='vue.js'], script[src*='vue.min.js'], script[src*='vue.esm']").Length() > 0 {
		f = append(f, "Vue.js")
	}

	return f, nil
}

func fetchHTML(url string) (string, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}