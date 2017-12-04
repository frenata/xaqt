package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testbox"
)

type TestResponse struct {
	Id          id     `json:"id"`
	Description string `json:"description"`
}

type SubmissionRequest struct {
	Id       id     `json:"id"`
	Language string `json:"language"`
	Code     string `json:"code"`
}

type SubmissionResponse struct {
	PassedTests map[string]bool `json:"passedTests"`
}

var box testbox.TestBox

func main() {
	port := os.Getenv("TEST_BOX_PORT")

	box = testbox.New("data/compilers.json")

	http.HandleFunc("/", getTest)
	http.HandleFunc("/submit/", submitTest)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func getTest(w http.ResponseWriter, r *http.Request) {
	var test Test
	var testid id

	// ranging as a quick way to get a random map entry
	for k, v := range challenges {
		testid = k
		test = v
		break
	}

	tr := TestResponse{testid, test.Description}
	json, _ := json.MarshalIndent(tr, "", "    ")

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func submitTest(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var submission SubmissionRequest
	err := decoder.Decode(&submission)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	test := challenges[submission.Id]
	stdin, stdout := test.StdIO()

	passed := box.Test(submission.Language, submission.Code, stdin, stdout)
	log.Println(passed)

	json, _ := json.MarshalIndent(SubmissionResponse{passed}, "", "   ")

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(json))
}

var challenges map[id]Test

func init() {
	bytes, err := ioutil.ReadFile("data/challenges.json")
	if err != nil {
		panic(err)
	}

	challenges = make(map[id]Test)
	err = json.Unmarshal(bytes, &challenges)
	if err != nil {
		panic(err)
	}

}
