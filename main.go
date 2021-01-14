package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Name  string
	Email string
}

func initialMigration() {
	db, err := gorm.Open("mysql", "root:example@tcp(192.168.0.102:3306)/dbpm?charset=utf8&parseTime=True")
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&User{})
}

func allUsers(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("mysql", "root:example@tcp(192.168.0.102:3306)/dbpm?charset=utf8&parseTime=True")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	var users []User
	db.Find(&users)
	fmt.Println("{}", users)

	json.NewEncoder(w).Encode(users)
}

func newUser(w http.ResponseWriter, r *http.Request) {
	type UserJson struct {
		Name  string `json:"Name"`
		Email string `json:"Email"`
	}

	fmt.Println("New User Endpoint Hit")

	db, err := gorm.Open("mysql", "root:example@tcp(192.168.0.102:3306)/dbpm?charset=utf8&parseTime=True")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	reqBody, _ := ioutil.ReadAll(r.Body)
	var userjson UserJson
	json.Unmarshal(reqBody, &userjson)

	db.Create(&User{Name: userjson.Name, Email: userjson.Email})
	fmt.Fprintf(w, "New User Successfully Created")
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	type UserJson struct {
		Id string `json:"Id"`
	}

	db, err := gorm.Open("mysql", "root:example@tcp(192.168.0.102:3306)/dbpm?charset=utf8&parseTime=True")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	reqBody, _ := ioutil.ReadAll(r.Body)
	var userjson UserJson
	json.Unmarshal(reqBody, &userjson)

	var user User
	db.Where("id = ?", userjson.Id).Find(&user)
	db.Delete(&user)

	fmt.Fprintf(w, "Successfully Deleted User")
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	type UserJson struct {
		Id    string `json:"Id"`
		Name  string `json:"Name"`
		Email string `json:"Email"`
	}

	db, err := gorm.Open("mysql", "root:example@tcp(192.168.0.102:3306)/dbpm?charset=utf8&parseTime=True")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	reqBody, _ := ioutil.ReadAll(r.Body)
	var userjson UserJson
	json.Unmarshal(reqBody, &userjson)

	var user User
	db.Where("id = ?", userjson.Id).Find(&user)

	user.Email = userjson.Email
	user.Name = userjson.Name

	db.Save(&user)
	fmt.Fprintf(w, "Successfully Updated User")
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/users", allUsers).Methods("GET")
	myRouter.HandleFunc("/user", deleteUser).Methods("DELETE")
	myRouter.HandleFunc("/user", updateUser).Methods("PUT")
	myRouter.HandleFunc("/user", newUser).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", myRouter))
}

func main() {
	fmt.Println("Go ORM Tutorial")

	// Add the call to our new initialMigration function
	initialMigration()

	handleRequests()
}
