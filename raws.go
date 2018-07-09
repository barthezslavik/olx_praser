package main

import (
	"database/sql"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	_ "github.com/lib/pq"
	"sync"
)

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "postgres"
	DB_NAME     = "olx_map_go"
)

var db *sql.DB
var storage = make(map[int]string)
var wg sync.WaitGroup

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func parse(page_url string, id int, ch chan string, chFinished chan bool) {
	defer func() {
		chFinished <- true
	}()

	doc, err := goquery.NewDocument(page_url)
	checkErr(err)

	if err != nil {
		wg.Done()
	}

	doc.Find("#offerdescription").Each(func(i int, s *goquery.Selection) {
		raw, _ := s.Html()

		if raw != "" {
			storage[id] = raw
		}
	})
	wg.Done()
}

func main() {
	chUrls := make(chan string)
	chFinished := make(chan bool)

	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	checkErr(err)

	rows, err := db.Query("SELECT id, href FROM ads WHERE raw IS NULL LIMIT 100")
	checkErr(err)

	for rows.Next() {
		var id int
		var href *string

		err = rows.Scan(&id, &href)
		checkErr(err)

		wg.Add(1)
		go parse(*href, id, chUrls, chFinished)
	}

	wg.Wait()

	//fmt.Println(storage)
	for k := range storage {
		_, err = db.Exec("UPDATE ads SET raw=$1 WHERE id=$2", storage[k], k)
		checkErr(err)
	}

	defer db.Close()
}
