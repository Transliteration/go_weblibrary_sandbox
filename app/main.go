package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"weblibrary_sandbox/config"

	"github.com/google/uuid"

	log "github.com/sirupsen/logrus"
)

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	io.WriteString(w, "This is my website!\n")
}

// Запрос:
// 1. добавить пользователя
func addUserHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /addUser request %s\n", r.Method)

	err := r.ParseForm()
	if err != nil {
		log.Error(trace())
		log.Error("ParseForm failed in %s", trace().funcName)
		// in case of any error
		return
	}

	name := r.Form.Get("name") // x will be "" if parameter is not set
	fmt.Println(name)
	ageString := r.Form.Get("age") // age will be "" if parameter is not set
	fmt.Println(ageString)

	age, err := strconv.ParseInt(ageString, 10, 64)
	if err != nil {
		log.Error(trace())
		log.Error("Could not parse to int user age = %v", ageString)
		///
		return
	}

	newUser := User{
		name:   name,
		userID: uuid.New(),
		age:    int(age),
	}

	err = addUser(newUser)
	if err != nil {
		log.Error(trace())
		log.Error("Could not add new user %v", newUser)
		///
		return
	}

	io.WriteString(w, "New user added successfully\n")
}

// 2. добавить ему книгу по айди (Не могут зафейлиться со стороны юзера)

// Запрос:
// 1. получить пользователя по айди (могут зафейлиться, похуй)

func getUserHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /getUser request %s\n", r.Method)

	err := r.ParseForm()
	if err != nil {
		log.Error(trace())
		log.Error("ParseForm failed")
		return
	}

	userIDField := r.Form.Get("userID") // x will be "" if parameter is not set
	userID, err := uuid.Parse(userIDField)
	if err != nil {
		log.Error(trace())
		log.Error("Could not parse user UUID %s", userIDField)
		return
	}

	user, err := getUser(userID)
	if err != nil {
		log.Error(trace())
		log.Error("Could not get user with ID = %v", userID)
		return
	}

	sas := config.Kek
	fmt.Println(sas)

	responseMessage := fmt.Sprintf("Your user name is {%s}\n", user.name)
	io.WriteString(w, responseMessage)
}

func main() {
	http.HandleFunc("/", getRoot)

	http.HandleFunc("/addUser/", addUserHandle)
	http.HandleFunc("/getUser/", getUserHandle)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
