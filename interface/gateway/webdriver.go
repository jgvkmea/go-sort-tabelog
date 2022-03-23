package gateway

import (
	"fmt"
	"math"
	"strconv"

	"github.com/jgvkmea/go-sort-tabelog/entity"
	"github.com/jgvkmea/go-sort-tabelog/usecase"
	"github.com/sclevine/agouti"
)

const (
	USER_AGENT          = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.13; rv:57.0) Gecko/20100101 Chrome/35.0.1916.114 Safari/537.36"
	TABELOG_TOP_URL     = "https://tabelog.com/"
	DEFAULT_SEARCH_PAGE = 5
	COUNT_PER_PAGE      = 20
)

type WebDriver struct {
	driver *agouti.WebDriver
	page   *agouti.Page
}

func (wd *WebDriver) GetShopList(area string, keyword string) (entity.Shops, error) {
	if err := wd.start(); err != nil {
		return []entity.Shop{}, err
	}
	defer wd.stop()

	if err := wd.newPage(); err != nil {
		return []entity.Shop{}, fmt.Errorf("failed to create page: %v", err)
	}
	if err := wd.goToTabelogTop(); err != nil {
		return []entity.Shop{}, fmt.Errorf("failed to go to tabelog page: %v", err)
	}
	if err := wd.search(area, keyword); err != nil {
		return []entity.Shop{}, fmt.Errorf("failed to search by area and keyword: %v", err)
	}

	shops := []entity.Shop{}
	var err error

	searchPage, err := wd.getSearchPageCount()
	if err != nil {
		return nil, err
	}

	for i := 1; i <= searchPage; i++ {
		// 1ページ目以外の場合は next ボタンを押してから店舗情報取得する
		if i > 1 {
			url, err := wd.page.FindByClass("c-pagination__arrow--next").Attribute("href")
			if err != nil {
				return nil, fmt.Errorf("failed to go next page URL: %v", err)
			}
			err = wd.page.Navigate(url)
			if err != nil {
				return nil, fmt.Errorf("failed to go next page: %v", err)
			}
		}

		shops, err = wd.getShopListFromPage(shops)
		if err != nil {
			return nil, fmt.Errorf("failed to get shop list from page: %v", err)
		}
	}

	return shops, nil
}

func NewWebDriver() usecase.ShopDataSource {
	return &WebDriver{
		driver: agouti.ChromeDriver(
			agouti.ChromeOptions("args", []string{
				"--headless",
				"--no-sandbox",
				fmt.Sprintf("--user-agent=%s", USER_AGENT),
			}),
		),
	}
}

func (d *WebDriver) start() error {
	return d.driver.Start()
}

func (d *WebDriver) stop() error {
	return d.driver.Stop()
}

func (d *WebDriver) newPage() error {
	p, err := d.driver.NewPage()
	if err != nil {
		return err
	}
	d.page = p
	return nil
}

func (d *WebDriver) navigate(url string) error {
	if d.page == nil {
		return fmt.Errorf("Driver.page is nil")
	}
	return d.page.Navigate(url)
}

func (d *WebDriver) goToTabelogTop() error {
	return d.navigate(TABELOG_TOP_URL)
}

func (d *WebDriver) search(area string, keyword string) error {
	url, err := d.page.URL()
	if err != nil {
		return fmt.Errorf("failed to get URL: %v", err)
	}
	if url != TABELOG_TOP_URL {
		return fmt.Errorf("current page is not tabelog")
	}

	err = d.page.FindByID("sa").Fill(area)
	if err != nil {
		return fmt.Errorf("failed to fill area: %v", err)
	}
	err = d.page.FindByID("sk").Fill(keyword)
	if err != nil {
		return fmt.Errorf("failed to fill keyword: %v", err)
	}
	err = d.page.FindByID("js-global-search-btn").Click()
	if err != nil {
		return fmt.Errorf("failed to click submit: %v", err)
	}
	return nil
}

func (d *WebDriver) getSearchPageCount() (int, error) {
	shopCountText, err := d.page.All(".c-page-count__num").At(2).Find("strong").Text()
	if err != nil {
		return 0, fmt.Errorf("failed to get shop count: %v", err)
	}

	shopCount, err := strconv.Atoi(shopCountText)
	if err != nil {
		return 0, fmt.Errorf("failed to convert string to int: %v", err)
	}

	shopPageCount := math.Ceil(float64(shopCount) / COUNT_PER_PAGE)
	if shopPageCount >= DEFAULT_SEARCH_PAGE {
		return DEFAULT_SEARCH_PAGE, nil
	}
	return int(shopPageCount), nil
}

func (d *WebDriver) getShopListFromPage(shops []entity.Shop) (entity.Shops, error) {
	shop := entity.Shop{}

	shopElements := d.page.All(".rstlist-info > .js-rst-cassette-wrap")
	shopCount, err := shopElements.Count()
	if err != nil {
		return nil, fmt.Errorf("failed to get shop count in page")
	}

	for i := 0; i < shopCount; i++ {
		n, err := shopElements.At(i).FindByClass("list-rst__rst-name-target").Text()
		if err != nil {
			return nil, fmt.Errorf("failed to get shopname: %v", err)
		}
		shop.SetName(n)

		if r, err := shopElements.At(i).FindByClass("list-rst__rating-val").Text(); err != nil {
			shop.SetRating(0)
		} else {
			f, err := strconv.ParseFloat(r, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to convert float: %v", err)
			}
			shop.SetRating(int(f * 100))
		}

		ref, err := shopElements.At(i).FindByClass("list-rst__rst-name-target").Attribute("href")
		if err != nil {
			return nil, fmt.Errorf("failed to get link: %v", err)
		}
		shop.SetURL(ref)

		shops = append(shops, shop)
	}

	return shops, nil
}
