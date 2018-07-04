package xaqt

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

// Compilers maps language names to the details of how to execute code in that language.
type Compilers map[string]CompilerDetails

// ExecutionDetails specifies how to execute certain code.
type ExecutionDetails struct {
	Compiler           string `json:"compiler"`
	SourceFile         string `json:"sourceFile"`
	OptionalExecutable string `json:"optionalExecutable"`
	CompilerFlags      string `json:"compilerFlags"`
	Disabled           string `json:"disabled"`
}

// CompositionDetails specifies how to write code in a given language
type CompositionDetails struct {
	Boilerplate   string `json:"boilerplate"`
	CommentPrefix string `json:"commentPrefix"`
}

// CompilerDetails contains everything XAQT knows about handling a certain language
type CompilerDetails struct {
	ExecutionDetails
	CompositionDetails
}

// ReadCompilers reads a compilers map from a file.
func ReadCompilers(filename string) Compilers {
	compilerMap := make(Compilers, 0)

	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Fatal: failed to read language file: %s", err)
	}

	err = json.Unmarshal(bytes, &compilerMap)
	if err != nil {
		log.Fatalf("Fatal: failed to parse JSON: %s", err)
	}

	return compilerMap
}

// availableLanguages returns a list of currently supported languages.
func (c Compilers) availableLanguages() map[string]CompositionDetails {
	fmt.Printf("Received languages request...")
	langs := make(map[string]CompositionDetails)

	// make a list of currently supported languages
	for k, v := range c {
		if v.Disabled != "true" {
			langs[k] = v.CompositionDetails
		}
	}

	log.Printf("currently supporting %d of %d known languages\n", len(langs), len(c))
	return langs
}
