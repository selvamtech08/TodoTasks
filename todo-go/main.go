package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/selvamtech08/todogo/controller"
	"github.com/selvamtech08/todogo/store"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbName         = "todo-go"
	collectionName = "tasks"
	addr           = "localhost:8030"
)

var (
	taskHandler    controller.TaskHandler
	taskStore      store.TaskStoreager
	taskCollection *mongo.Collection
	client         *mongo.Client
	err            error
	ctx            context.Context
)

func init() {
	ctx = context.TODO()
	mongoOption := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err = mongo.Connect(ctx, mongoOption)
	if err != nil {
		log.Fatalln(err.Error())
	}
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalln(err.Error())
	}
	taskCollection = client.Database(dbName).Collection(collectionName)
	taskStore = store.NewTaskStore(taskCollection, ctx)
	taskHandler = controller.NewTaskController(taskStore)
}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("GET /task", taskHandler.GetAll)
	mux.HandleFunc("GET /task/pending", taskHandler.GetPending)
	mux.HandleFunc("GET /task/{name}", taskHandler.Get)
	mux.HandleFunc("POST /task", taskHandler.Create)
	mux.HandleFunc("PUT /task", taskHandler.Update)
	mux.HandleFunc("DELETE /task/{name}", taskHandler.Remove)

	server := &http.Server{
		Addr:         addr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		Handler:      mux,
	}

	defer client.Disconnect(ctx)

	log.Println("server running on", addr)
	log.Fatalln(server.ListenAndServe())
}
