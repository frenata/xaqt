package main

import (
	"fmt"
	"log"
	"sandbox"
)

func main() {
	box := sandbox.New("data/compilers.json")
	compilerTests := make(map[string]string)

	// need test for Haskell
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
	// FIXME
	// compilerTests["Scala"] = "object HelloWorld {def main(args: Array[String]) = println(\"Hello\")}"
	// FIXME compiler
	// compilerTests["Rust"] = "fn main() {\n\tprintln!(\"Hello\");\n}"
	// FIXME
	// compilerTests["PHP"] = "<?php\n$ho = fopen('php://stdout', \"w\");\n\nfwrite($ho, \"Hello\");\n\n\nfclose($ho);\n"

	/*
		"MySQL":"create table myTable(name varchar(10));\ninsert into myTable values(\"Hello\");\nselect * from myTable;",
		"Objective-C": "#include <Foundation/Foundation.h>\n\n@interface Challenge\n+ (const char *) classStringValue;\n@end\n\n@implementation Challenge\n+ (const char *) classStringValue;\n{\n    return \"Hey!\";\n}\n@end\n\nint main(void)\n{\n    printf(\"%s\\n\", [Challenge classStringValue]);\n    return 0;\n}",
		"VB.NET": "Imports System\n\nPublic Class Challenge\n\tPublic Shared Sub Main() \n    \tSystem.Console.WriteLine(\"Hello\")\n\tEnd Sub\nEnd Class",
	*/

	stdin := "" + sandbox.Seperator
	expected := "Hello" + sandbox.Seperator
	langResults := make(map[string]string)
	for lang, code := range compilerTests {
		out, msg := box.CompileAndChallenge(lang, code, stdin, expected)
		// oOut, oMsg := box.CompileAndPrint(lang, code, "test")
		log.Println(out, msg)
		// log.Println(oOut, oMsg)
		langResults[lang] = out[""]
		if out[""] == "Pass" {
			log.Printf("%s passed 'Hello' test.", lang)
		} else {
			log.Println(out)
			log.Printf("%s failed 'Hello' test.", lang)
		}
		fmt.Println("-----------------------------------------------------")
	}

	for lang, result := range langResults {
		fmt.Printf("%s -> %s", lang, result)
	}
}
