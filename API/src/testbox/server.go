package testbox

import (
	"log"
	"strings"
)

// input is n test calls seperated by newlines
// input and expected MUST end in newlines
func Test(language, code, input, expected string) map[string]bool {

	lang := LanguageMap[language]

	sb := NewSandbox(lang, code, input, DefaultSandboxOptions())

	output, err := sb.Run()
	if err != nil {
		log.Println(err)
		return nil
	}

	splitOutput := strings.SplitN(output, "*-COMPILEBOX::ENDOFOUTPUT-*", 2)
	timeTaken := splitOutput[1]
	_ = timeTaken
	result := splitOutput[0]

	return compareLineByLine(input, expected, result)
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
