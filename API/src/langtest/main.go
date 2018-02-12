package main

import (
	"fmt"
	"log"
	"sandbox"
)

func main() {
	box := sandbox.New("data/compilers.json")
	compilerTests := make(map[string]string)

	// currently passing:

	compilerTests["C++"] = "#include <iostream>\nusing namespace std;\n\nint main() {\n\tcout<<\"Hello\";\n\treturn 0;\n}"
	compilerTests["Java"] = "\n\nimport java.io.*;\n\nclass myCode\n{\n\tpublic static void main (String[] args) throws java.lang.Exception\n\t{\n\t\t\n\t\tSystem.out.println(\"Hello\");\n\t}\n}"
	compilerTests["C#"] = "using System;\n\npublic class Challenge\n{\n\tpublic static void Main()\n\t{\n\t\t\tConsole.WriteLine(\"Hello\");\n\t}\n}"
	compilerTests["Clojure"] = "(println \"Hello\")"
	compilerTests["Perl"] = "use strict;\nuse warnings\n;use v5.14; say 'Hello';"
	compilerTests["Golang"] = "package main\nimport \"fmt\"\n\nfunc main(){\n  \n\tfmt.Printf(\"Hello\")\n}"
	compilerTests["JavaScript"] = "console.log(\"Hello\");"
	compilerTests["Python"] = "print(\"Hello\")"
	compilerTests["Ruby"] = "puts \"Hello\""
	compilerTests["Bash"] = "echo 'Hello' "
	compilerTests["PHP"] = "<?php\n$ho = fopen('php://stdout', \"w\");\n\nfwrite($ho, \"Hello\");\n\n\nfclose($ho);\n"

	// currently broken:

	// Haskell ghc missing, maybe need to rebuild docker file
	// compilerTests["Haskell"] = "module Main where\nmain = putStrLn \"Hello\""

	// Scala: don't understand the error this generates
	// compilerTests["Scala"] = "object HelloWorld {def main(args: Array[String]) = println(\"Hello\")}"

	// Rust seems to be missing and there's a problem setting environment variables
	// compilerTests["Rust"] = "fn main() {\n\tprintln!(\"Hello\");\n}"

	/*
		"MySQL":"create table myTable(name varchar(10));\ninsert into myTable values(\"Hello\");\nselect * from myTable;",
		"Objective-C": "#include <Foundation/Foundation.h>\n\n@interface Challenge\n+ (const char *) classStringValue;\n@end\n\n@implementation Challenge\n+ (const char *) classStringValue;\n{\n    return \"Hey!\";\n}\n@end\n\nint main(void)\n{\n    printf(\"%s\\n\", [Challenge classStringValue]);\n    return 0;\n}",
		"VB.NET": "Imports System\n\nPublic Class Challenge\n\tPublic Shared Sub Main() \n    \tSystem.Console.WriteLine(\"Hello\")\n\tEnd Sub\nEnd Class",
	*/

	stdin := ""
	expected := "Hello"
	langResults := make(map[string]string)
	for lang, code := range compilerTests {
		stdouts, msg := box.EvalWithStdins(lang, code, []string{stdin})
		// oOut, oMsg := box.CompileAndPrint(lang, code, "test")
		log.Println(stdouts[0], msg)
		// log.Println(oOut, oMsg)
		if stdouts[0] == expected {
			log.Printf("%s passed 'Hello' test.", lang)
			langResults[lang] = "Pass"
		} else {
			log.Println(stdouts)
			log.Printf("%s failed 'Hello' test.", lang)
			langResults[lang] = "Fail"
		}

		fmt.Println("-----------------------------------------------------")
	}

	for lang, result := range langResults {
		fmt.Printf("%s -> %s\n", lang, result)
	}
}
