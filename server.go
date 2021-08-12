package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

const (
	MoneyForView  = 0.1
	MoneyForClick = 0.3
)

type AdsServer struct {
	userStorage      UserStorage
	bannerStorage    BannerStorage
	analyticsStorage AnalyticsStorage
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

func returnHTTPError(errorCode int, w http.ResponseWriter) {
	http.Error(w, http.StatusText(errorCode), errorCode)
}

func PreInnitiallizeStuff(w http.ResponseWriter, r *http.Request) {
	fmt.Println("got request with method", r.Method, counter, "URL:", r.URL)
	counter++
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}

func (a *AdsServer) sendExtensionIDHandler(w http.ResponseWriter, r *http.Request) {
	PreInnitiallizeStuff(w, r)

	if r.Method != "POST" {
		returnHTTPError(http.StatusBadRequest, w)
	}

	rawData, err := ioutil.ReadAll(r.Body)
	checkForError(err, http.StatusBadRequest, w)

	var id_request ExtensionIDRequest
	err = json.Unmarshal(rawData, &id_request)
	checkForError(err, http.StatusBadRequest, w)

	var id_response TelegramIDRequest
	id_response.TelegramID = a.userStorage.returnUserIDFromExtensionID(id_request.ExtensionID)

	bytes, err := json.Marshal(id_response)
	w.Write(bytes)
}

func (a *AdsServer) deleteBannerHandler(w http.ResponseWriter, r *http.Request) {
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

func (a *AdsServer) sendBannerHandler(w http.ResponseWriter, r *http.Request) {
	PreInnitiallizeStuff(w, r)

	ads := a.bannerStorage.getRandomBanner()

	bytes, err := json.Marshal(ads)
	checkForError(err, http.StatusInternalServerError, w)

	fmt.Fprint(w, string(bytes))
}

/*func (a *AdsServer) receivePostHandler(w http.ResponseWriter, r *http.Request) {
	PreInnitiallizeStuff(w, r)

	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	rawBody, err := ioutil.ReadAll(r.Body)
	fmt.Println(string(rawBody))
	checkForError(err, http.StatusBadRequest, w)
	var newBanner Banner
	err = json.Unmarshal(rawBody, &newBanner)
	if err != nil {
		fmt.Println(err)
		return
	}
	a.bannerStorage.addBanner(newBanner)

	a.bannerStorage.putBannerIntoDB(newBanner.BannerID)

	bytes, err := json.Marshal(Test)
	checkForError(err, http.StatusInternalServerError, w)

	fmt.Fprint(w, string(bytes))

}*/

func (a *AdsServer) bannerClickedHandler(w http.ResponseWriter, r *http.Request) {
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

	var addView BannerGotInteractedRequest
	err = json.Unmarshal(rawBody, &addView)

	a.analyticsStorage.addClickToDB(addView.BannerID, addView.TelegramID)
	a.userStorage.addMoney(addView.TelegramID, MoneyForClick)
}

func (a *AdsServer) sendAnalyticsHandler(w http.ResponseWriter, r *http.Request) {
	PreInnitiallizeStuff(w, r)

	id := r.URL.Query().Get("id")
	analytics := a.analyticsStorage.AnalyticsMap[id]

	bytes, err := json.Marshal(analytics)
	checkForError(err, http.StatusInternalServerError, w)

	fmt.Fprint(w, string(bytes))

}

func (a *AdsServer) sendFaviconHandler(w http.ResponseWriter, r *http.Request) {
	PreInnitiallizeStuff(w, r)

	http.ServeFile(w, r, "favicon.ico")
}

var newBanner Banner

func (a *AdsServer) receiveBannerFromAdmin1(w http.ResponseWriter, r *http.Request) {
	PreInnitiallizeStuff(w, r)

	rawData, err := ioutil.ReadAll(r.Body)
	checkForError(err, http.StatusBadRequest, w)

	var newBannerRequest BannerRequest
	if err := json.Unmarshal(rawData, &newBannerRequest); err != nil {
		fmt.Println(err)
		return
	}

	newBanner.BannerID = RandomString(20)
	newBanner.Domains = newBannerRequest.Domains
	newBanner.DomainURL = newBannerRequest.URL
	newBanner.Image = ""
	newBanner.ImageBase64 = true

	var IDResponse BannerIDRequest
	IDResponse.ID = newBanner.BannerID
	bytes, err := json.Marshal(IDResponse)
	checkForError(err, http.StatusInternalServerError, w)

	a.bannerStorage.addBanner(newBanner)
	w.Write(bytes)
}

func (a *AdsServer) receiveBannerFromAdmin2(w http.ResponseWriter, r *http.Request) {
	PreInnitiallizeStuff(w, r)

	rawData, err := ioutil.ReadAll(r.Body)
	checkForError(err, http.StatusBadRequest, w)

	bannerID := r.URL.Query().Get("id")

	a.bannerStorage.changeBannerImage(bannerID, string(rawData))
	a.bannerStorage.putBannerIntoDB(bannerID)

}

/*func (a *AdsServer) receiveBannerImageHandler(w http.ResponseWriter, r *http.Request) {
	PreInnitiallizeStuff(w, r)

	rawData, err := ioutil.ReadAll(r.Body)
	checkForError(err, http.StatusBadRequest, w)

	var newImage BannerRequest
	if err := json.Unmarshal(rawData, &newImage); err != nil {
		fmt.Println(err)
		return
	}

	var newBanner Banner
	newBanner.BannerID = RandomString(20)
	newBanner.Domains = newImage.Domains
	newBanner.DomainURL = newImage.URL
	newBanner.ImageBase64 = true

	a.bannerStorage.addBanner(newBanner)
	a.bannerStorage.putBannerIntoDB(newBanner.BannerID)

	w.WriteHeader(http.StatusOK)

}*/

func (a *AdsServer) getUserMoneyHandler(w http.ResponseWriter, r *http.Request) {
	PreInnitiallizeStuff(w, r)
	if r.Method != "POST" {
		returnHTTPError(http.StatusBadRequest, w)
		return
	}

	var newRequest TelegramIDRequest
	rawBytes, err := ioutil.ReadAll(r.Body)
	fmt.Println(string(rawBytes))
	checkForError(err, http.StatusBadRequest, w)

	err = json.Unmarshal(rawBytes, &newRequest)
	checkForError(err, http.StatusBadRequest, w)
	fmt.Println(newRequest)

	user := a.userStorage.getUserByID(newRequest.TelegramID)

	var money MoneyResponse
	money.Money = GtToMoney(user.Gotubles, user.Gopeykis)
	money.Username = user.Username
	money.PhotoURL = user.PhotoURL
	bytes, err := json.Marshal(money)
	w.Write(bytes)
}

func (a *AdsServer) bannerWatchedHandler(w http.ResponseWriter, r *http.Request) {
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

	var addClicks BannerGotInteractedRequest
	err = json.Unmarshal(rawBody, &addClicks)

	a.analyticsStorage.addViewToDB(addClicks.BannerID, addClicks.TelegramID)
	a.userStorage.addMoney(addClicks.TelegramID, MoneyForView)
}

func (a *AdsServer) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	PreInnitiallizeStuff(w, r)

	if r.Method != "POST" {
		returnHTTPError(http.StatusBadRequest, w)
		return
	}

	rawBytes, err := ioutil.ReadAll(r.Body)
	checkForError(err, http.StatusBadRequest, w)

	var NewUserRequest NewUserRequest
	err = json.Unmarshal(rawBytes, &NewUserRequest)
	checkForError(err, http.StatusBadRequest, w)

	var newUser User
	newUser.Firstname = NewUserRequest.FirstName
	newUser.Lastname = NewUserRequest.LastName
	newUser.ID = NewUserRequest.ID
	newUser.Account = NewUserRequest.ID
	newUser.Money = 0.0

	a.userStorage.addUserToDB(newUser)

}

func (a *AdsServer) sendMoneyToUserHandler(w http.ResponseWriter, r *http.Request) {
	PreInnitiallizeStuff(w, r)

	if r.Method != "POST" {
		returnHTTPError(http.StatusBadRequest, w)
		return
	}

	rawBytes, err := ioutil.ReadAll(r.Body)
	checkForError(err, http.StatusBadRequest, w)

	var userToSendMoney TelegramIDRequest
	err = json.Unmarshal(rawBytes, &userToSendMoney)
	checkForError(err, http.StatusBadRequest, w)

	var moneyAm = GtToMoney(
		a.userStorage.getUserByID(userToSendMoney.TelegramID).Gotubles,
		a.userStorage.getUserByID(userToSendMoney.TelegramID).Gopeykis)
	var statusOK = false
	response, err := sendMoneyToUser(userToSendMoney.TelegramID, moneyAm)

	if err != nil || response.StatusCode != http.StatusOK{
		returnHTTPError(http.StatusInternalServerError, w)
		return
	} else {
		a.userStorage.resetUserMoney(userToSendMoney.TelegramID)
		statusOK = true
	}

	Test.Body = "OK"

	a.analyticsStorage.addTransactionToDB(userToSendMoney.TelegramID, moneyAm, statusOK)
	bytes, err := ioutil.ReadAll(response.Body)
	checkForError(err, http.StatusInternalServerError, w)
	fmt.Println(string(bytes))
	w.Write(bytes)
}

func (a *AdsServer) linkExtensionIDToUserHandler(w http.ResponseWriter, r *http.Request) {
	PreInnitiallizeStuff(w, r)

	if r.Method != "POST" {
		http.Error(w,
			http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest)
		return
	}
	
	rawBytes, err := ioutil.ReadAll(r.Body)
	checkForError(err, http.StatusBadRequest, w)

	var linkRequest LinkExtensionIDRequest
	err = json.Unmarshal(rawBytes, &linkRequest)
	checkForError(err, http.StatusBadRequest, w)

	a.userStorage.linkExtensionID(linkRequest.ExtensionIDRequest, linkRequest.UserID)
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
	GoloAdsServer := AdsServer{UserStorage{}, TestAdvertisementStorage, TestAnalyticsStorage}

	// initializing PostgreSQL database

	InitializeDB()

	// initializing http handlers

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(GoloAdsServer.sendBannerHandler))
	mux.Handle("/delete", http.HandlerFunc(GoloAdsServer.deleteBannerHandler))
	mux.Handle("/add/image", http.HandlerFunc(GoloAdsServer.receiveBannerFromAdmin2))
	mux.Handle("/favicon.ico", http.HandlerFunc(GoloAdsServer.sendFaviconHandler))
	mux.Handle("/add", http.HandlerFunc(GoloAdsServer.receiveBannerFromAdmin1))
	mux.Handle("/analytics", http.HandlerFunc(GoloAdsServer.sendAnalyticsHandler))
	mux.Handle("/clicked", http.HandlerFunc(GoloAdsServer.bannerClickedHandler))
	mux.Handle("/watched", http.HandlerFunc(GoloAdsServer.bannerWatchedHandler))
	mux.Handle("/info/get", http.HandlerFunc(GoloAdsServer.getUserMoneyHandler))
	mux.Handle("/info/withdraw", http.HandlerFunc(GoloAdsServer.sendMoneyToUserHandler))
	mux.Handle("/user", http.HandlerFunc(GoloAdsServer.sendExtensionIDHandler))
	mux.Handle("/register", http.HandlerFunc(GoloAdsServer.registerUserHandler))
	mux.Handle("/user/link", http.HandlerFunc(GoloAdsServer.linkExtensionIDToUserHandler))

	log.Fatal(http.ListenAndServeTLS("doats.ml:8080", "certificate.crt", "private.key", mux))
}
