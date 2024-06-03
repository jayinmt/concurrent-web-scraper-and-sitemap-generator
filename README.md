# Concurrent Web Scraper and Sitemap Generator

This Go script is a concurrent web scraper and sitemap generator. It takes a website URL as input and concurrently scrapes the website to generate a sitemap. The program utilizes Go's concurrency features, such as goroutines and channels, to efficiently crawl and process multiple pages simultaneously.

## Features

- Concurrent web scraping using goroutines and channels
- Recursive crawling to discover and process all pages of the website
- URL normalization and duplicate URL filtering
- Respect for robots.txt rules
- Generating a structured sitemap in XML format
- Error handling and reporting
- Command-line interface for easy execution

## Prerequisites

- Go programming language (version 1.16 or higher)

## Installation

1. Clone the repository:

   ```
   git clone https://github.com/jaydxyz/concurrent-web-scraper-and-sitemap-generator.git
   ```

2. Navigate to the project directory:

   ```
   cd concurrent-web-scraper
   ```

3. Build the executable:

   ```
   go build sitemap_generator.go
   ```

## Usage

To generate a sitemap for a website, run the following command:

```
./sitemap_generator -website "https://example.com"
```

Replace `"https://example.com"` with the URL of the website you want to scrape and generate a sitemap for.

The generated sitemap will be printed to the console in XML format.

## Example

```
./sitemap_generator -website "https://example.com"
```

Output:
```xml
<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url>
    <loc>https://example.com/</loc>
  </url>
  <url>
    <loc>https://example.com/about</loc>
  </url>
  <url>
    <loc>https://example.com/contact</loc>
  </url>
  ...
</urlset>
```

## Contributing

Contributions are welcome! If you find any issues or have suggestions for improvements, please open an issue or submit a pull request.

## License

This project is licensed under the [MIT License](LICENSE).

## Acknowledgements

- The script utilizes the `robotstxt` package by temoto: [https://github.com/temoto/robotstxt](https://github.com/temoto/robotstxt)
