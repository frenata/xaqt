package compilebox

var LanguageMap map[string]Language

type Language struct {
	Compiler           string //`json:"compiler"`
	SourceFile         string //`json:"sourceFile"`
	OptionalExecutable string //`json:"optionalExecutable"`
	CompilerFlags      string //`json:"compilerFlags"`
	Boilerplate        string `json:"boilerplate"`
	CommentPrefix      string `json:"commentPrefix"`
	Disabled           string //`json:"disabled"`
}
