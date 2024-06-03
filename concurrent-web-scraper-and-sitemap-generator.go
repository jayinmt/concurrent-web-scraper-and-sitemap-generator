package main

import (
    "encoding/xml"
    "flag"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "os"
    "strings"
    "sync"

    "github.com/temoto/robotstxt"
)

type loc struct {
    Value string `xml:"loc"`
}

type urlset struct {
    Urls  []loc  `xml:"url"`
    Xmlns string `xml:"xmlns,attr"`
}

func main() {
    var website string
    flag.StringVar(&website, "website", "", "The website URL to scrape")
    flag.Parse()

    if website == "" {
        fmt.Println("Please provide a website URL using the -website flag")
        os.Exit(1)
    }

    sitemapUrls := make(chan string)
    var wg sync.WaitGroup

    wg.Add(1)
    go func() {
        defer wg.Done()
        scrapeWebsite(website, sitemapUrls)
    }()

    go func() {
        wg.Wait()
        close(sitemapUrls)
    }()

    var urls []loc
    for url := range sitemapUrls {
        urls = append(urls, loc{Value: url})
    }

    toXML := urlset{
        Urls:  urls,
        Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
    }

    fmt.Print(xml.Header)
    enc := xml.NewEncoder(os.Stdout)
    enc.Indent("", "  ")
    if err := enc.Encode(toXML); err != nil {
        fmt.Println("Error encoding XML:", err)
        os.Exit(1)
    }
}

func scrapeWebsite(website string, sitemapUrls chan<- string) {
    baseURL, err := url.Parse(website)
    if err != nil {
        fmt.Println("Error parsing website URL:", err)
        return
    }

    robotsURL := baseURL.Scheme + "://" + baseURL.Host + "/robots.txt"
    robotsData, err := http.Get(robotsURL)
    if err != nil {
        fmt.Println("Error retrieving robots.txt:", err)
        return
    }
    defer robotsData.Body.Close()

    robotsBytes, err := io.ReadAll(robotsData.Body)
    if err != nil {
        fmt.Println("Error reading robots.txt:", err)
        return
    }

    robotsGroup := robotstxt.Group(baseURL, robotsBytes)

    var wg sync.WaitGroup
    visited := make(map[string]bool)

    var crawl func(url string)
    crawl = func(url string) {
        defer wg.Done()

        if !robotsGroup.Test(url) {
            return
        }

        normalizedURL := normalizeURL(baseURL, url)
        if visited[normalizedURL] {
            return
        }
        visited[normalizedURL] = true

        resp, err := http.Get(url)
        if err != nil {
            fmt.Println("Error fetching URL:", err)
            return
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
            fmt.Printf("Skipping URL with status code %d: %s\n", resp.StatusCode, url)
            return
        }

        sitemapUrls <- normalizedURL

        if contentType := resp.Header.Get("Content-Type"); !strings.HasPrefix(contentType, "text/html") {
            return
        }

        body, err := io.ReadAll(resp.Body)
        if err != nil {
            fmt.Println("Error reading response body:", err)
            return
        }

        links := extractLinks(string(body))
        for _, link := range links {
            if isInternalLink(link, baseURL) {
                wg.Add(1)
                go crawl(link)
            }
        }
    }

    wg.Add(1)
    go crawl(baseURL.String())

    wg.Wait()
}

func normalizeURL(baseURL *url.URL, href string) string {
    link, err := url.Parse(href)
    if err != nil {
        return ""
    }
    return baseURL.ResolveReference(link).String()
}

func isInternalLink(link string, baseURL *url.URL) bool {
    parsedLink, err := url.Parse(link)
    if err != nil {
        return false
    }
    return parsedLink.Host == baseURL.Host
}

func extractLinks(body string) []string {
    var links []string
    for _, match := range linkRegex.FindAllString(body, -1) {
        link := strings.TrimSpace(match)
        link = strings.TrimPrefix(link, "href=\"")
        link = strings.TrimSuffix(link, "\"")
        links = append(links, link)
    }
    return links
}

var linkRegex = regexp.MustCompile(`href="([^"]+)"`)
