package tabelog

import (
	"fmt"
	"strconv"

	"github.com/sclevine/agouti"
	"github.com/sirupsen/logrus"
)

const (
	userAgent     = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.13; rv:57.0) Gecko/20100101 Chrome/35.0.1916.114 Safari/537.36"
	tabelogTopURL = "https://tabelog.com/"
	maxSearchPage = 3
)

type SearchParams struct {
	Area    string
	Keyword string
}

type Shop struct {
	Name   string
	Rating int
	Url    string
}

func NewWebDriver() *agouti.WebDriver {
	return agouti.ChromeDriver(
		agouti.ChromeOptions("args", []string{
			"--headless",
			"no sandbox",
			fmt.Sprintf("--user-agent=%s", userAgent),
		}),
	)
}

func NewPage(driver *agouti.WebDriver) (*agouti.Page, error) {
	p, err := driver.NewPage()
	if err != nil {
		return nil, err
	}
	return p, nil
}

func GoToTabelogTop(page *agouti.Page) error {
	return page.Navigate(tabelogTopURL)
}

func Search(page *agouti.Page, params SearchParams, log *logrus.Logger) error {
	url, err := page.URL()
	if err != nil {
		log.Errorf("failed to get url: %v", err)
		return err
	}
	if url != tabelogTopURL {
		log.Errorln("this page is not tabelog top")
		return fmt.Errorf("this page is not tabelog top")
	}

	area := page.FindByID("sa")
	keyword := page.FindByID("sk")
	submit := page.FindByID("js-global-search-btn")

	err = area.Fill(params.Area)
	if err != nil {
		log.Errorf("failed to fill area: %v", err)
		return err
	}
	err = keyword.Fill(params.Keyword)
	if err != nil {
		log.Errorf("failed to fill keyword: %v", err)
		return err
	}
	err = submit.Click()
	if err != nil {
		log.Errorf("failed to click submit: %v", err)
		return err
	}

	return nil
}

func GetShopList(page *agouti.Page, log *logrus.Logger) ([]Shop, error) {
	shops := []Shop{}
	var err error

	for i := 1; i <= maxSearchPage; i++ {
		if i > 1 {
			url, err := page.FindByClass("c-pagination__arrow--next").Attribute("href")
			if err != nil {
				log.Errorln("failed to go next page URL: ", err)
				return nil, err
			}
			err = page.Navigate(url)
			if err != nil {
				log.Errorln("failed to go next page: ", err)
				return nil, err
			}
		}

		shops, err = getShopListFromPage(page, shops, log)
		if err != nil {
			log.Errorln("failed to get shop list from page: ", err)
			return nil, err
		}
	}

	return shops, nil
}

func getShopListFromPage(page *agouti.Page, shops []Shop, log *logrus.Logger) ([]Shop, error) {
	shop := Shop{}

	shopElements := page.All(".rstlist-info > .js-rst-cassette-wrap")
	shopCount, err := shopElements.Count()
	if err != nil {
		log.Errorln("failed to get shop count in page")
		return nil, err
	}

	for i := 0; i < shopCount; i++ {
		n, err := shopElements.At(i).FindByClass("list-rst__rst-name-target").Text()
		if err != nil {
			log.Errorf("failed to get shopname: %v", err)
			return nil, err
		}
		shop.Name = n

		if r, err := shopElements.At(i).FindByClass("list-rst__rating-val").Text(); err != nil {
			shop.Rating = 0
		} else {
			f, err := strconv.ParseFloat(r, 64)
			if err != nil {
				log.Errorf("failed to convert float: %v", err)
				return nil, err
			}
			shop.Rating = int(f * 100)
		}

		ref, err := shopElements.At(i).FindByClass("list-rst__rst-name-target").Attribute("href")
		if err != nil {
			log.Errorf("failed to get link: %v", err)
			return nil, err
		}
		shop.Url = ref

		shops = append(shops, shop)
	}

	return shops, nil
}
