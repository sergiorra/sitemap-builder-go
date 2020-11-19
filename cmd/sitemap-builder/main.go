package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"os"

	"github.com/sergiorra/sitemap-builder-go/internal/builder"
)

func main() {
	urlFlag := flag.String("url", "https://gophercises.com", "The url that you want to build a sitemap for")
	maxDepth := flag.Int("depth", 10, "The maximum number of links deep to traverse")
	flag.Parse()

	pages := builder.Bfs(*urlFlag, *maxDepth)
	toXml := builder.Urlset{
		Xmlns: builder.Xmlns,
	}
	for _, page := range pages {
		toXml.Urls = append(toXml.Urls, builder.Loc{page})
	}

	fmt.Print(xml.Header)
	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", "  ")
	if err := enc.Encode(toXml); err != nil {
		panic(err)
	}
	fmt.Println()
}
