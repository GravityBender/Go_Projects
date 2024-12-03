package models

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dbName = "Go_Practise"
const collectionName = "snippets"

type Snippet struct {
	ID      int `bson:"_id"`
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	Client  *mongo.Client
	Context context.Context
}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {

	coll := m.Client.Database(dbName).Collection(collectionName)

	currentTime := time.Now()
	expireTime := currentTime.AddDate(0, 0, expires)

	var lastDoc Snippet
	sort := bson.D{{"_id", -1}}
	opts := options.FindOne().SetSort(sort)
	err := coll.FindOne(context.TODO(), bson.D{}, opts).Decode(&lastDoc)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("No documents found")
		} else {
			panic(err)
		}
	}

	doc := Snippet{ID: lastDoc.ID + 1, Title: title, Content: content, Created: time.Now(), Expires: expireTime}

	result, err := coll.InsertOne(context.TODO(), doc)
	if err != nil {
		return 0, err
	}

	fmt.Printf("Inserted document with _id: %v\n", result.InsertedID)
	return doc.ID, nil
}

// This will return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (*Snippet, error) {

	coll := m.Client.Database(dbName).Collection(collectionName)
	filter := bson.M{"_id": id}

	var result Snippet

	err := coll.FindOne(m.Context, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("No documents found")
		} else {
			fmt.Printf("Error while finding snippet with id: %v\n", id)
		}
		return nil, err
	}
	res, _ := bson.MarshalExtJSON(result, false, false)
	fmt.Println(string(res))
	return &result, nil
}

// This will return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	coll := m.Client.Database(dbName).Collection(collectionName)

	filter := bson.D{{"expires", bson.D{{"$gt", time.Now()}}}}
	opts := options.Find().SetLimit(10)

	snippets := []*Snippet{}

	cursor, err := coll.Find(context.TODO(), filter, opts)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("No documents found")
		} else {
			fmt.Printf("Error while finding latest snippets\n")
		}
		return nil, err
	}

	if err = cursor.All(context.TODO(), &snippets); err != nil {
		panic(err)
	}

	for _, result := range snippets {
		res, _ := bson.MarshalExtJSON(result, false, false)
		fmt.Println(string(res))
	}
	return snippets, nil
}
