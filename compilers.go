package xaqt

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// Compilers maps language names to the details of how to execute code in that language.
type Compilers map[string]ExecutionDetails

// ExecutionDetails specifies how to execute certain code.
type ExecutionDetails struct {
	Compiler           string `json:"compiler"`
	SourceFile         string `json:"sourceFile"`
	OptionalExecutable string `json:"optionalExecutable"`
	CompilerFlags      string `json:"compilerFlags"`
	Disabled           string `json:"disabled"`
}

// Reads a compilers map from a file.
func ReadCompilers(filename string) Compilers {
	compilerMap := make(Compilers, 0)

	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Failed to read language file: %s", err)
	}

	err = json.Unmarshal(bytes, &compilerMap)
	if err != nil {
		log.Fatalf("Failed to parse JSON: %s", err)
	}

	return compilerMap
}

// availableLanguages returns a list of currently supported languages.
func (c Compilers) availableLanguages() []string {
	langs := make([]string, 0)

	// make a list of currently supported languages
	for k, v := range c {
		if v.Disabled != "true" {
			langs = append(langs, k)
		}
	}

	log.Printf("currently supporting %d of %d known languages\n", len(langs), len(c))
	return langs
}
