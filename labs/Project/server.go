package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
const (
	mongodbEndpoint = "mongodb://172.17.0.2:27017" // Find this from the Mongo container
)

type Post struct {
	ID           primitive.ObjectID `bson:"_id"`
	TaskId       string             `bson:"taskId"`
	Task         string             `bson:"body"`
	Tags         []string           `bson:"tags"`
	Asignee   	 string             `bson:"asignee"`
	TimeEstimate string			    `bson:timeEstimate`
	Priority     string				`bson:priority`
	CreatedAt    time.Time          `bson:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at"`
}

func db() *mongo.Client {
	clientOptions := options.Client().ApplyURI(mongodbEndpoint)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

var userDB = db().Database("tasks").Collection("posts")

func main() {
	http.HandleFunc("/list", list)
	http.HandleFunc("/resolve", resolve)
	http.HandleFunc("/create", create)
	http.HandleFunc("/update", update)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

//list
func list(w http.ResponseWriter, req *http.Request) {
	cursor, err := userDB.Find(context.TODO(), bson.M{})
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		var results bson.M
		if err = cursor.Decode(&results); err != nil {
			log.Fatal(err)
		}
		fmt.Fprint(w, results)
	}
}

//create?taskId=4&task=task4&asignee=bill&timeEstimate=3days&priority=2
func create(w http.ResponseWriter, req *http.Request) {
	taskId := req.URL.Query().Get("taskId")
	task := req.URL.Query().Get("task")
	asignee := req.URL.Query().Get("asignee")
	timeEstimate := req.URL.Query().Get("timeEstimate")
	priority := req.URL.Query().Get("priority")
	fmt.Fprintf(w, taskId)

	res, _ := userDB.InsertOne(context.TODO(), &Post{
		ID:           primitive.NewObjectID(),
		TaskId:       taskId,
		Tags:         []string{"CRITICAL"},
		Task:         task,
		Asignee:      asignee,
		TimeEstimate: timeEstimate,
		Priority:     priority,
		CreatedAt:    time.Now(),
	})
	fmt.Fprint(w, " was inserted\n", res.InsertedID.(primitive.ObjectID).Hex())
}

//update?taskId=3&task=newtask3&asignee=bill&timeEstimate=3days&priority=2
func update(w http.ResponseWriter, req *http.Request) {
	taskId := req.URL.Query().Get("taskId")
	task := req.URL.Query().Get("task")
	// asignee := req.URL.Query().Get("asignee")
	// timeEstimate := req.URL.Query().Get("timeEstimate")
	// priority := req.URL.Query().Get("priority")
	fmt.Fprintf(w, "updating"+taskId)

	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}
	filter := bson.M{"taskId": taskId} //unsure about changing these two string literals to id and task
	update := bson.M{"$set": bson.M{"task": task}} 
	res := userDB.FindOneAndUpdate(context.TODO(), filter, update, &opt)
	_ = res
}

//resolve?taskId=1
func resolve(w http.ResponseWriter, req *http.Request) {
	taskId:= req.URL.Query().Get("taskId")
	res, _ := userDB.DeleteOne(context.TODO(), bson.D{{"taskId", taskId}})
	_ = res
	fmt.Fprintf(w, "deleted docs: "+ taskId)
}