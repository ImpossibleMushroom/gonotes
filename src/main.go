package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"./models"
	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
)

// In a production environment, pass through command line args or environment variables
var store = sessions.NewCookieStore([]byte("super_secret_key"))

type flashMsg struct {
	Class  string
	Prefix string
	Text   string
}

// TODO: Expire time, expire after view etc
// TODO: encryption?
// TODO: Simplify the way to access sessions? A lot of re-used code with error handling
// TODO:		Especially AddFlash, setting user errors looks messy and large for now

func randomString(n int) string {
	buff := make([]byte, n)
	rand.Read(buff)
	str := base64.StdEncoding.EncodeToString(buff)
	// lazy / wrong / dirty way to keep this string from having / in it
	str = strings.ReplaceAll(str, "/", "a")
	return str[:n]
}

// Handler to view existing notes
func viewHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	nid := p.ByName("noteId")

	session, err := store.Get(r, "gonotes")
	hasSession := true
	if err != nil {
		// This error is not essential to the app
		log.Println(err)
		hasSession = false
	}

	if len(nid) < 12 {
		if hasSession {
			session.AddFlash(flashMsg{"alert-danger", "Oops!", "You tried to goto an invalid note!"}, "messages")
			session.Save(r, w)
		}

		http.Redirect(w, r, "/", 302)
		return
	}

	note, err := models.GetNote(nid)
	if err != nil {
		log.Println(err)
		if hasSession {
			session.AddFlash(flashMsg{"alert-danger", "Oh no!", "We couldn't get that note for you!"}, "messages")
			session.Save(r, w)
		}
		http.Redirect(w, r, "/", 302)
		return
	}

	note.AddView(r)

	// Check if user is the creator & get flash messages
	var messages []interface{}
	if hasSession {
		if session.Values["ownershipToken"] == nil {
			session.Values["ownershipToken"] = randomString(16)
		}

		if session.Values["ownershipToken"].(string) == note.OwnerID {
			note.UserIsCreator = true
		}
		messages = session.Flashes("messages")
		session.Save(r, w)
	}

	serveTemplate(w, "./templates/view.html", "view.html", struct {
		Note     models.Note
		Messages []interface{}
	}{note, messages})
}

// POST route to make a new note
func makeNoteHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	nid := randomString(12)

	session, err := store.Get(r, "gonotes")
	hasSession := true
	if err != nil {
		// This error is not essential to the app
		log.Println(err)
		hasSession = false
	}

	ownershipToken := randomString(16)
	if hasSession {
		if session.Values["ownershipToken"] == nil {
			session.Values["ownershipToken"] = ownershipToken
			session.Save(r, w)
		} else {
			ownershipToken = session.Values["ownershipToken"].(string)
		}
	}

	postContent := r.FormValue("post_content")
	// Post must be between 1 and 4000 characters
	plen := len(postContent)
	if plen <= 0 || plen >= 4000 {
		if hasSession {
			session.AddFlash(flashMsg{"alert-warning", "Sorry!", "That post is too long! It must be <= 4000 characters. Yours is " + strconv.Itoa(plen)}, "messages")
			session.Save(r, w)
		}
		http.Redirect(w, r, "/", 303)
		return
	}

	postLang := r.FormValue("langMode")

	inserted, err := models.MakeNote(models.Note{
		ID:      nid,
		Content: postContent,
		OwnerID: ownershipToken,
		Views:   0,
		Lang:    postLang})
	if err != nil || inserted == false {
		log.Println(err)
		if hasSession {
			session.AddFlash(flashMsg{"alert-danger", "Sorry!", "We couldn't make your post."}, "messages")
			session.Save(r, w)
		}
		http.Redirect(w, r, "/", 302)
		return
	}

	if hasSession {
		err = session.Save(r, w)
		if err != nil {
			log.Println(err)
		}
	}

	http.Redirect(w, r, "/view/"+nid, 302)
}

// POST route to delete a note, checks (using session for now) if the user is the owner of it
func deleteNoteHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	nid := r.FormValue("noteId")

	session, err := store.Get(r, "gonotes")
	hasSession := true
	if err != nil {
		// This error is not essential to the app
		log.Println(err)
		hasSession = false
	}

	if len(nid) < 12 {
		if hasSession {
			session.AddFlash(flashMsg{"alert-danger", "Oops!", "You tried delete an invalid note!"}, "messages")
			session.Save(r, w)
		}
		http.Redirect(w, r, "/", 302)
		return
	}

	// Check for ownership
	if hasSession {
		if session.Values["ownershipToken"] != nil {
			deleted, err := models.DeleteNote(nid, session.Values["ownershipToken"].(string))
			if err != nil || deleted != true {
				log.Println(err)
				session.AddFlash(flashMsg{"alert-danger", "Sorry!", "We weren't able to delete your note!"}, "messages")
				session.Save(r, w)
				http.Redirect(w, r, "/view/"+nid, 302)
				return
			}

			session.AddFlash(flashMsg{"alert-success", "Yay!", "We deleted that note!"}, "messages")
			session.Save(r, w)
			http.Redirect(w, r, "/", 302)
			return
		}

		session.AddFlash(flashMsg{"alert-danger", "Wait!", "You don't own that note!"}, "messages")
		session.Save(r, w)
	}

	http.Redirect(w, r, "/view/"+nid, 302)
}

// 404 route
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	serveTemplate(w, "./templates/404.html", "404.html", struct {
		URL string
	}{r.URL.EscapedPath()})
}

// Default page to make new notes on
func homeHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var messages []interface{}
	session, err := store.Get(r, "gonotes")
	if err != nil {
		log.Println(err)
	} else {
		messages = session.Flashes("messages")
		if session.Values["ownershipToken"] == nil {
			session.Values["ownershipToken"] = randomString(16)
		}
		session.Save(r, w)
	}

	serveTemplate(w, "./templates/editor.html", "editor.html", struct {
		Messages []interface{}
	}{messages})
}

// Serve template to given http.ResponseWriter, help prevent repeat code/error checking
func serveTemplate(w http.ResponseWriter, templatePath string, templateName string, data interface{}) {
	page, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}

	err = page.ExecuteTemplate(w, templateName, data)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}

// Application entry point
func main() {
	err := models.ConnectToDatabase()
	if err != nil {
		log.Fatal(err)
	}

	gob.Register(&flashMsg{})

	router := httprouter.New()

	router.GET("/view/:noteId", viewHandler)
	router.GET("/", homeHandler)

	router.POST("/make_note/", makeNoteHandler)
	router.POST("/delete_note/", deleteNoteHandler)

	router.ServeFiles("/assets/*filepath", http.Dir("./public"))
	router.NotFound = http.HandlerFunc(notFoundHandler)

	fmt.Println("Running on port 80...")
	log.Fatal(http.ListenAndServe(":80", router))
}
