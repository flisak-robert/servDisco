package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

type Label struct {
	Job             string `json:"job"`
	Env             string `json:"env"`
	MetricsEndpoint string `json:"__metrics_path__"`
}

type Target struct {
	Targets []string `json:"targets"`
	Labels  Label    `json:"labels"`
}

var target_json = "targets.json"
var target_ssl_json = "targets-ssl.json"

var infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Service discovery home page")
	fmt.Println("Endpoint Hit: homePage")
}

func handleRequests() {

	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/add/{key}", addNewTarget).Methods("POST")
	myRouter.HandleFunc("/targets/{key}", getAllTargets)
	log.Println("Starting server on :4000")
	log.Fatal(http.ListenAndServe(":4000", myRouter))
}

func openJson(r *http.Request) string {
	filename_from_url := strings.Split(r.URL.Path, "/")
	filename := filename_from_url[len(filename_from_url)-1]
	return filename
}

func getAllTargets(w http.ResponseWriter, r *http.Request) {
	filename := openJson(r)
	var file_to_open string
	if filename == "targets" {
		file_to_open = target_json
	} else if filename == "targets-ssl" {
		file_to_open = target_ssl_json
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "404 Not found")
		return
	}

	jsonFile, err := os.Open(file_to_open)

	if err != nil {
		fmt.Println(err)
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var Targets []Target
	json.Unmarshal(byteValue, &Targets)
	infoLog.Printf("Getting all targets: %v", Targets)

	json.NewEncoder(w).Encode(Targets)
	defer jsonFile.Close()
}

func addNewTarget(w http.ResponseWriter, r *http.Request) {
	filename := openJson(r)
	var file_to_open string
	if filename == "targets" {
		file_to_open = target_json
	} else if filename == "targets-ssl" {
		file_to_open = target_ssl_json
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "404 Not found")
		return
	}

	jsonFile, err := os.OpenFile(file_to_open, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	var Targets []Target
	byteValue, err := ioutil.ReadFile(file_to_open)
	json.Unmarshal(byteValue, &Targets)
	if err != nil {
		log.Fatal(err)
	}

	var t Target
	err = json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for counter, b := range Targets {
		for _, i := range b.Targets {
			if i == t.Targets[0] {
				w.WriteHeader(http.StatusOK)
				fmt.Fprintf(w, "Target already exists")
				infoLog.Printf("Target: %s already exists!", i)
				return
			}
		}

		if b.Labels.Job == t.Labels.Job && b.Labels.Env == t.Labels.Env {

			b.Targets = append(b.Targets, t.Targets...)

			p := &Targets[counter]
			p.Targets = b.Targets

			writeJsonToFile(Targets, jsonFile)
			fmt.Fprintf(w, "Successfully added new target")
			infoLog.Printf("%s successfully added new target to job: %s, env: %s", b.Targets[len(b.Targets)-1], b.Labels.Job, b.Labels.Env)

			return
		}
	}

	Targets = append(Targets, t)

	writeJsonToFile(Targets, jsonFile)
	fmt.Fprintf(w, "Successfully added new target")
	infoLog.Printf("%v successfully added to targets.json", t)

}

func writeJsonToFile(jsonArray []Target, file *os.File) {

	indented_data, _ := json.MarshalIndent(jsonArray, "", " ")
	file.Write(indented_data)
	defer file.Close()

}

func main() {

	handleRequests()

}
