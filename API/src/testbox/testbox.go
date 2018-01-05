package testbox

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

const Seperator = "\n*-BRK-*\n"

type TestBox struct {
	LanguageMap map[string]Language
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
	log.Printf("TestBox run called with input: %v", input)
	lang, ok := t.LanguageMap[strings.ToLower(language)]
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

// func (t TestBox) CompileAndPrint(language, code, input string) (map[string]string, Message) {
func (t TestBox) CompileAndPrint(language, code, input string) (string, Message) {
	input = input + Seperator
	result, msg := t.run(language, code, input)

	// input = strings.Split(input, Seperator)[0]
	result = strings.Split(result, Seperator)[0]

	log.Printf("CompileAndPrint result: %v", result)
	return result, msg
	// return mapInToOut(input, result), msg
}

func (t TestBox) CompileAndChallenge(language, code, input, expected string) (map[string]string, Message) {

	result, msg := t.run(language, code, input)

	return compareBlockByBlock(input, expected, result), msg
}

func compareBlockByBlock(input, exp, res string) map[string]string {
	inpSlice := strings.Split(input, Seperator)
	expSlice := strings.Split(exp, Seperator)
	resSlice := strings.Split(res, Seperator)

	results := make(map[string]string)
	// log.Printf("compBbyb, input: %v", input)
	// log.Printf("compBbyb, exp: %v", exp)
	// log.Printf("compBbyb, res: %v", res)

	// log.Printf("compBbyb, inpSlice: %v", inpSlice)
	// log.Printf("compBbyb, expSlice: %v", expSlice)
	// log.Printf("compBbyb, resSlice: %v", resSlice)

	// TODO deal with partial success but incorrect result couont
	/*if len(expSlice) != len(resSlice) {
		return results
	}*/

	for i := 0; i < len(inpSlice)-1; i++ {

		if i > len(expSlice)-1 || i > len(resSlice)-1 {
			results[inpSlice[i]] = "Fail"
		} else {
			log.Printf("Input:\n%v\nOutput:\n%v\nResult:\n%v\n", inpSlice[i], expSlice[i], resSlice[i])
			results[inpSlice[i]] = passFail(expSlice[i], resSlice[i])
		}
	}

	return results
}

func passFail(a, b string) string {
	if a == b {
		return "Pass"
	}
	return "Fail"
}

func mapInToOut(input, res string) map[string]string {
	inpSlice := strings.Split(input, Seperator)
	resSlice := strings.Split(res, Seperator)

	results := make(map[string]string, len(inpSlice))

	// TODO deal with partial success but incorrect result count

	// log.Printf("Input: %v\nSliced: %v\nRes:%v\nSliced:%v\nlen(inpSlice):%v", input, inpSlice, res, resSlice, len(inpSlice))

	for i := range inpSlice {
		// log.Printf("resSlice: %s\n", resSlice[i])
		// log.Printf("resSlice trimmed: %s\n", strings.TrimSpace(resSlice[i]))
		// log.Printf("Seperator: %s\n", Seperator)
		if resSlice[i] == "" {
			resSlice[i] = "NO OUTPUT"
		}
		if inpSlice[i] == "" {
			inpSlice[i] = "NO INPUT"
		}
		results[inpSlice[i]] = resSlice[i]
	}

	return results
}

// func compareLineByLine(input, exp, res string) map[string]string {
// 	inpSlice := strings.Split(input, "\n")
// 	expSlice := strings.Split(exp, "\n")
// 	resSlice := strings.Split(res, "\n")

// 	results := make(map[string]string)

// 	// TODO: remove for prod!
// 	if strings.HasPrefix(resSlice[0], "godmode ") {
// 		nStr := strings.TrimPrefix(resSlice[0], "godmode ")
// 		n, e := strconv.Atoi(nStr)
// 		if e != nil {
// 			log.Println(e)
// 			panic("Bad dog!")
// 		}
// 		if n > len(inpSlice) {
// 			n = len(inpSlice)
// 		}

// 		i := 0
// 		log.Println(i, n)
// 		for ; i < n; i++ {
// 			results[inpSlice[i]] = "true"
// 		}
// 		for ; i < len(inpSlice)-1; i++ {
// 			results[inpSlice[i]] = "false"
// 		}

// 		return results
// 	}

// 	// TODO deal with partial success but incorrect result couont
// 	/*if len(expSlice) != len(resSlice) {
// 		return results
// 	}*/

// 	for i := 0; i < len(inpSlice)-1; i++ {
// 		//log.Println("compare: ", inpSlice[i], expSlice[i], resSlice[i])
// 		if i > len(expSlice)-1 || i > len(resSlice)-1 {
// 			results[inpSlice[i]] = "false"
// 		} else {
// 			results[inpSlice[i]] = fmt.Sprintf("%v", expSlice[i] == resSlice[i])
// 		}
// 	}

// 	return results
// }
