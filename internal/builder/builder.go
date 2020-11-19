package builder

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/sergiorra/sitemap-builder-go/internal/parser"
)

const Xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

type Loc struct {
	Value string `xml:"loc"`
}

type Urlset struct {
	Urls  []Loc  `xml:"url"`
	Xmlns string `xml:"xmlns,attr"`
}

type empty struct{}

// Bfs searches all non-repeated pages within a max depth using the Breadth-First Search algorithm
func Bfs(urlStr string, maxDepth int) []string {
	seen := make(map[string]empty)
	var q map[string]empty
	nq := map[string]empty{
		urlStr: {},
	}
	for i := 0; i <= maxDepth; i++ {
		q, nq = nq, make(map[string]empty)
		if len(q) == 0 {
			break
		}
		for url, _ := range q {
			if _, ok := seen[url]; ok {
				continue
			}
			seen[url] = empty{}
			for _, link := range get(url) {
				nq[link] = empty{}
			}
		}
	}
	ret := make([]string, 0, len(seen))
	for url, _ := range seen {
		ret = append(ret, url)
	}
	return ret
}

// get fetches the data from an URL and returns all the links founded
func get(urlStr string) []string {
	resp, err := http.Get(urlStr)
	if err != nil {
		return []string{}
	}
	defer resp.Body.Close()
	reqUrl := resp.Request.URL
	baseUrl := &url.URL{
		Scheme: reqUrl.Scheme,
		Host:   reqUrl.Host,
	}
	base := baseUrl.String()
	return filter(hrefs(resp.Body, base), withPrefix(base))
}

// hrefs parses an HTML document and returns all formatted links in the HTML
func hrefs(r io.Reader, base string) []string {
	links, _ := parser.Parse(r)
	var ret []string
	for _, l := range links {
		switch {
		case strings.HasPrefix(l.Href, "/"):
			ret = append(ret, base+l.Href)
		case strings.HasPrefix(l.Href, "http"):
			ret = append(ret, l.Href)
		}
	}
	return ret
}

// filter filters all the links that meet the function received
func filter(links []string, keepFn func(string) bool) []string {
	var ret []string
	for _, link := range links {
		if keepFn(link) {
			ret = append(ret, link)
		}
	}
	return ret
}

// withPrefix checks that a link has a certain prefix using a closure to be more flexible
func withPrefix(pfx string) func(string) bool {
	return func(link string) bool {
		return strings.HasPrefix(link, pfx)
	}
}