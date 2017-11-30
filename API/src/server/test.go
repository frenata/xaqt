package main

import "strings"

type id string
type Test struct {
	Description string            `json:"description"`
	Io          map[string]string `json:"io"`
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

	return join(inputs), join(outputs)
}

func join(s []string) string {
	return strings.Join(s, "\n") + "\n"
}
