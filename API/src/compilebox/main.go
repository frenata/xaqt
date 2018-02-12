package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sandbox"
)

type CodeSubmission struct {
	Language string   `json:"language"`
	Code     string   `json:"code"`
	Stdins   []string `json:"input"`
}

func (s CodeSubmission) String() string {
	return fmt.Sprintf("( <CodeSubmission> {Language: %s, Code: Hidden, Stdins: %s} )", s.Language, s.Stdins)
}

type ExecutionResult struct {
	Stdouts []string        `json:"stdouts"`
	Message sandbox.Message `json:"message"`
}

type LanguagesResponse struct {
	Languages map[string]sandbox.Language `json:"languages"`
}

var box sandbox.Interface

func main() {
	port := getEnv("COMPILE_PORT", "31337")

	box = sandbox.New("data/compilers.json")

	http.HandleFunc("/languages/", getLangs)
	http.HandleFunc("/eval/", evalCode)

	log.Println("testbox listening on " + port)
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

	log.Printf("...along with %d stdin inputs", len(submission.Stdins))

	stdouts, msg := box.EvalWithStdins(submission.Language, submission.Code, submission.Stdins)
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
	log.Println("Received languages request")
	langs := make(map[string]sandbox.Language)

	for k, v := range box.LanguageMap {
		langs[k] = sandbox.Language{Boilerplate: v.Boilerplate, CommentPrefix: v.CommentPrefix}
	}

	// add boilerplate and comment info
	log.Println(langs)
	buf, _ := json.MarshalIndent(LanguagesResponse{langs}, "", "   ")

	w.Header().Set("Content-Type", "application/json")
	w.Write(buf)
}
