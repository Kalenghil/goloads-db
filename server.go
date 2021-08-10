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
	userStorage      UserStorage
	bannerStorage    BannerStorage
	analyticsStorage AnalyticsStorage
}

type BannerIDRequest struct {
	ID string `json:"id"`
}

type TelegramIDRequest struct {
	TelegramID int `json:"id"`
}

type MoneyResponse struct {
	Money float64 `json:"money"`
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

func returnHTTPError(errorCode int, w http.ResponseWriter){
	http.Error(w, http.StatusText(errorCode), errorCode)
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

	var id_request BannerIDRequest
	err = json.Unmarshal(rawData, &id_request)
	checkForError(err, http.StatusBadRequest, w)

	a.bannerStorage.deleteAdvertisement(id_request.ID)

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

	ads := a.bannerStorage.getAdvertisements()

	bytes, err := json.Marshal(ads)
	checkForError(err, http.StatusInternalServerError, w)

	fmt.Fprint(w, string(bytes))
}

func (a *adsServer) receiveBannerHandler(w http.ResponseWriter, r *http.Request) {
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
	a.bannerStorage.addAdvertisement(newAdvertisement)

	a.bannerStorage.putAdvertisementIntoDB(newAdvertisement.BannerID)

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

	var addClicks BannerIDRequest
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

	a.bannerStorage.addAdvertisement(newAdvertisement)
	a.bannerStorage.putAdvertisementIntoDB(newAdvertisement.BannerID)

	w.WriteHeader(http.StatusOK)

}

func (a *adsServer) getUserMoneyHandler(w http.ResponseWriter, r *http.Request) {
	PreInnitiallizeStuff(w, r)
	if r.Method != "POST" {
	 	returnHTTPError(http.StatusBadRequest, w)
	 	return
	}

	var newRequest TelegramIDRequest
	rawBytes, err := ioutil.ReadAll(r.Body)
	checkForError(err, http.StatusBadRequest, w)

	err = json.Unmarshal(rawBytes, &newRequest)
	checkForError(err, http.StatusBadRequest, w)

	user := a.userStorage.getUserByID(newRequest.TelegramID)

	var moneyAm MoneyResponse
	moneyAm.Money = user.Money
	bytes, err := json.Marshal(moneyAm)
	w.Write(bytes)
}

func main() {

	// initializing test objects

	arrayLength := 14
	TestAnalytics := Analytics{
		BannerID:     "nbn9ewnd",
		Clicks:       RandomArray(arrayLength),
		UniqueClicks: RandomArray(arrayLength),
		Views:        RandomArray(arrayLength),
		UniqueViews:  RandomArray(arrayLength),
	}
	TestAnalyticsStorage := AnalyticsStorage{map[string]Analytics{TestAnalytics.BannerID: TestAnalytics}}
	AdsServer := adsServer{UserStorage{}, BannerStorage{}, TestAnalyticsStorage}
	var bannerStorage BannerStorage = BannerStorage{map[string]Banner{}}
	for _, banner := range AdsServer.bannerStorage.getAdvertisementsFromDB(){
		bannerStorage.BannerMap[banner.BannerID] = banner
	}

	// initializing PostgreSQL database

	InnitializeDB()

	// initializing http handlers
	fmt.Println(bannerStorage.BannerMap)
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(AdsServer.sendBannerHandler))
	mux.Handle("/delete", http.HandlerFunc(AdsServer.deleteBannerHandler))
	mux.Handle("/favicon.ico", http.HandlerFunc(AdsServer.sendFaviconHandler))
	mux.Handle("/add", http.HandlerFunc(AdsServer.receiveBannerHandler))
	mux.Handle("/add/image", http.HandlerFunc(AdsServer.receiveBannerImageHandler))
	mux.Handle("/analytics", http.HandlerFunc(AdsServer.sendAnalyticsHandler))
	mux.Handle("/clicked", http.HandlerFunc(AdsServer.receiveClickHandler))
	mux.Handle("/info/get", http.HandlerFunc(AdsServer.getUserMoneyHandler))
	log.Fatal(http.ListenAndServeTLS("doats.ml:8080", "certificate.crt", "private.key", mux))
}
