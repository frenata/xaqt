package main

import "strings"
import "testbox"

type id = string
type Challenge struct {
	ID          string            `json:"id"`
	Description string            `json:"description"`
	IO          map[string]string `json:"io"`
	SampleIO    string            `json:"sampleIO"`
}

func (t Challenge) StdIO() (string, string) {
	inputs := make([]string, len(t.IO))
	outputs := make([]string, len(t.IO))

	i := 0
	for k, v := range t.IO {
		inputs[i] = k
		outputs[i] = v
		i++
	}

	return joinAndAppend(inputs, testbox.Seperator), joinAndAppend(outputs, testbox.Seperator)
}

func joinAndAppend(sl []string, endChar string) string {
	return strings.Join(sl, endChar) + endChar
}
