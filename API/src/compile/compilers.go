package compile

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

var LanguageMap map[string]Language

type Language struct {
	Compiler           string `json:"compiler"`
	SourceFile         string `json:"sourceFile"`
	OptionalExecutable string `json:"optionalExecutable"`
	CompilerFlags      string `json:"compilerFlags"`
}

func init() {
	bytes, err := ioutil.ReadFile("data/compilers.json")
	if err != nil {
		log.Fatal(err)
	}

	LanguageMap = make(map[string]Language, 0)
	err = json.Unmarshal(bytes, &LanguageMap)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(LanguageMap)
}
