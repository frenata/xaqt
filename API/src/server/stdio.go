package main

import "strings"
import "testbox"

type id = string
type Test struct {
	Description string            `json:"description"`
	Io          map[string]string `json:"io"`
	SampleIO    string            `json:"sampleIO"`
}

func (t Test) StdIO() (string, string) {
	inputs := make([]string, len(t.Io))
	outputs := make([]string, len(t.Io))

	i := 0
	for k, v := range t.Io {
		inputs[i] = k
		outputs[i] = v
		i++
	}

	return joinAndAppend(inputs, testbox.Seperator), joinAndAppend(outputs, testbox.Seperator)
}

func joinAndAppend(sl []string, endChar string) string {
	return strings.Join(sl, endChar) + endChar
}
