package entity

import (
	"sort"
)

const (
	DEFAULT_OUTPUT_COUNT = 7
)

type Shop struct {
	name   string
	rating int
	url    string
}

type Shops []Shop

func (s *Shop) SetName(n string) {
	s.name = n
	return
}

func (s *Shop) GetName() string {
	return s.name
}

func (s *Shop) SetRating(r int) {
	s.rating = r
	return
}

func (s *Shop) GetRating() int {
	return s.rating
}

func (s *Shop) SetURL(u string) {
	s.url = u
	return
}

func (s *Shop) GetURL() string {
	return s.url
}

func (ss *Shops) SortByRating() {
	sort.Slice(*ss, func(i, j int) bool { return (*ss)[i].GetRating() > (*ss)[j].GetRating() })
}

func (ss *Shops) GetOutputCount() (count int) {
	if len(*ss) < DEFAULT_OUTPUT_COUNT {
		return len(*ss)
	}
	return DEFAULT_OUTPUT_COUNT
}

func (ss *Shops) GetTopShopList(count int) Shops {
	var topShopList Shops
	for i := 0; i < count; i++ {
		if len(*ss) <= i {
			break
		}
		topShopList = append(topShopList, (*ss)[i])
	}
	return topShopList
}
