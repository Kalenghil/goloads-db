package main

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"time"
)

/*const (
	host     = "localhost"
	port     = 5433
	user     = "postgres"
	password = "postgres"
	dbname   = "postgres"
)*/

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "goloads"
)

var psqlInfo = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
	host, port, user, password, dbname)

func InitializeDB() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	return db
}

func (b *BannerStorage) putBannerIntoDB(id string) {
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	var banner Banner
	banner = b.BannerMap[id]
	_, err = db.Query(`INSERT INTO "Banners" 
					VALUES ($1, $2, $3, $4, $5);`,
		banner.BannerID,
		banner.DomainURL,
		banner.Image,
		pq.Array(banner.Domains),
		banner.ImageBase64,
	)

	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = db.Query(`INSERT INTO "Analytics" 
					VALUES ($1, $2, $3, $4, $5);`,
		banner.BannerID,
		pq.Array([]int{}),
		pq.Array([]int{}),
		pq.Array([]int{}),
		pq.Array([]int{}),
		8080,
	)

	if err != nil {
		fmt.Println(err)
		return
	}
}

func (b *BannerStorage) getBannersFromDB() []Banner {
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	var banners []Banner
	rows, err := db.Query(`SELECT * FROM "Banners"`)
	if err != nil {
		fmt.Println(err)
		return []Banner{}
	}
	i := 0
	for rows.Next() {
		err = rows.Scan(&banners[i].BannerID, &banners[i].Image, &banners[i].DomainURL, &banners[i].Domains)
		if err != nil {
			fmt.Println(err)
			return []Banner{}
		}
		i++
	}

	return banners
}

func (a *AnalyticsStorage) addClickToDB(banner_id string, user_id int) {
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Query(`INSERT INTO "Clicks" VALUES ($1, $2, $3)`, banner_id, user_id, time.Now().Unix())
	if err != nil {
		fmt.Println(err)
		return
	}

}

func (a *AnalyticsStorage) addViewToDB(banner_id string, user_id int) {
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Query(`INSERT INTO "Views" VALUES ($1, $2, $3)`, banner_id, user_id, time.Now().Unix())
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (a *AnalyticsStorage) addUserToDB(user User) {
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Query(`INSERT INTO "User" ("Firstname", "Lastname", ID, "Account") VALUES ($1, $2, $3, $4)`,
		user.Firstname,
		user.Lastname,
		user.ID,
		user.Account,
	)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (a *UserStorage) getUserByID(telegramID int) User {
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var user User
	row := db.QueryRow(`SELECT * FROM "Users" WHERE ID=$1`, telegramID)
	if err := row.Scan(
		&user.Firstname,
		&user.Lastname,
		&user.ID,
		&user.Account,
		&user.Money,
		&user.PhotoURL,
		&user.Username,
		&user.Hash,
	); err != nil {
		fmt.Println(err)
		return User{}
	}

	return user
}

func (u *UserStorage) resetUserMoney(telegramID int) {
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Query(`UPDATE "Users"
			SET "Money"=0.0
			WHERE ID=$1`, telegramID)

}

func (u *UserStorage) addMoney(telegramID int, moneyAmount float64) {
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Query(`UPDATE "Users"
			SET "Money"="Money"+$1
			WHERE ID=$2`,moneyAmount, telegramID)

}
// func (a *BannerStorage) getAdvertisementFromDB (id string)
