package main

import (
	"log"
	"strings"
	"testbox"
)

type id = string
type Challenge struct {
	ID          string            `json:"id"`
	Description string            `json:"description"`
	IO          map[string]string `json:"io"`
	SampleIO    string            `json:"sampleIO"`
}

func (c *Challenge) getCases() (string, string) {
	inputs := make([]string, len(c.IO))
	outputs := make([]string, len(c.IO))

	log.Printf("getCases, challenge: %v\n", c)
	i := 0
	for k, v := range c.IO {
		inputs[i] = k
		outputs[i] = v
		i++
	}

	return joinAndAppend(inputs, testbox.Seperator), joinAndAppend(outputs, testbox.Seperator)
}

func joinAndAppend(sl []string, endChar string) string {
	return strings.Join(sl, endChar) + endChar
}
