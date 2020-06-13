package paginator

import (
	"context"
	"errors"
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
	Url     string    `bson:"url"`
	Mdate   time.Time `bson:"mdate"` // metadata adding time
	Ddate   time.Time `bson:"ddate"` // details info about car time
	Fdate   time.Time `bson:"fdate"` // fetched damage report time
	Fetched bool      `bson:"fetched"`
	Dir     string    `bson:"dir"`     // dir or arch path with reports
	Checked bool      `bson:"checked"` // means data checked (vin present)
}

func (car *Car) SaveId(mclient *mongoClient) error {
	db := (*mclient).client.Database(DB)
	coll := db.Collection(EXLE_CARS)
	var findResult *Car
	err := coll.FindOne(context.TODO(), bson.M{"_id": (*car).Id}).Decode(findResult)
	if err == mongo.ErrNoDocuments {
		_, err = coll.InsertOne(context.TODO(), car)
		if err != nil {
			log.Error.Fatalln(err)
		}
		return nil
	}

	if (*findResult).Meta.Checked {
		return nil
	}
	if (*findResult).Meta.Fetched {
		// change flag for new fetched, may be new data present
		filter := bson.M{"_id": (*car).Id}
		update := bson.M{"$set": bson.M{"meta.fetched": false}}
		_, err := coll.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			log.Error.Fatalln(err)
		}
	}
	return errors.New("unknown error")
}

// cleaning url from query,
// seting id
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
