package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/frenata/xaqt"
)

type CodeSubmission struct {
	Language string   `json:"language"`
	Code     string   `json:"code"`
	Stdins   []string `json:"stdins"`
}

func (s CodeSubmission) String() string {
	return fmt.Sprintf("( <CodeSubmission> {Language: %s, Code: Hidden, Stdins: %s} )", s.Language, s.Stdins)
}

type ExecutionResult struct {
	Stdouts []string           `json:"stdouts"`
	Message xaqt.Message `json:"message"`
}

// TODO: move into main rather than a global
var box xaqt.Compilers

func main() {
	port := getEnv("COMPILEBOX_PORT", "31337")

	box = xaqt.New("data/compilers.json")

	http.HandleFunc("/languages/", getLangs)
	http.HandleFunc("/evaluate/", evalCode)

	log.Println("xaqt listening on " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	log.Printf("Environment variable %s not found, setting to %s", key, fallback)
	os.Setenv(key, fallback)
	return fallback
}

func evalCode(w http.ResponseWriter, r *http.Request) {
	log.Println("Received code subimssion...")
	decoder := json.NewDecoder(r.Body)
	var submission CodeSubmission
	err := decoder.Decode(&submission)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	//log.Println(submission)
	stdouts, msg := box.Evaluate(submission.Language, submission.Code, submission.Stdins)
	log.Println(stdouts, msg)

	if len(stdouts) == 0 {
		log.Println("Code produced no output")
		stdouts = append(stdouts, "ZERO OUTPUTS")
	}

	buf, _ := json.MarshalIndent(ExecutionResult{
		Stdouts: stdouts,
		Message: msg,
	}, "", "   ")

	w.Header().Set("Content-Type", "application/json")
	w.Write(buf)
}

func getLangs(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received languages request...")

	workingLangs := box.AvailableLanguages()
	log.Printf("currently supporting %d of %d known languages\n", len(workingLangs), len(box))

	// encode language list
	buf, _ := json.MarshalIndent(workingLangs, "", "   ")

	// write working language list back to client
	w.Header().Set("Content-Type", "application/json")
	w.Write(buf)
}
