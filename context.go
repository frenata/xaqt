package xaqt

import (
	"fmt"
	"log"
	"strings"
	"time"
)

// Context creates a execution context for evaluating user code.
type Context struct {
	compilers Compilers
	options
}

// Message represents details on success or failure of execution.
type Message struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

// options to control how the sandbox is executed.
type options struct {
	folder  string // path to folder where results should be recorded
	path    string // path to execution script
	image   string // name of docker image to run
	timeout time.Duration
}

// Uses default sandbox options.
func newDefault(compilers Compilers) *Context {
	c := &Context{compilers, options{}}
	_ = defaultOptions(c)
	return c
}

// NewContext creates a context from a map of compilers and
// some user provided options.
func NewContext(compilers Compilers, options ...option) (*Context, error) {
	context := newDefault(compilers)

	for _, option := range options {
		err := option(context)
		if err != nil {
			return nil, err
		}
	}

	return context, nil
}

// Evaluate code in a given language and for a set of 'stdin's.
func (c *Context) Evaluate(language, code string, stdins []string) ([]string, Message) {
	stdinGlob := glob(stdins)
	results, msg := c.run(language, code, stdinGlob)

	return unglob(results), msg
}

// input is n test calls seperated by newlines
// input and expected MUST end in newlines
func (c *Context) run(language, code, stdinGlob string) (string, Message) {
	log.Printf("launching new %s sandbox", language)
	// log.Printf("launching sandbox...\nLanguage: %s\nStdin: %sCode: Hidden\n", language, stdinGlob)

	lang, ok := c.compilers[strings.ToLower(language)]
	if !ok || lang.Disabled == "true" {
		return "", Message{"error", "language not supported"}
	}

	if code == "" {
		return "", Message{"error", "no code submitted"}
	}

	sb, err := newSandbox(lang.ExecutionDetails, code, stdinGlob, c.options)
	if err != nil {
		log.Printf("sandbox initialization error: %v", err)
		return "", Message{"error", fmt.Sprintf("%s", err)}
	}

	// run the new sandbox
	output, err := sb.run()
	if err != nil {
		log.Printf("sandbox run error: %v", err)
		return output, Message{"error", fmt.Sprintf("%s", err)}
	}

	splitOutput := strings.SplitN(output, "*-COMPILEBOX::ENDOFOUTPUT-*", 2)
	timeTaken := splitOutput[1]
	result := splitOutput[0]

	return result, Message{"success", "compilation took " + timeTaken + " seconds"}
}

// Languages returns a list of available language names.
func (c *Context) Languages() map[string]CompositionDetails { return c.compilers.availableLanguages() }
