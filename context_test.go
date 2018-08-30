package xaqt_test

import (
	"testing"
	"time"

	"github.com/frenata/xaqt"
)

var box *xaqt.Context
var tests map[string]string

// Test that each compiler, given the appropriate code, can print "Hello"
func TestCompilers(t *testing.T) {
	langResults := make(map[string]bool)
	for lang, code := range tests {
		langResults[lang] = printsHello(t, lang, code)
	}

	pass := true
	for _, result := range langResults {
		pass = pass && result
	}

	if !pass {
		t.Fatal("Not all languages printed 'Hello' correctly.")
	}
}

func printsHello(t *testing.T, lang, code string) bool {
	stdin := ""
	expected := "Hello"
	stdouts, msg := box.Evaluate(lang, code, []string{stdin})
	//log.Println(stdouts[0], msg)

	// check for timeout
	if msg.Data == "Timed out" { // TODO make proper error type
		t.Logf("%s timed out during 'Hello' test.", lang)
		return false
	}

	if stdouts[0] != expected {
		t.Log(stdouts)
		t.Logf("%s failed 'Hello' test.", lang)
		return false
	}

	t.Logf("%s passed 'Hello' test.", lang)
	return true
}

func init() {
	box, _ = xaqt.NewContext(
		xaqt.GetCompilers(),
		xaqt.Timeout(time.Second*10))

	tests = map[string]string{

		// currently passing:
		"C++":        "#include <iostream>\nusing namespace std;\n\nint main() {\n\tcout<<\"Hello\";\n\treturn 0;\n}",
		"Java":       "\n\nimport java.io.*;\n\nclass myCode\n{\n\tpublic static void main (String[] args) throws java.lang.Exception\n\t{\n\t\t\n\t\tSystem.out.println(\"Hello\");\n\t}\n}",
		"C#":         "using System;\n\npublic class Challenge\n{\n\tpublic static void Main()\n\t{\n\t\t\tConsole.WriteLine(\"Hello\");\n\t}\n}",
		"Clojure":    "(println \"Hello\")",
		"Perl":       "use strict;\nuse warnings\n;use v5.14; say 'Hello';",
		"Golang":     "package main\nimport \"fmt\"\n\nfunc main(){\n  \n\tfmt.Printf(\"Hello\")\n}",
		"JavaScript": "console.log(\"Hello\");",
		"Python":     "print(\"Hello\")",
		"Ruby":       "puts \"Hello\"",
		"Bash":       "echo 'Hello' ",
		"PHP":        "<?php\n$ho = fopen('php://stdout', \"w\");\n\nfwrite($ho, \"Hello\");\n\n\nfclose($ho);\n",
		"Haskell":    "module Main where\nmain = putStrLn \"Hello\"",

		// currently failing
		// Scala: don't understand the error this generates
		// "Scala" : "object HelloWorld {def main(args: Array[String]) = println(\"Hello\")}",

		// Rust seems to be missing and there's a problem setting environment variables
		// "Rust" : "fn main() {\n\tprintln!(\"Hello\");\n}",

		//"MySQL":"create table myTable(name varchar(10));\ninsert into myTable values(\"Hello\");\nselect * from myTable;",

		//"Objective-C": "#include <Foundation/Foundation.h>\n\n@interface Challenge\n+ (const char *) classStringValue;\n@end\n\n@implementation Challenge\n+ (const char *) classStringValue;\n{\n    return \"Hey!\";\n}\n@end\n\nint main(void)\n{\n    printf(\"%s\\n\", [Challenge classStringValue]);\n    return 0;\n}",

		//"VB.NET": "Imports System\n\nPublic Class Challenge\n\tPublic Shared Sub Main() \n    \tSystem.Console.WriteLine(\"Hello\")\n\tEnd Sub\nEnd Class",
	}
}
