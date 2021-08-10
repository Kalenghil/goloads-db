package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type DataBase interface {
	getUserFromId()
}

type adsServer struct {
	userStorage          UserStorage
	advertisementStorage BannerStorage
	analyticsStorage     AnalyticsStorage
}

type IDRequest struct {
	ID string `json:"id"`
}

type BannerRequest struct {
	URL     string   `json:"url"`
	Image   string   `json:"image"`
	Domains []string `json:"domains"`
}

type test struct {
	Body string `json:"body"`
}

var Test = test{Body: "OK"}

var counter int = 0

func checkForError(err error, errorCode int, w http.ResponseWriter) {
	if err != nil {
		if errorCode == 0 {
			fmt.Println(err)
			if errorCode != 0 {
				http.Error(w, http.StatusText(errorCode), errorCode)
			}
			return
		}
	}
}

func PreInnitiallizeStuff(w http.ResponseWriter, r *http.Request) {
	fmt.Println("got request with method", r.Method, counter)
	counter++
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "*")
}

func (a *adsServer) deleteBannerHandler(w http.ResponseWriter, r *http.Request) {
	PreInnitiallizeStuff(w, r)

	if r.Method != "DELETE" {
		http.Error(w,
			http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest)
		return
	}

	rawData, err := ioutil.ReadAll(r.Body)
	checkForError(err, http.StatusInternalServerError, w)
	fmt.Println(string(rawData))

	var id_request IDRequest
	err = json.Unmarshal(rawData, &id_request)
	checkForError(err, http.StatusBadRequest, w)

	a.advertisementStorage.deleteAdvertisement(id_request.ID)

	bytes, err := json.Marshal(Test)
	checkForError(err, http.StatusInternalServerError, w)
	_, err = w.Write(bytes)
	if err != nil {
		return
	}
	// fmt.Fprint(w, string(bytes))
}

func (a *adsServer) sendBannerHandler(w http.ResponseWriter, r *http.Request) {
	PreInnitiallizeStuff(w, r)

	ads := a.advertisementStorage.getAdvertisements()

	bytes, err := json.Marshal(ads)
	checkForError(err, http.StatusInternalServerError, w)

	fmt.Fprint(w, string(bytes))
}

func (a *adsServer) receivePostHandler(w http.ResponseWriter, r *http.Request) {
	PreInnitiallizeStuff(w, r)

	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	rawBody, err := ioutil.ReadAll(r.Body)
	fmt.Println(string(rawBody))
	checkForError(err, http.StatusBadRequest, w)
	var newAdvertisement Banner
	err = json.Unmarshal(rawBody, &newAdvertisement)
	if err != nil {
		fmt.Println(err)
		return
	}
	a.advertisementStorage.addAdvertisement(newAdvertisement)

	a.advertisementStorage.putAdvertisementIntoDB(newAdvertisement.BannerID)

	bytes, err := json.Marshal(Test)
	checkForError(err, http.StatusInternalServerError, w)

	fmt.Fprint(w, string(bytes))

}

func (a *adsServer) receiveClickHandler(w http.ResponseWriter, r *http.Request) {
	PreInnitiallizeStuff(w, r)
	if r.Method != "POST" {
		http.Error(w,
			http.StatusText(http.StatusBadRequest),
			http.StatusOK)
		return
	}

	rawBody, err := ioutil.ReadAll(r.Body)
	fmt.Println(string(rawBody))
	checkForError(err, http.StatusBadRequest, w)

	var addClicks IDRequest
	err = json.Unmarshal(rawBody, &addClicks)

	a.analyticsStorage.addClick(addClicks.ID)
	a.analyticsStorage.addClickToDB(addClicks.ID)
}

func (a *adsServer) sendAnalyticsHandler(w http.ResponseWriter, r *http.Request) {
	PreInnitiallizeStuff(w, r)

	id := r.URL.Query().Get("id")
	analytics := a.analyticsStorage.AnalyticsMap[id]

	bytes, err := json.Marshal(analytics)
	checkForError(err, http.StatusInternalServerError, w)

	fmt.Fprint(w, string(bytes))

}

func (a *adsServer) sendFaviconHandler(w http.ResponseWriter, r *http.Request) {
	PreInnitiallizeStuff(w, r)

	http.ServeFile(w, r, "favicon.ico")
}

func (a *adsServer) receiveBannerImageHandler(w http.ResponseWriter, r *http.Request) {
	PreInnitiallizeStuff(w, r)

	rawData, err := ioutil.ReadAll(r.Body)
	checkForError(err, http.StatusBadRequest, w)

	var newImage BannerRequest
	if err := json.Unmarshal(rawData, &newImage); err != nil {
		fmt.Println(err)
		return
	}

	var newAdvertisement Banner
	newAdvertisement.BannerID = RandomString(19)
	newAdvertisement.Image = newImage.Image
	newAdvertisement.Domains = newImage.Domains
	newAdvertisement.DomainURL = newImage.URL
	newAdvertisement.ImageBase64 = true

	a.advertisementStorage.addAdvertisement(newAdvertisement)
	a.advertisementStorage.putAdvertisementIntoDB(newAdvertisement.BannerID)

	w.WriteHeader(http.StatusOK)

}

func main() {

	// initializing test objects

	TestAdvertisement := Banner{
		BannerID:    "nbn9ewnd",
		Image:       "https://klike.net/uploads/posts/2019-05/1556708032_1.jpg",
		DomainURL:   "yandex.ru",
		Domains:     []string{"stackoverflow.com"},
		ImageBase64: false,
	}

	TestAdvertisementStorage := BannerStorage{map[string]Banner{TestAdvertisement.BannerID: TestAdvertisement}}

	arrayLength := 14
	TestAnalytics := Analytics{
		BannerID:     "nbn9ewnd",
		Clicks:       RandomArray(arrayLength),
		UniqueClicks: RandomArray(arrayLength),
		Views:        RandomArray(arrayLength),
		UniqueViews:  RandomArray(arrayLength),
	}
	TestAnalyticsStorage := AnalyticsStorage{map[string]Analytics{TestAnalytics.BannerID: TestAnalytics}}
	AdsServer := adsServer{UserStorage{}, TestAdvertisementStorage, TestAnalyticsStorage}

	// initializing PostgreSQL database

	InnitializeDB()

	// initializing http handlers

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(AdsServer.sendBannerHandler))
	mux.Handle("/delete", http.HandlerFunc(AdsServer.deleteBannerHandler))
	mux.Handle("/favicon.ico", http.HandlerFunc(AdsServer.sendFaviconHandler))
	mux.Handle("/add", http.HandlerFunc(AdsServer.receivePostHandler))
	mux.Handle("/analytics", http.HandlerFunc(AdsServer.sendAnalyticsHandler))
	mux.Handle("/clicked", http.HandlerFunc(AdsServer.receiveClickHandler))
	log.Fatal(http.ListenAndServeTLS("doats.ml:8080", "certificate.crt", "private.key", mux))
}
