package spn

import (
	"net/url"
	"reflect"
	"strconv"
)

type CaptureOptions struct {
	// Capture a web page with errors (HTTP status=4xx or 5xx). By default SPN2 captures only status=200 URLs.
	CaptureAll bool `spn:"capture_all"`
	// Capture web page outlinks automatically. This also applies to PDF, JSON, RSS and MRSS feeds.
	CaptureOutlinks int `spn:"capture_outlinks"`
	// Capture full page screenshot in PNG format. This is also stored in the Wayback Machine as a different capture.
	CaptureScreenshot bool `spn:"capture_screenshot"`
	// The capture becomes available in the Wayback Machine after ~12 hours instead of immediately. This option helps reduce the load on our systems. All API responses remain exactly the same when using this option.
	DelayWBAvailability bool `spn:"delay_wb_availability"`
	// Force the use of a simple HTTP GET request to capture the target URL. By default SPN2 does a HTTP HEAD on the target URL to decide whether to use a headless browser or a simple HTTP GET request. force_get overrides this behavior.
	ForceGet bool `spn:"force_get"`
	// Skip checking if a capture is a first if you don’t need this information. This will make captures run faster.
	SkipFirstArchive bool `spn:"skip_first_archive"`
	// if_not_archived_within=<timedelta>
	//
	// Capture web page only if the latest existing capture at the Archive is older than the <timedelta> limit.  Its  format could be any datetime expression like “3d 5h 20m” or just a number of seconds, e.g. “120”. If there is a capture within the defined timedelta, SPN2 returns that as a recent capture. The default system <timedelta> is 45 min.
	//
	// if_not_archived_within=<timedelta1>,<timedelta2>
	//
	// When using 2 comma separated <timedelta> values, the first one applies to the main capture and the second one applies to outlinks.
	IfNotArchivedWithin string `spn:"if_not_archived_within"`
	// Return the timestamp of the last capture for all outlinks.
	OutlinksAvailability bool `spn:"outlinks_availability"`
	// Send an email report of the captured URLs to the user’s email.
	EmailResult bool `spn:"email_result"`
	// Run JS code for <N> seconds after page load to trigger target page functionality like image loading on mouse over, scroll down to load more content, etc. The default system <N> is 5 sec.
	//
	// More details on the JS code we execute:
	// https://github.com/internetarchive/brozzler/blob/master/brozzler/behaviors.yaml
	//
	// WARNING: The max <N> value that applies is 30 sec.
	//
	// NOTE: If the target page doesn’t have any JS you need to run, you can use js_behavior_timeout=0 to speed up the capture.
	JsBehaviorTimeout string `spn:"js_behavior_timeout"` // It's hard to determine if int 0 is user input or default value, so we use string instead
	// Use extra HTTP Cookie value when capturing the target page.
	CaptureCookie string `spn:"capture_cookie"`
	// Use custom HTTP User-Agent value when capturing the target page.
	UseUserAgent string `spn:"use_user_agent"`

	// target_username=<XXX>
	// target_password=<YYY>
	//
	// Use your own username and password in the target page’s login forms.
	TargetUsername string `spn:"target_username"`
	// target_username=<XXX>
	// target_password=<YYY>
	//
	// Use your own username and password in the target page’s login forms.
	TargetPassword string `spn:"target_password"`
}

// converts CaptureOptions to url.Values
func (opts CaptureOptions) Encode() url.Values {
	urlValues := url.Values{}
	typ := reflect.TypeOf(opts)
	val := reflect.ValueOf(opts)
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)
		if value.Kind() == reflect.Bool {
			if value.Bool() {
				urlValues.Add(field.Tag.Get("spn"), "1")
			}
		} else if value.Kind() == reflect.String {
			if value.String() != "" {
				urlValues.Add(field.Tag.Get("spn"), value.String())
			}
		} else if value.Kind() == reflect.Int {
			urlValues.Add(field.Tag.Get("spn"), strconv.FormatInt(value.Int(), 10))
		} else {
			panic("Unknown field type")
		}
	}

	return urlValues
}
