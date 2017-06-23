package store

import (
	b64 "encoding/base64"
	"fmt"
	"log"
	"net/mail"
	"strings"

	"google.golang.org/api/gmail/v1"
	"gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"
)

var (
	db *mgo.Session
)

func init() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	db = session
}

// Add interface for adding messages
func Add(raw *gmail.Message) {
	// dump(raw)
	msg := Message{ID: raw.Id}
	msg.HistoryID = raw.HistoryId
	msg.InternalDate = raw.InternalDate
	msg.ThreadID = raw.ThreadId
	// msg.Raw = raw.Raw
	msg.Snippet = raw.Snippet
	msg.LabelIDs = raw.LabelIds

	email := parseEmail(raw.Raw)

	header := email.Header

	msg.From = header.Get("From")
	msg.To = header.Get("To")
	msg.Date = header.Get("Date")
	msg.Subject = header.Get("Subject")

	c := db.DB("gmail").C("message")
	_ = c.Insert(msg)
}

func parseEmail(raw string) *mail.Message {
	sDec, _ := b64.URLEncoding.DecodeString(raw)
	rawMsg := string(sDec)

	reader := strings.NewReader(rawMsg)
	email, err := mail.ReadMessage(reader)
	if err != nil {
		log.Fatal(err)
	}

	return email
}

func dump(msg *gmail.Message) {
	email := parseEmail(msg.Raw)
	header := email.Header
	fmt.Println("Id:", msg.Id)
	fmt.Println("HistoryId:", msg.HistoryId)
	fmt.Println("ThreadId:", msg.ThreadId)
	fmt.Println("InternalDate:", msg.InternalDate)
	fmt.Println("Date:", header.Get("Date"))
	fmt.Println("From:", header.Get("From"))
	fmt.Println("To:", header.Get("To"))
	fmt.Println("Subject:", header.Get("Subject"))

	fmt.Println("Snippet:", msg.Snippet)
}
