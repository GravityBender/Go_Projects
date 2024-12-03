package main

import (
	"context"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	"snippetbox/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
}

const uri = "mongodb://localhost:27017"
const dbName = "Go_Practise"

func main() {

	addr := flag.String("addr", ":4000", "HTTP network address")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	ctx := context.TODO()
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		panic(err)
	}

	defer client.Disconnect(ctx)

	var result bson.M
	if err := client.Database(dbName).RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		panic(err)
	}

	templateCache, templateCacheErr := newTemplateCache()
	if templateCacheErr != nil {
		errorLog.Fatal(templateCacheErr)
	}

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		snippets: &models.SnippetModel{
			Client:  client,
			Context: ctx,
		},
		templateCache: templateCache,
	}

	// val, _ := client.ListDatabases(ctx, bson.D{}, options.ListDatabases())
	// infoLog.Printf("Client: %v", val.Databases)
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		// Call the new app.routes() method to get the servemux containing our routes.
		Handler: app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)

	serverErr := srv.ListenAndServe()
	errorLog.Fatal(serverErr)
}
