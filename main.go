package main

import (
	"net/http"
	"fmt"
	"log"
	"github.com/gorilla/mux"
	"github.com/gorilla/context"
	"io/ioutil"
	"encoding/json"
	"os"
	"strings"
)

type Potato struct {
	Text    string  `json:"text"`
	History []Entry `json:"history"`
}

type Entry struct {
	Node string `json:"node"`
	Text string `json:"text"`
	Desc string `json:"desc"`
}

func work(w http.ResponseWriter, r *http.Request)  {
	defer r.Body.Close()
	body,err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	var data Potato

	err = json.Unmarshal(body, &data)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Could not unmarshall that shit"))
		return
	}

	newData := Entry{}
	newData.Desc = "just asking simon"
	newData.Node,err = os.Hostname()
	if err != nil{
		w.WriteHeader(500)
		w.Write([]byte("Cannot get hostname"))
		return
	}

	splited := strings.Split(data.Text," ")
	length := len(splited)
	if length % 2 != 0 {
		newData.Text = "simon says " + data.Text
	} else {
		newData.Text = "simon does not like you"
	}

	data.History = append(data.History,newData)
	data.Text = newData.Text

	outputBytes, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("WTF cannot marshal"))
		return
	}

	w.WriteHeader(200)
	w.Write(outputBytes)
	return
}

	func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/process", work).Methods("POST")
	http.Handle("/", router)
	fmt.Printf("Server ready\n")
	log.Fatal(http.ListenAndServe(":80",context.ClearHandler(router)))
}
