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
	Input    string `json:"input"`
}

type SubmissionResponse struct {
	Output map[string]string `json:"passFail"`
	Error  testbox.Message   `json:"error"`
}

var box testbox.TestBox

func main() {
	port := os.Getenv("TEST_BOX_PORT")

	box = testbox.New("data/compilers.json")

	http.HandleFunc("/", getTest)
	http.HandleFunc("/submit/", submitTest)
	http.HandleFunc("/stdout/", getStdout)
	// TODO: add languages endpoint

	log.Println("TestBox listening on " + port)
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

func getStdout(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var submission SubmissionRequest
	err := decoder.Decode(&submission)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	output, msg := box.StdOut(submission.Language, submission.Code, submission.Input)
	log.Println(output, msg)

	buf, _ := json.MarshalIndent(SubmissionResponse{output, msg}, "", "   ")

	w.Header().Set("Content-Type", "application/json")
	w.Write(buf)
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

	passed, msg := box.Test(submission.Language, submission.Code, stdin, stdout)
	log.Println(passed, msg)

	buf, _ := json.MarshalIndent(SubmissionResponse{passed, msg}, "", "   ")

	w.Header().Set("Content-Type", "application/json")
	w.Write(buf)
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
