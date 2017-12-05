package testbox

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

type TestBox struct {
	languageMap map[string]Language
}

type Message struct {
	Type   string `json:"type"`
	Sender string `json:"sender"`
	Data   string `json:"data"`
}

func New(languagesFile string) TestBox {
	bytes, err := ioutil.ReadFile(languagesFile)
	if err != nil {
		log.Fatal(err)
	}

	languageMap := make(map[string]Language, 0)
	err = json.Unmarshal(bytes, &languageMap)
	if err != nil {
		log.Fatal(err)
	}

	return TestBox{languageMap}
}

// input is n test calls seperated by newlines
// input and expected MUST end in newlines
func (t TestBox) Test(language, code, input, expected string) (map[string]bool, Message) {
	lang := t.languageMap[language]

	if code == "" {
		return nil, Message{"error", "testBox", "no code submitted"}
	}

	sb := NewSandbox(lang, code, input, DefaultSandboxOptions())

	output, err := sb.Run()
	if err != nil {
		log.Println(err)
		return nil, Message{"error", "testBox", fmt.Sprintf("%s", err)}
	}

	splitOutput := strings.SplitN(output, "*-COMPILEBOX::ENDOFOUTPUT-*", 2)
	timeTaken := splitOutput[1]
	_ = timeTaken
	result := splitOutput[0]

	return compareLineByLine(input, expected, result), Message{"success", "testBox", "compilation took " + timeTaken + " seconds"}
}

func compareLineByLine(input, exp, res string) map[string]bool {
	inpSlice := strings.Split(input, "\n")
	expSlice := strings.Split(exp, "\n")
	resSlice := strings.Split(res, "\n")

	results := make(map[string]bool, len(expSlice))

	// TODO deal with partial success but incorrect result couont
	/*if len(expSlice) != len(resSlice) {
		return results
	}*/

	for i := 0; i < len(inpSlice)-1; i++ {
		//log.Println("compare: ", inpSlice[i], expSlice[i], resSlice[i])
		results[inpSlice[i]] = expSlice[i] == resSlice[i]
	}

	return results
}
