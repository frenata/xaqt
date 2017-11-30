package main

import (
	"encoding/json"
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

func main() {
	port := os.Getenv("TEST_BOX_PORT")

	http.HandleFunc("/", getTest)
	http.HandleFunc("/submit/", submitTest)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func getTest(w http.ResponseWriter, r *http.Request) {
	var test Test
	var testid id

	// ranging as a quick way to get a random map entry
	for k, v := range tests {
		testid = k
		test = v
		break
	}

	tr := TestResponse{testid, test.description}
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

	test := tests[submission.Id]
	stdin, stdout := test.StdIO()

	passed := testbox.Test(submission.Language, submission.Code, stdin, stdout)
	log.Println(passed)

	json, _ := json.MarshalIndent(SubmissionResponse{passed}, "", "   ")

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(json))
}

var tests map[id]Test

func init() {
	// later read from a file
	tests = make(map[id]Test)
	tests["1"] = Test{"echo", map[string]string{"": "", "test": "test"}}
}
