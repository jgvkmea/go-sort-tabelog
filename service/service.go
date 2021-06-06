package service

import (
	"sort"

	"github.com/jgvkmea/go-sort-tabelog/tabelog"
	"github.com/sirupsen/logrus"
)

const (
	defaultOutputCount = 7
)

var (
	log = logrus.New()
)

func GetShopsOrderByRating(area string, keyword string) ([]tabelog.Shop, error) {
	driver := tabelog.NewWebDriver()
	if err := driver.Start(); err != nil {
		return nil, err
	}
	defer driver.Stop()

	page, err := tabelog.NewPage(driver)
	if err != nil {
		log.Errorf("failed to create page: %v", err)
		return nil, err
	}

	tabelog.GoToTabelogTop(page)
	err = tabelog.Search(page, area, keyword, log)
	if err != nil {
		return nil, err
	}

	shops, err := tabelog.GetShopList(page, log)
	if err != nil {
		return nil, err
	}
	sort.Slice(shops, func(i, j int) bool { return shops[i].Rating > shops[j].Rating })

	return shops, nil
}

func getOutputCount(shops []tabelog.Shop) (count int) {
	if len(shops) < defaultOutputCount {
		return len(shops)
	}
	return defaultOutputCount
}
