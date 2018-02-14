package sandbox

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

// Seperator is used to delineate inputs when running code through the sandbox
// Seperator must be maintained in the scripts in /payload/ if the values don't match things break
const Seperator = "\n*-BRK-*\n"

// Interface provides an interface to interact with the sandbox
type Interface struct {
	LanguageMap map[string]Language
}

// Message ...
type Message struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

// New creates a new sandbox handler that can compile code and return results
func New(languagesFile string) Interface {
	bytes, err := ioutil.ReadFile(languagesFile)
	if err != nil {
		log.Fatal(err)
	}

	languageMap := make(map[string]Language, 0)
	err = json.Unmarshal(bytes, &languageMap)
	if err != nil {
		log.Fatal(err)
	}

	return Interface{languageMap}
}

// input is n test calls seperated by newlines
// input and expected MUST end in newlines
func (t Interface) run(language, code, stdinGlob string) (string, Message) {
	log.Printf("sandbox launching sandbox...\nLanguage: %s\nStdin: %sCode: Hidden\n", language, stdinGlob)
	lang, ok := t.LanguageMap[strings.ToLower(language)]
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

// EvalWithStdins evaluates 'code' with each of 'stdins'
func (t Interface) EvalWithStdins(language, code string, stdins []string) ([]string, Message) {
	stdinGlob := glob(stdins)
	results, msg := t.run(language, code, stdinGlob)

	return unglob(results), msg
}

// glob and unglob combine and seperate groups of input or output (compiler deals with globs of text seperated by separator but outside compilation a []string is preferred)
func glob(stdins []string) string {
	glob := strings.Join(stdins, Seperator) + Seperator
	return glob
}

func unglob(glob string) []string {
	stdins := strings.Split(glob, Seperator)
	return stdins[:len(stdins)-1]
}
