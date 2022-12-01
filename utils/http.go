package utils

import (
	"errors"
	"github.com/muety/wakapi/models"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	cacheMaxAgePattern = `max-age=(\d+)`
)

var (
	cacheMaxAgeRe *regexp.Regexp
)

func init() {
	cacheMaxAgeRe = regexp.MustCompile(cacheMaxAgePattern)
}

func IsNoCache(r *http.Request, cacheTtl time.Duration) bool {
	cacheControl := r.Header.Get("cache-control")
	if strings.Contains(cacheControl, "no-cache") {
		return true
	}
	if match := cacheMaxAgeRe.FindStringSubmatch(cacheControl); match != nil && len(match) > 1 {
		if maxAge, _ := strconv.Atoi(match[1]); time.Duration(maxAge)*time.Second <= cacheTtl {
			return true
		}
	}
	return false
}

func ParsePageParams(r *http.Request) *models.PageParams {
	pageParams := &models.PageParams{}
	page := r.URL.Query().Get("page")
	pageSize := r.URL.Query().Get("page_size")
	if p, err := strconv.Atoi(page); err == nil {
		pageParams.Page = p
	}
	if p, err := strconv.Atoi(pageSize); err == nil && pageParams.Page > 0 {
		pageParams.PageSize = p
	}
	return pageParams
}

func ParsePageParamsWithDefault(r *http.Request, page, size int) *models.PageParams {
	pageParams := ParsePageParams(r)
	if pageParams.Page == 0 {
		pageParams.Page = page
	}
	if pageParams.PageSize == 0 {
		pageParams.PageSize = size
	}
	return pageParams
}

func ParseUserAgent(ua string) (string, string, error) {
	re := regexp.MustCompile(`(?iU)^wakatime\/(?:v?[\d+.]+|unset)\s\((\w+)-.*\)\s.+\s([^\/\s]+)-wakatime\/.+$`)
	groups := re.FindAllStringSubmatch(ua, -1)
	if len(groups) == 0 || len(groups[0]) != 3 {
		return "", "", errors.New("failed to parse user agent string")
	}
	return groups[0][1], groups[0][2], nil
}
