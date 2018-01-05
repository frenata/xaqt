package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testbox"
)

// type TestResponse struct {
// 	Id          id     `json:"id"`
// 	Description string `json:"description"`
// 	SampleIO    string `json:"sampleIO"`
// 	// ShortName   string   `json:"shortName"`
// 	// Tags        []string `json:"tags"`
// }

type SubmissionRequest struct {
	Id       id     `json:"id"`
	Language string `json:"language"`
	Code     string `json:"code"`
	Input    string `json:"input"`
}

func (s SubmissionRequest) String() string {
	return fmt.Sprintf("( <SubmissionRequest> {ID: %s, Language: %s, Code: Hidden, Input: %s} )", s.Id, s.Language, s.Input)
}

type CompileResult struct {
	Raw     string            `json:"raw"`
	Graded  map[string]string `json:"graded"`
	Message testbox.Message   `json:"message"`
}

type LanguagesResponse struct {
	Languages map[string]testbox.Language `json:"languages"`
}

var box testbox.TestBox

func main() {
	port := os.Getenv("TEST_BOX_PORT")

	box = testbox.New("data/compilers.json")

	http.HandleFunc("/", getChallenge)
	http.HandleFunc("/submit/", submitTest)
	http.HandleFunc("/stdout/", getStdout)
	http.HandleFunc("/languages/", getLangs)

	log.Println("TestBox listening on " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func getChallenge(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request for test...")

	// rand.Seed(time.Now().UTC().UnixNano())
	// n := rand.Intn(len(challenges))

	// temporary hack to check multi-line:
	challengeID := "1"
	// challengeID := strconv.Itoa(n)
	challenge := challenges[challengeID]

	json, _ := json.MarshalIndent(challenge, "", "    ")

	log.Printf("Handing out test ID: %s, Desc:%s\n", challenge.ID, challenge.Description)

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func getStdout(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request for stdout")
	decoder := json.NewDecoder(r.Body)
	var submission SubmissionRequest
	err := decoder.Decode(&submission)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	output, msg := box.CompileAndPrint(submission.Language, submission.Code, submission.Input)
	log.Println(output, msg)

	if output == "" {
		output = "NO OUTPUT"
	}

	buf, _ := json.MarshalIndent(CompileResult{
		Raw:     output,
		Message: msg,
	}, "", "   ")

	w.Header().Set("Content-Type", "application/json")
	w.Write(buf)
}

func submitTest(w http.ResponseWriter, r *http.Request) {
	log.Println("Received challenge submission")
	decoder := json.NewDecoder(r.Body)
	var submission SubmissionRequest
	err := decoder.Decode(&submission)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	log.Printf("submitTest, submission: %s\n", submission)

	if submission.Id == "" {
		// lack of response crashes client
		log.Panic("Submission has no challengeID")
		return
	}

	challenge := challenges[submission.Id]
	stdin, stdout := challenge.getCases()
	log.Printf("submitTest, challenge: %v\n", challenge)

	passed, msg := box.CompileAndChallenge(submission.Language, submission.Code, stdin, stdout)
	log.Println(passed, msg)

	buf, _ := json.MarshalIndent(CompileResult{
		Graded:  passed,
		Message: msg,
	}, "", "   ")

	w.Header().Set("Content-Type", "application/json")
	w.Write(buf)
}

func getLangs(w http.ResponseWriter, r *http.Request) {
	log.Println("Received languages request")
	langs := make(map[string]testbox.Language)

	for k, v := range box.LanguageMap {
		langs[k] = testbox.Language{Boilerplate: v.Boilerplate, CommentPrefix: v.CommentPrefix}
	}

	// add boilerplate and comment info
	log.Println(langs)
	buf, _ := json.MarshalIndent(LanguagesResponse{langs}, "", "   ")

	w.Header().Set("Content-Type", "application/json")
	w.Write(buf)
}

var challenges map[id]Challenge

func init() {
	log.Println("Reading challenges file...")
	bytes, err := ioutil.ReadFile("data/challenges.json")
	if err != nil {
		panic(err)
	}

	challenges = make(map[id]Challenge)
	err = json.Unmarshal(bytes, &challenges)
	if err != nil {
		panic(err)
	}
	log.Println("Challenges file loaded.")
	// for k, v := range challenges {
	// 	log.Printf("ID: %s maps to %s", k, v.ID)
	// }
}
