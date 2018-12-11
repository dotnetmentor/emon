package main

import (
	"fmt"
	"net/url"
	"strings"
)

func cleanURL(urlString string) string {
	uri, _ := url.Parse(urlString)
	p, _ := uri.User.Password()
	pf := fmt.Sprintf(":%s@", p)
	urlString = strings.Replace(urlString, pf, ":****@", -1)
	return urlString
}
