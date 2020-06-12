package models

import (
	"log"
	"net/http"
	"strings"
	"time"
)

// Note - object to process Notes from DB into
// 	In the future, maybe store creators IP someway to check if they are the creator for UserIsCreator
type Note struct {
	ID            string
	Content       string
	Lang          string
	OwnerID       string
	Views         int
	UserIsCreator bool
}

// GetNote - Get a note from the database with given ID
func GetNote(id string) (Note, error) {
	var n Note
	err := db.QueryRow("SELECT \"Views\", \"Content\", \"Id\", \"Lang\", \"OwnerId\" FROM \"Notes\" WHERE \"Id\" = $1", id).Scan(&n.Views, &n.Content, &n.ID, &n.Lang, &n.OwnerID)
	return n, err
}

// MakeNote - Put a given note (n) in the db
func MakeNote(n Note) (bool, error) {
	res, err := db.Exec("INSERT INTO \"Notes\" (\"Id\", \"Content\", \"Views\", \"Lang\", \"OwnerId\") VALUES($1, $2, $3, $4, $5)", n.ID, n.Content, n.Views, n.Lang, n.OwnerID)
	if err != nil {
		return false, err
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	if ra <= 0 {
		return false, err
	}

	return true, nil
}

// DeleteNote - Delete a note with given nid
func DeleteNote(nid string, ownershipToken string) (bool, error) {
	// Delete Note
	res, err := db.Exec("DELETE FROM \"Notes\" WHERE \"Id\" = $1 AND \"OwnerId\" = $2", nid, ownershipToken)
	if err != nil {
		return false, err
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	if ra <= 0 {
		return false, err
	}

	// Delete Views
	res, err = db.Exec("DELETE FROM \"Views\" WHERE \"PostId\" = $1", nid)
	if err != nil {
		return false, err
	}
	ra, err = res.RowsAffected()
	if err != nil {
		return false, err
	}

	// No errors encountered, return true
	return true, nil
}

// AddView - Adds a view to a note
func (n *Note) AddView(r *http.Request) {
	// Get user IP
	ip := r.Header.Get("X-FORWARDED-FOR")
	if ip == "" {
		ip = r.RemoteAddr
	}
	if strings.Contains(ip, ":") {
		ip = ip[0:strings.Index(ip, ":")]
	}

	// Check if view already exists
	var uip string
	err := db.QueryRow("SELECT \"Ip\" FROM \"Views\" WHERE \"Ip\" = $1 AND \"PostId\" = $2", ip, n.ID).Scan(&uip)
	if err != nil && err != NoRows {
		log.Println(err)
		return
	}

	// Insert view if not exists
	if err == NoRows || uip != ip {
		res, err := db.Exec("INSERT INTO \"Views\" (\"Ip\", \"PostId\") VALUES($1, $2)", ip, n.ID)
		if err != nil {
			log.Println(err)
			return
		}

		inserted, err := res.RowsAffected()
		if err != nil {
			log.Println(err)
			return
		}

		if inserted >= 1 {
			// Change view on Notes table
			res, err = db.Exec("UPDATE \"Notes\" SET \"Views\" = $1, \"LastView\" = $2 WHERE \"Id\" = $3", (n.Views + 1), time.Now().Unix(), n.ID)
			if err != nil {
				log.Println(err)
				return
			}

			n.Views++
		}
	}
}
