package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "file::memory:?mode=memory&cache=shared")
	if err != nil {
		fmt.Println(err)
	}

	const Schema = `
	create table if not exists datatable (
		id integer primary key autoincrement,
		uuid text not null,
		vidid text not null,
		rating integer not null
	)
	`

	db.Exec(Schema)

	http.HandleFunc("/valget", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")

		value := r.FormValue("value")
		videoid := r.FormValue("videoid")
		userid := r.FormValue("userid")

		_, err := db.Exec("delete from datatable where uuid = ? and vidid = ?", userid, videoid)
		if err != nil {
			fmt.Println(err)
		}
		_, err = db.Exec("insert into datatable (uuid, vidid, rating) values (?, ?, ?)", userid, videoid, value)
		if err != nil {
			fmt.Println(err)
		}
	})

	http.HandleFunc("/yourget", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
    
		videoid := r.FormValue("videoid")
		userid := r.FormValue("userid")
		rows, err := db.Query("select rating from datatable where uuid = ? and vidid = ?", userid, videoid)
		if err != nil {
			fmt.Println(err)
		}
		defer rows.Close()
		for rows.Next() {
			var rating float64
			err = rows.Scan(&rating)
			if err != nil {
				fmt.Println(err)
			}
			w.Write([]byte(fmt.Sprintf("%d", int(rating))))
		}
	})

	http.HandleFunc("/avgget", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
    
		videoid := r.FormValue("videoid")
		rows, err := db.Query("select avg(rating) from datatable where vidid = ?", videoid)
		if err != nil {
			fmt.Println(err)
		}
		defer rows.Close()
		for rows.Next() {
			var avg float64
			err = rows.Scan(&avg)
			if err != nil {
				fmt.Println(err)
			}
			roundavg := fmt.Sprintf("%.1f", avg)
			w.Write([]byte(roundavg))
		}
	})

	http.HandleFunc("/countget", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
    
		videoid := r.FormValue("videoid")
		rows, err := db.Query("select count(*) from datatable where vidid = ?", videoid)
		if err != nil {
			fmt.Println(err)
		}
		defer rows.Close()
		for rows.Next() {
			var count int
			err = rows.Scan(&count)
			if err != nil {
				fmt.Println(err)
			}
			w.Write([]byte(fmt.Sprintf("%d", count)))
		}
	})
	
	http.HandleFunc("/deletedata", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
    
		userid := r.FormValue("userid")
		value := r.FormValue("value")
		if value == "1" {
			db.Exec("delete from datatable where uuid = ?", userid)
		}
	})

	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", nil)
}

