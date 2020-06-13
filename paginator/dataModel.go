package paginator

import (
	"context"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/victron/exleacar/paginator/details"
	log "github.com/victron/simpleLogger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Car struct {
	Id   int          `bson:"_id"`
	Meta Meta         `bson:"meta"`
	Data details.Data `bson:"data"`
}

type Meta struct {
	Url   string    `bson:"url"`
	Mdate time.Time `bson:"mdate"` // metadata adding time
	// TODO: add after details
	Ddate   time.Time `bson:"ddate"`   // details info about car time
	Fdate   time.Time `bson:"fdate"`   // fetched damage report time
	Fetched bool      `bson:"fetched"` // means report fetched
	Dir     string    `bson:"dir"`     // dir or arch path with reports
	Checked bool      `bson:"checked"` // means vin present and data fetched
}

// inserting doc. no any checking if doc present
func (car *Car) InsertFullDoc(mclient *mongoClient) error {
	db := (*mclient).client.Database(DB)
	coll := db.Collection(EXLE_CARS)
	_, err := coll.InsertOne(context.TODO(), car)
	if err != nil {
		log.Error.Fatalln(err)
	}
	return nil
}

// check if data for Id present. To avoid double collecting
// if id not found adding doc to DB, at this moment only meta can present in doc, and return false
// if flag `checked==true` all data exist - return true
// if data fetched but not checked - set flag fetched to false for new fetching
// insert: if true insert doc; needed for possible work in separate service. No need to insert if next spep is to collect details.
func (car *Car) IdPresent(mclient *mongoClient, insert bool) bool {
	db := (*mclient).client.Database(DB)
	coll := db.Collection(EXLE_CARS)
	findResult := Car{}
	err := coll.FindOne(context.TODO(), bson.M{"_id": (*car).Id}).Decode(&findResult)
	if err == mongo.ErrNoDocuments {
		log.Debug.Println("no doc with id=", (*car).Id, "in DB")
		if insert == true {
			_, err = coll.InsertOne(context.TODO(), car)
			if err != nil {
				log.Error.Fatalln(err)
			}
		}
		return false
	}

	if findResult.Meta.Checked == true {
		return true
	}
	if findResult.Meta.Fetched {
		// change flag for new fetched, may be new data present
		filter := bson.M{"_id": (*car).Id}
		update := bson.M{"$set": bson.M{"meta.fetched": false}}
		_, err := coll.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			log.Error.Fatalln(err)
		}
		return false
	}
	return false
}

// cleaning url from query,
// seting id
// TODO: create test
func (car *Car) ParseUrl() error {
	u, err := url.Parse((*car).Meta.Url)
	if err != nil {
		log.Error.Fatal(err)
		return err
	}
	u.RawQuery = "" // remove all params
	(*car).Meta.Url = u.String()

	path := strings.Split((*car).Meta.Url, "/")
	if (*car).Id, err = strconv.Atoi(path[len(path)-1]); err != nil {
		log.Error.Fatal(err)
		return err
	}
	return nil
}
