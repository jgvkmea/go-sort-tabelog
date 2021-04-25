package service

import (
	"fmt"
	"sort"

	"github.com/jgvkmea/go-sort-tabelog/tabelog"
	"github.com/sirupsen/logrus"
)

const (
	outputCount = 7
)

var (
	log = logrus.New()
)

func GetShopsOrderRating(params tabelog.SearchParams) {
	driver := tabelog.NewWebDriver()
	if err := driver.Start(); err != nil {
		return
	}
	defer driver.Stop()

	page, err := tabelog.NewPage(driver)
	if err != nil {
		log.Errorf("failed to create page: %v", err)
		return
	}

	tabelog.GoToTabelogTop(page)
	err = tabelog.Search(page, params, log)
	if err != nil {
		return
	}

	shops, err := tabelog.GetShopList(page, log)
	if err != nil {
		return
	}
	sort.Slice(shops, func(i, j int) bool { return shops[i].Rating > shops[j].Rating })

	count := outputCount
	if len(shops) < outputCount {
		count = len(shops)
	}
	for i := 0; i < count; i++ {
		fmt.Printf("%d位 rating:%g %s URL: %s\n", i+1, float64(shops[i].Rating)/100, shops[i].Name, shops[i].Url)
	}
}