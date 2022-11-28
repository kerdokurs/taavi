package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"github.com/robfig/cron/v3"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DatabaseService struct {
	db            *firestore.Client
	stopListening context.CancelFunc
}

func (ds *DatabaseService) Init() {
	ctx := context.Background()
	var app *firebase.App
	var err error
	if os.Getenv("ENV") != "prod" {
		opt := option.WithCredentialsFile("firebase_creds.json")
		app, err = firebase.NewApp(context.Background(), nil, opt)
	} else {
		app, err = firebase.NewApp(context.Background(), nil)
	}
	if err != nil {
		log.Fatalf("Could not create Firebase app: %v\n", err)
	}
	ds.db, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalf("Could not create Firestore client: %v\n", err)
	}

	ds.setup()
}

func (ds *DatabaseService) Stop() {
	ds.stopListening()
	ds.db.Close()
}

func (ds *DatabaseService) setup() {
	var ctx context.Context
	ctx, ds.stopListening = context.WithCancel(context.Background())

	// Cron Info
	go ds.setupInfos(ctx)

	fmt.Println("All database hooks are set up!")
}

func (ds *DatabaseService) setupInfos(ctx context.Context) {
	it := ds.db.Collection("jobs").Snapshots(ctx)
	for {
		snap, err := it.Next()
		if status.Code(err) == codes.DeadlineExceeded {
			log.Printf("[CronInfo Snapshot] cancelled: %v\n", err)
			return
		}

		if err != nil {
			log.Printf("[CronInfo Snapshot] error: %v\n", err)
			continue
		}

		fmt.Printf("[CronInfo Snapshot] Update\n")

		docs, err := snap.Documents.GetAll()
		if err != nil {
			log.Printf("Could not get all documents: %v\n", err)
			continue
		}

		infos := make([]CronInfo, len(docs))
		for i, doc := range docs {
			doc.DataTo(&infos[i])
			infos[i].Id = doc.Ref.ID
		}
		cronInfos.Set(infos)
	}
}

type CronType int

const (
	Message CronType = iota
	Update  CronType = iota
)

type CronInfo struct {
	Id   string   `json:"id" firestore:"id"`
	Type CronType `json:"type" firestore:"type"`

	CronTime string `json:"cron_time" firestore:"cron_time"`

	StreamId string   `json:"stream_id" firestore:"stream_id"`
	TopicId  string   `json:"topic_id" firestore:"topic_id"`
	Emails   []string `json:"emails" firestore:"emails"` // Used for private messages
	Url      string   `json:"url" firestore:"url"`       // Used for custom url jobs
	Content  string   `json:"content" firestore:"content"`

	Enabled bool `json:"enabled" firestore:"enabled"`

	CreatedAt time.Time `json:"created_at" firestore:"created_at,serverTimestamp"`

	EntryId cron.EntryID
}

func (ci *CronInfo) Changed(other *CronInfo) bool {
	return ci.Type != other.Type || ci.CronTime != other.CronTime || ci.StreamId != other.StreamId || ci.TopicId != other.TopicId || SliceEq(ci.Emails, other.Emails)
}
