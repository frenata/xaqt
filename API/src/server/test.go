package main

import "strings"

type id string
type Test struct {
	description string
	io          map[string]string
}

func (t Test) StdIO() (string, string) {
	inputs := make([]string, len(t.io))
	outputs := make([]string, len(t.io))

	i := 0
	for k, v := range t.io {
		inputs[i] = k
		outputs[i] = v
		i++
	}

	return join(inputs), join(outputs)
}

func join(s []string) string {
	return strings.Join(s, "\n") + "\n"
}
