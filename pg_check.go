package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"time"
)

func main() {
	db, err := sql.Open("postgres", "user=postgres dbname=olx_map_go sslmode=disable")
	checkErr(err)
	//select_from_ads()
	db.Close()
}

func insert_urls() {
	err = db.QueryRow("INSERT INTO userinfo(username,departname,created) VALUES($1,$2,$3) returning uid;", "astaxie", "研发部门", "2012-12-09").Scan(&lastInsertId)
	checkErr(err)
	fmt.Println("last inserted id =", lastInsertId)
}

func select_from_ads() {
	rows, err := db.Query("SELECT * FROM ads")
	checkErr(err)

	for rows.Next() {
		var id int
		var description *string
		var price *string
		var created_at *time.Time
		var updated_at *time.Time
		var href *string
		var raw *string
		var pos *string
		var address *string
		var lng *string
		var lat *string
		var ad_type *string
		var room *string
		var posted_at *time.Time
		var title *string
		var filter *string
		var actual_at *string

		err = rows.Scan(&id, &description, &price, &created_at, &updated_at, &href,
			&raw, &pos, &address, &lng, &lat, &ad_type, &room, &posted_at, &title, &filter, &actual_at)
		checkErr(err)
		fmt.Println(id)
		fmt.Println(*href)
		if room != nil {
			fmt.Println(*room)
		}
		if ad_type != nil {
			fmt.Println(*ad_type)
		}
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
