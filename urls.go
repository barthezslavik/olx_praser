package main

import (
	"database/sql"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	_ "github.com/lib/pq"
	"os"
	"strings"
	"time"
)

var db *sql.DB
var storage []string

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func parse(page_url string, ch chan string, chFinished chan bool) {
	defer func() {
		chFinished <- true
	}()

	doc, err := goquery.NewDocument(page_url)
	checkErr(err)

	doc.Find("a.link.linkWithHash.detailsLink").Each(func(i int, s *goquery.Selection) {
		href, ok := s.Attr("href")
		if ok {
			left := strings.Split(href, "#")
			storage = append(storage, left[0])
			ch <- left[0]
		}
	})
}

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "postgres"
	DB_NAME     = "olx_map_go"
)

func main() {
	foundUrls := make(map[string]bool)
	chUrls := make(chan string)
	chFinished := make(chan bool)

	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	checkErr(err)

	defer db.Close()

	go parse("http://www.olx.ua/nedvizhimost/arenda-kvartir/kiev", chUrls, chFinished)

	for page := 2; page <= 400; page++ {
		url_name := fmt.Sprint("http://www.olx.ua/nedvizhimost/arenda-kvartir/kiev/?page=", page)
		go parse(url_name, chUrls, chFinished)
	}

	go parse("http://www.olx.ua/nedvizhimost/arenda-komnat/kiev/", chUrls, chFinished)

	for page := 2; page <= 100; page++ {
		url_name := fmt.Sprint("http://www.olx.ua/nedvizhimost/arenda-komnat/kiev/?page=", page)
		go parse(url_name, chUrls, chFinished)
	}

	for c := 0; c <= 500; {
		select {
		case url := <-chUrls:
			foundUrls[url] = true
		case <-chFinished:
			c++
		}
		if c == 500 {
			for _, href := range storage {
				t := time.Now()
				_, err = db.Exec("INSERT INTO ads(href, created_at, updated_at) VALUES($1, $2, $3)", href, t, t)
				checkErr(err)
			}
			os.Exit(0)
		}
	}

	close(chUrls)
}
