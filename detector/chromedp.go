package detector

import (
	"context"
	"time"

	"github.com/chromedp/chromedp"
)

func DetectChrome(url string) ([]string, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Headless,
		chromedp.DisableGPU,
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 25*time.Second)
	defer cancel()

	js := `(function() {
		const f = [];

		// React
		if (
			!!window.React ||
			!!document.querySelector("[data-reactroot], [data-reactid]") ||
			Array.from(document.querySelectorAll("*")).some(function(e) {
				return e._reactRootContainer !== undefined ||
					Object.keys(e).some(function(k) {
						return k.startsWith("__reactContainer");
					});
			})
		) f.push("React.js");

		// Next.js
		if (!!document.querySelector("script[id='__NEXT_DATA__']"))
			f.push("Next.js");

		// Gatsby
		if (!!document.querySelector("[id='___gatsby']"))
			f.push("Gatsby.js");

		// AngularJS (v1)
		if (
			!!window.angular ||
			!!document.querySelector(".ng-binding, [ng-app], [data-ng-app], [ng-controller], [data-ng-controller], [ng-repeat], [data-ng-repeat]") ||
			!!document.querySelector("script[src*='angular.js'], script[src*='angular.min.js']")
		) f.push("Angular.js");

		// Angular (v2+)
		if (
			!!window.getAllAngularRootElements ||
			!!(window.ng && window.ng.coreTokens && window.ng.coreTokens.NgZone)
		) f.push("Angular");

		// fq.js
		if (
			!!document.querySelector("script[src*='fq.js']") ||
			!!document.querySelector("#fq-root")
		) f.push("fq.js");

		// Svelte / SvelteKit
		if (
			!!document.querySelector("[data-svelte-h]") ||
			!!document.querySelector("sveltekit-endpoint, sveltekit-app")
		) f.push("Svelte.js/SvelteKit");

		// Window globals
		if (!!window.Backbone) f.push("Backbone.js");
		if (!!window.Ember)    f.push("Ember.js");
		if (!!window.Vue)      f.push("Vue.js");
		if (!!window.Meteor)   f.push("Meteor.js");
		if (!!window.Zepto)    f.push("Zepto.js");
		if (!!window.jQuery)   f.push("jQuery.js");
		if (!!window.can)      f.push("can.js");

		return f;
	})()`

	var frameworks []string

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.Sleep(2*time.Second),
		chromedp.Evaluate(js, &frameworks),
	)
	if err != nil {
		return nil, err
	}

	return frameworks, nil
}