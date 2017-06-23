package client

import (
	"fmt"
	"log"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"

	"github.com/juju/ratelimit"
	"github.com/suneil/gmail/auth"
	"github.com/suneil/gmail/store"
)

var idList = make([]string, 0)
var throttle *ratelimit.Bucket

func showLabels(srv *gmail.Service) {
	user := "me"

	r, err := srv.Users.Labels.List(user).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve labels. %v", err)
	}

	if len(r.Labels) > 0 {
		fmt.Print("Labels:\n")
		for _, l := range r.Labels {
			fmt.Printf("- %s\n", l.Name)
		}
	} else {
		fmt.Print("No labels found.")
	}
}

func fetchInbox(srv *gmail.Service) {
	userID := "me"
	jobs := make(chan string)
	done := make(chan bool)

	go func(srv *gmail.Service) {
		for {
			pageToken, more := <-jobs
			log.Println("Received jobs item", pageToken, more)
			if more {
				fmt.Println("received job", pageToken)
				r, err := srv.Users.Messages.List(userID).Q("in:inbox").PageToken(pageToken).MaxResults(100).Do()
				if err != nil {
					log.Fatalf("Unable to retrieve inbox. %v", err)
				}
				for _, m := range r.Messages {
					// fmt.Printf("id :: %s\n", m.Id)
					idList = append(idList, m.Id)
					// go getMessage(m.Id, srv)

				}
				if r.NextPageToken != "" {
					go func() { jobs <- r.NextPageToken }()
					log.Println("Sent next page token", r.NextPageToken)
				} else {
					log.Println("Closing jobs")
					close(jobs)
				}

			} else {
				log.Println("received all jobs")
				done <- true
				return
			}
		}
	}(srv)

	log.Println("Starting first job")
	jobs <- ""

	log.Println("Waiting for done channel")
	<-done
	log.Println("All done")
}

func getMessage(id string, srv *gmail.Service, output chan<- bool) {
	r, err := srv.Users.Messages.Get("me", id).Format("raw").Do()
	if err != nil {
		log.Fatalf("Unable to get message %s: %s", id, err)
	}

	if r != nil {
		store.Add(r)
	}

	output <- true
}

// GetClient returns a client
func GetMail() {
	throttle = ratelimit.NewBucketWithRate(140.0, 10)

	ctx := context.Background()

	b, err := auth.GetClientSecrets()
	if err != nil {
		return
	}

	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := auth.GetClient(ctx, config)

	srv, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve gmail Client %v", err)
	}

	fetchInbox(srv)
	output := make(chan bool)

	for _, value := range idList {
		throttle.Wait(1)
		go getMessage(value, srv, output)
	}

	for i := 0; i < len(idList); i++ {
		<-output
	}
}
