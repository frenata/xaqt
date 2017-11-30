package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"testbox"
)

type TestRequest struct {
	Language string `json:"language"`
	Code     string `json:"code"`
	Stdin    string `json:"stdin"`
	Expected string `json:"expected"`
}

type TestResponse struct {
	PassedTests []bool `json:"passedTests"`
}

func main() {
	port := os.Getenv("TEST_BOX_PORT")

	http.HandleFunc("/", tester)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func tester(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var req TestRequest
	err := decoder.Decode(&req)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	passed := testbox.Test(req.Language, req.Code, req.Stdin, req.Expected)
	log.Println(passed)

	json, _ := json.MarshalIndent(TestResponse{passed}, "", "   ")

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(json))
}
