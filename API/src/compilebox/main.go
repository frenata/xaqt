package main

import (
	"compile"
)

func main() {
	lang := compile.LanguageMap["Nodejs"]

	code := "require('express'); console.log('hello')"

	sb := compile.NewSandbox(lang, code, "", compile.DefaultSandboxOptions())

	sb.Run()
}
