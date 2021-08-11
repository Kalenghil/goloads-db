package main

import (
	"bytes"
	"encoding/json"
	"net/http"
)

var GoLoAdsToken string

type MoneyRequest struct {
	Token       string `json:"token"`
	AccountID   int    `json:"account_id"`
	Amount      int    `json:"amount"`
	Description string `json:"description"`
}

func sendMoneyToUser(user_id int, money_am float64) {
	var moneyRequest = MoneyRequest{
		Token:       GoLoAdsToken,
		AccountID:   user_id,
		Amount:      int(money_am),
		Description: "Вывод средств со счёта GoloAds на счет пользователя",
	}

	postBody, _ := json.Marshal(moneyRequest)
	responseBody := bytes.NewBuffer(postBody)
	http.Post("https://bank.goto.msk.ru/api/send", "application/json", responseBody)
}
