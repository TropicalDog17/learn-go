package trendcrawler

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

type Crawler struct {
}
type Repo struct {
	SpokenLanguage      string
	ProgrammingLanguage string
	StarCount           uint16
	ForkCount           uint16
	StarGained          uint16
	Description         string
	DateRange           string
	URL                 string
}

var (
	BASE_URL            = "https://github.com"
	GITHUB_TRENDING_URL = "https://github.com/trending/%s?since=%s&spoken_language_code=%s"
)

// Parse numerical text "num" to uint16 text, also handle comma-seperated number.
func parseCommaNumber(num string) (uint16, error) {
	num = strings.Replace(num, ",", "", -1)
	parsedNum, err := strconv.Atoi(num)
	if err != nil {
		return 0, err
	}
	return uint16(parsedNum), nil
}

// Parse a string in the format '<NUMBER> stars today/last week/last month.' and convert the numerical part into a uint16 value. If the initial portion of the string is not a valid numerical representation, return an error."
func parseStarGained(s string) (uint16, error) {
	s = strings.TrimSpace(s)
	// TODO: handle if empty string
	s = strings.Split(s, " ")[0]
	// Try to convert to uint
	count, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return uint16(count), nil
}
func (r Repo) String() string {
	return fmt.Sprintf("SpokenLanguage: %s\nProgrammingLanguage: %s\nStarCount: %d\nForkCount: %d\nStarGained: %d\nDescription: %s\nDateRange: %s\nURL: %s\n", r.SpokenLanguage, r.ProgrammingLanguage, r.StarCount, r.ForkCount, r.StarGained, r.Description, r.DateRange, r.URL)
}

func (r *Crawler) Crawl(spokenLang string, programmingLang string, dateRange string) ([]Repo, error) {
	repos := make([]Repo, 20)
	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: coursera.org, www.coursera.org
		colly.AllowedDomains("github.com", "www.github.com"),

		// Cache responses to prevent multiple download of pages
		// even if the collector is restarted
		colly.CacheDir("./github_cache"),
	)

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		log.Println("visiting", r.URL.String())
	})
	c.OnHTML("div.Box > div > article.Box-row", func(e *colly.HTMLElement) {
		r := Repo{}
		url := BASE_URL + e.ChildAttr("h2 > a", "href")

		starCountText := strings.TrimSpace(e.ChildText("div > span + a"))
		forkCountText := strings.TrimSpace(e.ChildText("div > span + a + a"))
		starGainedText := e.ChildText("div > span + a + a + span + span")
		description := e.ChildText("h2 + p")

		// StarGained: x stars today/last week/last month
		starGained, _ := parseStarGained(starGainedText)
		starCount, _ := parseCommaNumber(starCountText)
		forkCount, _ := parseCommaNumber(forkCountText)

		r.StarGained = starGained
		r.StarCount = starCount
		r.ForkCount = forkCount
		r.ProgrammingLanguage = programmingLang
		r.DateRange = dateRange
		r.Description = description
		r.SpokenLanguage = spokenLang
		r.URL = url
		repos = append(repos, r)
	})
	// Start scraping on Github Trending
	c.Visit(
		fmt.Sprintf(
			GITHUB_TRENDING_URL,
			programmingLang, dateRange, spokenLang,
		),
	)

	// Return the scraped repos
	return repos, nil
}
