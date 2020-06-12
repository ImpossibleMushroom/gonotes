/**
	When ran, this will clear all expired and inactive notes from the postgres database
	This can be ran in a cron job or a bash script etc to repeat it when needed
**/

package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	// Postgres driver for database/sql
	_ "github.com/lib/pq"
)

// Database details
//	In a production environment, pass these through command line args or environment variables
const (
	host     = "localhost"
	port     = 5432
	user     = "gonotes"
	password = "root"
	dbname   = "GoNotes"
)

var db *sql.DB

func deletePostsViews(post string) {
	res, err := db.Exec("DELETE FROM \"Views\" WHERE \"PostId\" = $1", post)
	if err != nil {
		log.Println(err)
		return
	}
	ra, err := res.RowsAffected()
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Printf("Deleted %d views belonging to %s\n", ra, post)
}

func main() {
	sqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	database, err := sql.Open("postgres", sqlInfo)
	if err != nil {
		log.Fatalln(err)
		return
	}
	db = database

	expiredTime := time.Now().AddDate(0, 0, -15)
	fmt.Printf("Deleting posts with the last view <= 15 days ago (%d)\n", expiredTime.Unix())
	rows, err := db.Query("DELETE FROM \"Notes\" WHERE \"LastView\" <= $1 RETURNING \"Id\"", expiredTime.Unix())
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var post string
		if err := rows.Scan(&post); err != nil {
			log.Fatal(err)
		}

		deletePostsViews(post)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
