package paginator

import (
	"context"
	"errors"
	"net/url"
	"strconv"
	"strings"
	"time"

	log "github.com/victron/simpleLogger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Car struct {
	Meta Meta `bson:"meta"`
	Data Data `bson:"data"`
}

type Meta struct {
	Id      int       `bson:"_id"`
	Url     string    `bson:"url"`
	Mdate   time.Time `bson:"mdate"` // metadata adding time
	Fdate   time.Time `bson:"fdate"` // fetched info about car time
	Fetched bool      `bson:"fetched"`
	Checked bool      `bson:"checked"` // means data checked (vin present)
}

type Data struct {
	Vin string `bson:"vin"`
}

func (car *Car) SaveId(mclient *mongoClient) error {
	db := (*mclient).client.Database(DB)
	coll := db.Collection(EXLE_CARS)
	var findResult *Car
	err := coll.FindOne(context.TODO(), bson.M{"_id": (*car).Meta.Id}).Decode(findResult)
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
		filter := bson.M{"_id": (*car).Meta.Id}
		update := bson.M{"$set": bson.M{"meta.fetched": false}}
		_, err := coll.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			log.Error.Fatalln(err)
		}
	}
	return errors.New("unknown error")
}

// parse Id from url
func ParseId(carUrl string) (int, error) {
	u, err := url.Parse(carUrl)
	if err != nil {
		log.Error.Fatal(err)
		return 0, err
	}
	path := strings.Split(u.Path, "/")
	if id, err := strconv.Atoi(path[len(path)-1]); err != nil {
		log.Error.Fatal(err)
		return 0, err
	} else {
		return id, nil
	}
}
