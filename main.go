package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jamiealquiza/envy"
	dd "github.com/zorkian/go-datadog-api"
)

var Params struct {
	APIKey string
	AppKey string
	Query  string
	Span   int
	ByTags []string
}

func init() {
	flag.StringVar(&Params.APIKey, "api-key", "", "Datadog API key")
	flag.StringVar(&Params.AppKey, "app-key", "", "Datadog app key")
	q := flag.String("query", "avg:system.load.1{*}", "Datadog metric query")
	flag.IntVar(&Params.Span, "span", 300, "Query range in seconds (now - span)")
	tags := flag.String("by-tags", "host", "Metric tags to reference data by (comma delimited)")

	envy.Parse("DDQ")
	flag.Parse()

	// Get tags.
	Params.ByTags = strings.Split(*tags, ",")

	// Complete query string.
	var b bytes.Buffer
	b.WriteString(*q)
	b.WriteString(fmt.Sprintf(".rollup(avg, %d)", Params.Span))
	Params.Query = b.String()
}

func main() {
	// Init, validate.
	client := dd.NewClient(Params.APIKey, Params.AppKey)

	ok, err := client.Validate()
	exitOnErr(err)

	if !ok {
		fmt.Println("invalid API or app key")
		os.Exit(1)
	}

	// Query.
	fmt.Printf("submitting %s\n\n", Params.Query)

	start := time.Now().Add(-time.Duration(Params.Span) * time.Second).Unix()
	o, err := client.QueryMetrics(start, time.Now().Unix(), Params.Query)
	exitOnErr(err)

	// Parse/print.
	var t bytes.Buffer
	for _, ts := range o {
		// Get the value for each tag, append
		// to our output key.
		for i, tag := range Params.ByTags {
			v := tagValFromScope(ts.GetScope(), tag)

			// If the tag isn't present.
			if v == "" {
				v = "null"
			}
			t.WriteString(v)
			if i < len(Params.ByTags)-1 {
				t.WriteString(",")
			}
		}

		// Write key, metric value.
		fmt.Printf("%20s: %.2f\n", t.String(), *ts.Points[0][1])
		t.Reset()
	}
}

// tagValFromScope takes a metric scope string
// and a tag and returns that tag's value.
func tagValFromScope(scope, tag string) string {
	ts := strings.Split(scope, ",")

	return valFromTags(ts, tag)
}

// valFromTags takes a []string of tags and
// a key, returning the value for the key.
func valFromTags(tags []string, key string) string {
	var v []string

	for _, tag := range tags {
		if strings.HasPrefix(tag, key+":") {
			v = strings.Split(tag, ":")
			break
		}
	}

	if len(v) > 1 {
		return v[1]
	}

	return ""
}

func exitOnErr(e error) {
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
}
