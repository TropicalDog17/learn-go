package trendcrawler

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCommaNumber(t *testing.T) {
	num, err := parseCommaNumber("123,45")
	if assert.NoError(t, err) {
		assert.Equal(t, uint16(12345), num)
	}
	num, err = parseCommaNumber("123")
	if assert.NoError(t, err) {
		assert.Equal(t, uint16(123), num)

	}
}
func TestStarGained(t *testing.T) {
	count, err := parseStarGained("1222 stars today")
	if assert.NoError(t, err) {
		assert.Equal(t, uint16(1222), count)
	}
	count, err = parseStarGained("120 stars last week")
	if assert.NoError(t, err) {
		assert.Equal(t, uint16(120), count)
	}
	count, err = parseStarGained("     150 stars last month     ")
	if assert.NoError(t, err) {
		assert.Equal(t, uint16(150), count)
	}
	_, err = parseCommaNumber("     zero stars last month     ")
	assert.Error(t, err)
}

func TestCrawl(t *testing.T) {
	crawler := Crawler{}
	repos, _ := crawler.Crawl("en", "golang", "daily")
	fmt.Print(repos)

}
