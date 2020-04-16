package api

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/desteves/fle/fle"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

// Foobar is
type Foobar struct {
	ID      string `json:"_id" bson:"_id"`
	Name    string `json:"name" bson:"name"`
	Message string `json:"message" bson:"message"` //this field is encrypted if added via /foo endpoint, unencrypted if added via /bar endpoint
}

//CreateEncryptedFoobarHandler inserts a document to tutorial.foobar. It uses Field Encryption on the field "message" to insert a new Foobar JSON document to the tutorial Database.
func CreateEncryptedFoobarHandler(w http.ResponseWriter, r *http.Request) {

	var client *mongo.Client
	var doc Foobar

	defer func() {
		if r != nil {
			r.Body.Close()
		}
	}()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("! body - ", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &doc)
	if err != nil {
		log.Error("! unmarshal - ", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	client, err = fle.CreateMongoClient(os.Getenv("MONGODB_ATLAS_URI"), writeconcern.New(writeconcern.WMajority()), true) // Only difference from CreateFoobarHandler
	if err != nil {
		log.Error("! client - ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	result, err := client.Database("tutorial").Collection("foobar").InsertOne(context.TODO(), doc)
	if err != nil {
		log.Error("! insert - ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// success response
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
	return
}

//ReadEncryptedFoobarHandler reads a document with given id from  tutorial.foobar
func ReadEncryptedFoobarHandler(w http.ResponseWriter, r *http.Request) {

	var client *mongo.Client

	defer func() {
		if r != nil {
			r.Body.Close()
		}
	}()

	params := mux.Vars(r)
	ID, ok := params["id"]
	if !ok {
		log.Error("! request missing id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// database
	client, err := fle.CreateMongoClient(os.Getenv("MONGODB_ATLAS_URI"), writeconcern.New(writeconcern.WMajority()), true) // Only difference from ReadFoobarHandler
	if err != nil {
		log.Error("! client - ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.TODO())
	filter := bson.D{primitive.E{Key: "_id", Value: ID}}
	singleResult := client.Database("tutorial").Collection("foobar").FindOne(context.TODO(), filter)
	if singleResult.Err() != nil {
		log.Error("! find - ", singleResult.Err())
		if singleResult.Err() == mongo.ErrNoDocuments {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var doc bson.M
	err = singleResult.Decode(&doc)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// success response
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	json.NewEncoder(w).Encode(doc)
	return
}

//ReadFoobarHandler reads a document with given id from  tutorial.foobar
func ReadFoobarHandler(w http.ResponseWriter, r *http.Request) {

	var client *mongo.Client

	defer func() {
		if r != nil {
			r.Body.Close()
		}
	}()

	params := mux.Vars(r)
	ID, ok := params["id"]
	if !ok {
		log.Error("! request missing id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// database
	client, err := fle.CreateMongoClient(os.Getenv("MONGODB_ATLAS_URI"), writeconcern.New(writeconcern.WMajority()), false) // Only difference from ReadFoobarHandler
	if err != nil {
		log.Error("! client - ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.TODO())
	filter := bson.D{primitive.E{Key: "_id", Value: ID}}
	singleResult := client.Database("tutorial").Collection("foobar").FindOne(context.TODO(), filter)
	if singleResult.Err() != nil {
		log.Error("! find - ", singleResult.Err())
		if singleResult.Err() == mongo.ErrNoDocuments {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var doc bson.M
	err = singleResult.Decode(&doc)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// success response
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	json.NewEncoder(w).Encode(doc)
	return
}

// CreateFoobarHandler is
func CreateFoobarHandler(w http.ResponseWriter, r *http.Request) {

	var client *mongo.Client
	var doc Foobar

	defer func() {
		if r != nil {
			r.Body.Close()
		}
	}()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("! body - ", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &doc)
	if err != nil {
		log.Error("! unmarshal - ", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	client, err = fle.CreateMongoClient(os.Getenv("MONGODB_ATLAS_URI"), writeconcern.New(writeconcern.WMajority()), false)
	if err != nil {
		log.Error("! client - ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	result, err := client.Database("tutorial").Collection("foobar").InsertOne(context.TODO(), doc)
	if err != nil {
		log.Error("! insert - ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// success response
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
	return
}
