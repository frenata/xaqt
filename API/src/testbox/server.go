package testbox

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
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
func (t TestBox) run(language, code, input string) (string, Message) {
	lang, ok := t.languageMap[language]
	if !ok {
		return "", Message{"error", "testBox", "language not recognized"}
	}

	if code == "" {
		return "", Message{"error", "testBox", "no code submitted"}
	}

	sb := NewSandbox(lang, code, input, DefaultSandboxOptions())

	output, err := sb.Run()
	if err != nil {
		log.Println(err)
		return "", Message{"error", "testBox", fmt.Sprintf("%s", err)}
	}

	splitOutput := strings.SplitN(output, "*-COMPILEBOX::ENDOFOUTPUT-*", 2)
	timeTaken := splitOutput[1]
	result := splitOutput[0]

	return result, Message{"success", "testBox", "compilation took " + timeTaken + " seconds"}
}

func (t TestBox) StdOut(language, code, input string) (map[string]string, Message) {
	if !strings.HasSuffix(input, "\n") {
		input = input + "\n"
	}
	result, msg := t.run(language, code, input)

	return mapInToOut(input, result), msg
}

func (t TestBox) Test(language, code, input, expected string) (map[string]string, Message) {
	if !strings.HasSuffix(input, "\n") {
		input = input + "\n"
	}
	result, msg := t.run(language, code, input)

	return compareLineByLine(input, expected, result), msg
}

func compareLineByLine(input, exp, res string) map[string]string {
	inpSlice := strings.Split(input, "\n")
	expSlice := strings.Split(exp, "\n")
	resSlice := strings.Split(res, "\n")

	results := make(map[string]string)

	// TODO: remove for prod!
	if strings.HasPrefix(resSlice[0], "godmode ") {
		nStr := strings.TrimPrefix(resSlice[0], "godmode ")
		n, e := strconv.Atoi(nStr)
		if e != nil {
			log.Println(e)
			panic("Bad dog!")
		}
		if n > len(inpSlice) {
			n = len(inpSlice)
		}

		i := 0
		log.Println(i, n)
		for ; i < n; i++ {
			results[inpSlice[i]] = "true"
		}
		for ; i < len(inpSlice)-1; i++ {
			results[inpSlice[i]] = "false"
		}

		return results
	}

	// TODO deal with partial success but incorrect result couont
	/*if len(expSlice) != len(resSlice) {
		return results
	}*/

	for i := 0; i < len(inpSlice)-1; i++ {
		//log.Println("compare: ", inpSlice[i], expSlice[i], resSlice[i])
		if i > len(expSlice)-1 || i > len(resSlice)-1 {
			results[inpSlice[i]] = "false"
		} else {
			results[inpSlice[i]] = fmt.Sprintf("%v", expSlice[i] == resSlice[i])
		}
	}

	return results
}

func mapInToOut(input, res string) map[string]string {
	inpSlice := strings.Split(input, "\n")
	resSlice := strings.Split(res, "\n")

	results := make(map[string]string, len(inpSlice))

	// TODO deal with partial success but incorrect result couont
	/*if len(expSlice) != len(resSlice) {
		return results
	}*/

	for i := 0; i < len(inpSlice); i++ {
		//log.Println("compare: ", inpSlice[i], expSlice[i], resSlice[i])
		if i > len(resSlice)-1 {
			results[inpSlice[i]] = "NO OUTPUT"
		} else {
			results[inpSlice[i]] = resSlice[i]
		}
	}

	return results
}
