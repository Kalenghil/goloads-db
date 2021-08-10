package main

import (
	"math/rand"
)

type User struct {
	Firstname string  `json:"firstname"`
	Lastname  string  `json:"lastname"`
	ID        int     `json:"id"`
	Account   int     `json:"account"`
	Token     string  `json:"token"`
	Money     float64 `json:"money"`
}

type UserStorage struct {
	UserMap map[string]User
}

type Banner struct {
	BannerID    string   `json:"id"`
	Image       string   `json:"image"`
	DomainURL   string   `json:"url"`
	Domains     []string `json:"domains"`
	ImageBase64 bool     `json:"image-base64"`
}

type Analytics struct {
	BannerID     string `json:"id"`
	Clicks       []int  `json:"clicks"`
	UniqueClicks []int  `json:"unique_clicks"`
	Views        []int  `json:"views"`
	UniqueViews  []int  `json:"unique_views"`
}

type BannerStorage struct {
	BannerMap map[string]Banner
}

type AnalyticsStorage struct {
	AnalyticsMap map[string]Analytics
}

func (a *BannerStorage) addAdvertisement(ad Banner) {
	a.BannerMap[ad.BannerID] = ad
}

func (a *BannerStorage) getAdvertisements() Banner {
	var ads []Banner
	for _, ad := range a.BannerMap {
		ads = append(ads, ad)
	}
	return ads[rand.Intn(len(ads))]
}

func (a *BannerStorage) deleteAdvertisement(id string) {
	delete(a.BannerMap, id)
}

func (b BannerStorage) sendBanner(id string) Banner {
	return b.BannerMap[id]
}

func (a AnalyticsStorage) getAnalytics(id string) Analytics {
	return a.AnalyticsMap[id]
}

func (a AnalyticsStorage) addClick(id string) {
	a.AnalyticsMap[id].Clicks[0]++
}
