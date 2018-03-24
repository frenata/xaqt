package xaqt

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

// Compilers maps language names to the details of how to execute code in that language.
type Compilers map[string]ExecutionDetails

// ExecutionDetails specifies how to execute certain code.
type ExecutionDetails struct {
	Compiler           string //`json:"compiler"`
	SourceFile         string //`json:"sourceFile"`
	OptionalExecutable string //`json:"optionalExecutable"`
	CompilerFlags      string //`json:"compilerFlags"`
	Disabled           string //`json:"disabled"`
}

// Message represents details on success or failure of execution.
type Message struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

// New reads a JSON-encoded file of compilers.
func New(filename string) Compilers {
	compilers := make(Compilers, 0)

	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Failed to read language file: %s", err)
	}

	err = json.Unmarshal(bytes, &compilers)
	if err != nil {
		log.Fatalf("Failed to parse JSON: %s", err)
	}

	return compilers
}

// Evaluate code in a given language and for a set of 'stdin's.
func (c Compilers) Evaluate(language, code string, stdins []string) ([]string, Message) {
	stdinGlob := glob(stdins)
	results, msg := c.run(language, code, stdinGlob)

	return unglob(results), msg
}

// AvailableLanguages returns a list of currently supported languages.
func (c Compilers) AvailableLanguages() []string {
	langs := make([]string, 0)

	// make a list of currently supported languages
	for k, v := range c {
		if v.Disabled != "true" {
			langs = append(langs, k)
		}
	}

	return langs
}

// input is n test calls seperated by newlines
// input and expected MUST end in newlines
func (c Compilers) run(language, code, stdinGlob string) (string, Message) {
	log.Printf("sandbox launching sandbox...\nLanguage: %s\nStdin: %sCode: Hidden\n", language, stdinGlob)
	lang, ok := c[strings.ToLower(language)]
	if !ok || lang.Disabled == "true" {
		return "", Message{"error", "language not supported"}
	}

	if code == "" {
		return "", Message{"error", "no code submitted"}
	}

	sb := NewSandbox(lang, code, stdinGlob, DefaultSandboxOptions())

	output, err := sb.Run()
	if err != nil {
		log.Printf("sandbox run error: %v", err)
		return output, Message{"error", fmt.Sprintf("%s", err)}
	}

	splitOutput := strings.SplitN(output, "*-COMPILEBOX::ENDOFOUTPUT-*", 2)
	timeTaken := splitOutput[1]
	result := splitOutput[0]

	return result, Message{"success", "compilation took " + timeTaken + " seconds"}
}
