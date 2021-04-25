package main

import (
	"flag"

	"github.com/jgvkmea/go-sort-tabelog/service"
	"github.com/jgvkmea/go-sort-tabelog/tabelog"
)

var (
	flagArea    string
	flagKeyword string
)

func main() {
	flag.StringVar(&flagArea, "area", "", "")
	flag.StringVar(&flagKeyword, "keyword", "", "")

	flag.Parse()
	params := tabelog.SearchParams{Area: flagArea, Keyword: flagKeyword}

	service.GetShopsOrderRating(params)
}
