package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type karyawan struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Nama   string
	Email  string
	NoTlp  string
	Alamat string
}

var tmpl = template.Must(template.ParseGlob("view/*"))

const dbname = "form_db"
const collname = "karyawan"
const connlink = "mongodb://localhost:27017"

func index(w http.ResponseWriter, r *http.Request) {
	clientop := options.Client().ApplyURI(connlink)
	client, err := mongo.Connect(context.TODO(), clientop)
	if err != nil {
		log.Fatal(err)
	}
	coll := client.Database(dbname).Collection(collname)
	querry, err := coll.Find(context.TODO(), bson.D{{}}, options.Find())
	if err != nil {
		log.Fatal(err)
	}
	var re []*karyawan
	for querry.Next(context.TODO()) {
		var kar karyawan
		err := querry.Decode(&kar)
		if err != nil {
			log.Fatal(err)
		}
		re = append(re, &kar)
	}
	if err := querry.Err(); err != nil {
		log.Fatal(err)
	}
	querry.Close(context.TODO())
	tmpl.ExecuteTemplate(w, "index", re)
	fmt.Println("accessing index page")
}

func insert(w http.ResponseWriter, r *http.Request) {
	clientop := options.Client().ApplyURI(connlink)
	client, err := mongo.Connect(context.TODO(), clientop)
	if err != nil {
		log.Fatal(err)
	}
	coll := client.Database(dbname).Collection(collname)
	if r.Method == "POST" {
		kar := karyawan{Nama: r.FormValue("nama"), Email: r.FormValue("email"), NoTlp: r.FormValue("notelp"), Alamat: r.FormValue("alamat")}
		_, err := coll.InsertOne(context.TODO(), kar)
		if err != nil {
			log.Fatal(err)
		}
	}
	http.Redirect(w, r, "/", 301)
	fmt.Println("inserting", r.FormValue("nama"), "succeed")
}

func edit(w http.ResponseWriter, r *http.Request) {
	clientop := options.Client().ApplyURI(connlink)
	client, err := mongo.Connect(context.TODO(), clientop)
	if err != nil {
		log.Fatal(err)
	}
	coll := client.Database(dbname).Collection(collname)
	var kar karyawan
	rid := r.URL.Query().Get("id")
	rid = rid[10:34]
	brid, _ := primitive.ObjectIDFromHex(rid)
	err = coll.FindOne(context.TODO(), bson.M{"_id": brid}).Decode(&kar)
	if err != nil {
		log.Fatal(err)
	}
	tmpl.ExecuteTemplate(w, "edit", kar)
	fmt.Println("accessing id:", rid, "edit page")
}

func update(w http.ResponseWriter, r *http.Request) {
	var rid1 string
	clientop := options.Client().ApplyURI(connlink)
	client, err := mongo.Connect(context.TODO(), clientop)
	if err != nil {
		log.Fatal(err)
	}
	coll := client.Database(dbname).Collection(collname)
	if r.Method == "POST" {
		rid := r.FormValue("id")
		rid = rid[10:34]
		rid1 = rid
		brid, _ := primitive.ObjectIDFromHex(rid)
		ukar := bson.M{
			"$set": bson.M{
				"nama":   r.FormValue("nama"),
				"email":  r.FormValue("email"),
				"notelp": r.FormValue("notelp"),
				"alamat": r.FormValue("alamat"),
			},
		}
		_, err := coll.UpdateOne(context.TODO(), bson.M{"_id": brid}, ukar)
		if err != nil {
			log.Fatal(err)
		}
	}
	http.Redirect(w, r, "/", 301)
	fmt.Println("updating id:", rid1, "succeed")
}

func del(w http.ResponseWriter, r *http.Request) {
	var rid string
	clientop := options.Client().ApplyURI(connlink)
	client, err := mongo.Connect(context.TODO(), clientop)
	if err != nil {
		log.Fatal(err)
	}
	coll := client.Database(dbname).Collection(collname)
	rid = r.URL.Query().Get("id")
	rid = rid[10:34]
	brid, _ := primitive.ObjectIDFromHex(rid)
	_, err = coll.DeleteOne(context.TODO(), bson.M{"_id": brid})
	if err != nil {
		log.Fatal(err)
	}
	http.Redirect(w, r, "/", 301)
	fmt.Println("deleting id:", rid, "succeed")
}

func Add(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "Add", nil)
	fmt.Println("accessing new page")
}

func main() {
	fmt.Println("creating new mongodb client")
	clientop := options.Client().ApplyURI(connlink)
	client, err := mongo.Connect(context.TODO(), clientop)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("mongodb online")
	fmt.Println("creating router")
	r := mux.NewRouter()
	r.HandleFunc("/", index).Methods("GET")
	r.HandleFunc("/insert", insert).Methods("POST")
	r.HandleFunc("/Add", Add)
	r.HandleFunc("/edit", edit).Methods("GET")
	r.HandleFunc("/update", update).Methods("POST")
	r.HandleFunc("/delete", del).Methods("GET")
	fmt.Println("ready")
	http.ListenAndServe(":9090", r)
}
