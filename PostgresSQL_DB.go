package main

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"

	_ "github.com/lib/pq"
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

func InnitializeDB() *sql.DB {
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

	_, err = db.Query(`CREATE TABLE IF NOT EXISTS "Banners"
							(
    							"BannerID" text not null,
								"DomainURL" text not null,
								"Image"text,
								"Domains" text[]
							);`)
	if err != nil {
		return nil
	}

	fmt.Println("Successfully connected!")
	return db
}

func (b *BannerStorage) putAdvertisementIntoDB(id string) {
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

func (b *BannerStorage) getAdvertisementsFromDB() []Banner {
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
		fmt.Printf("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAa %v\n", banners)
		i++
		if err != nil {
			fmt.Println(err)
			return []Banner{}
		}
	}

	return banners
}

func (a *AnalyticsStorage) addClickToDB(id string) {
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Query(`UPDATE "Analytics" SET "Clicks"="clicks" + 1 WHERE BannerID=$1`, id)
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
		&user.Token,
		); err != nil {
		fmt.Println(err)
		return User{}
	}

	return user
}

// func (a *BannerStorage) getAdvertisementFromDB (id string)
