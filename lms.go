package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type student struct {
	Id   int64  `json:"id" bson:"id"`
	Name string `json:"name"`
}
type leave_data struct {
	Id     int64  `json:"id" bson:"id"`
	Reason string `json:"reason"`
	Frm    string `json:"frm"`
	To     string `json:"to"`
	Status string `json:"status"`
}

var s_col *mongo.Collection
var l_col *mongo.Collection

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func goDotEnvVariable(key string) string {
	err := godotenv.Load(".env")
	handleError(err)
	return os.Getenv(key)
}

func main() {
	goDotEnvVariable(".env")
	mongo_uri := goDotEnvVariable("mongourl")
	client, err := mongo.NewClient(options.Client().ApplyURI(mongo_uri))
	handleError(err)
	fmt.Println("connected")
	err = client.Connect(context.TODO())
	handleError(err)
	s_col = client.Database("student").Collection("data")
	l_col = client.Database("leave").Collection("data")
	r := mux.NewRouter()
	r.HandleFunc("/addstd", AddStd).Methods("POST")
	r.HandleFunc("/updstd", updateStd).Methods("POST")
	r.HandleFunc("/reqlve", reqLeave).Methods("POST")
	r.HandleFunc("/apprv", Approve).Methods("POST")
	r.HandleFunc("/all", allLeaves).Methods("GET")
	log.Fatal(http.ListenAndServe("0.0.0.0:9000", r))
}

func AddStd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var stud student
	json.NewDecoder(r.Body).Decode(&stud)
	id := stud.Id
	filter := bson.M{
		"id": id,
	}
	var result_data []student
	cursor, err := s_col.Find(context.TODO(), filter)
	handleError(err)
	cursor.All(context.Background(), &result_data)
	if len(result_data) != 0 {
		w.Write([]byte("User already exist"))
		return
	}
	s_col.InsertOne(context.TODO(), stud)
	fmt.Println(stud)
	w.Write([]byte("User added successfully"))
}
func updateStd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var stud student
	json.NewDecoder(r.Body).Decode(&stud)
	id := stud.Id
	filter := bson.M{
		"id": id,
	}
	update := bson.D{{"$set", bson.D{{"name", "hielo"}}}}
	_, err := s_col.UpdateOne(context.Background(), filter, update)
	handleError(err)
	var result_data []student
	cursor, err := s_col.Find(context.TODO(), filter)
	handleError(err)
	cursor.All(context.Background(), &result_data)
	fmt.Println(stud)
	w.Write([]byte("User added successfully"))
}
func reqLeave(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var leave leave_data
	json.NewDecoder(r.Body).Decode(&leave)
	id := leave.Id
	filter := bson.M{
		"id": id,
	}
	leave1 := &leave_data{
		Id:     leave.Id,
		Reason: leave.Reason,
		Frm:    leave.Frm,
		To:     leave.To,
		Status: "false",
	}
	var result_data []leave_data
	cursor, err := l_col.Find(context.TODO(), filter)
	handleError(err)
	cursor.All(context.Background(), &result_data)
	l_col.InsertOne(context.TODO(), leave1)
	fmt.Println(leave)
	w.Write([]byte("request sent successfully"))
}

func Approve(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var l leave_data
	json.NewDecoder(r.Body).Decode(&l)
	id := l.Id
	filter := bson.M{
		"id": id,
	}
	update := bson.D{{"$set", bson.D{{"status", "true"}}}}
	_, err := l_col.UpdateOne(context.Background(), filter, update)
	handleError(err)
	var result_data []leave_data
	cursor, err := l_col.Find(context.TODO(), filter)
	handleError(err)
	cursor.All(context.Background(), &result_data)
	fmt.Println(l)
	w.Write([]byte("req approved successfully"))
}
func allLeaves(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var leaves []leave_data
	cursor, err := l_col.Find(context.TODO(), bson.D{{}})
	handleError(err)
	cursor.All(context.Background(), &leaves)
	fmt.Println(leaves)
	for _, val := range leaves {
		s := fmt.Sprintf(" %d %s %s %s %s\n", val.Id, val.Reason, val.Frm, val.To, val.Status)
		w.Write([]byte(s))
	}
}
