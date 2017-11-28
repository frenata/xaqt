package testbox

import (
	"log"
	"strings"
)

// input is n test calls seperated by newlines
// input and expected MUST end in newlines
func Test(language, code, input, expected string) []bool {

	lang := LanguageMap[language]

	sb := NewSandbox(lang, code, input, DefaultSandboxOptions())

	output, err := sb.Run()
	if err != nil {
		log.Println(err)
		return []bool{false}
	}

	splitOutput := strings.SplitN(output, "*-COMPILEBOX::ENDOFOUTPUT-*", 2)
	timeTaken := splitOutput[1]
	_ = timeTaken
	result := splitOutput[0]

	return compareLineByLine(expected, result)
}

func compareLineByLine(exp, res string) []bool {
	expSlice := strings.Split(exp, "\n")
	resSlice := strings.Split(res, "\n")

	results := make([]bool, len(expSlice))

	// TODO deal with partial success but incorrect result couont
	if len(expSlice) != len(resSlice) {
		return results
	}

	for i := range expSlice {
		results[i] = expSlice[i] == resSlice[i]
	}

	return results[:len(results)-1]
}
