package xaqt

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/stretchr/testify/suite"
)

// test suite for unit testing the sandbox
type SandboxTestSuite struct {
	suite.Suite
}

func (s *SandboxTestSuite) TestPrepareTmpDir_WithSourceString() {
	var (
		err error
	)

	expectedSourceString := "f(g) => g()"
	expectedStdin := "h() => search('h2o')"
	srcFileName := "source.code"
	// setup test sandbox
	sandbox := &sandbox{
		ID:       "sahara",
		language: ExecutionDetails{SourceFile: srcFileName},
		code: Code{
			IsFile: false,
			String: expectedSourceString,
		},
		stdin: expectedStdin,
		options: options{
			execDir: "",
			path:    DataPath(),
		},
	}

	err = sandbox.PrepareTmpDir()
	s.NoError(err)

	// defer cleanup
	defer os.RemoveAll(sandbox.options.execDir)

	// input file exists?
	inputFile := filepath.Join(sandbox.options.execDir, STDIN_FILENAME)
	_, err = os.Stat(inputFile)
	s.NoError(err)

	// input file correct?
	actualStdinBytes, err := ioutil.ReadFile(inputFile)
	s.NoError(err)
	s.EqualValues(expectedStdin, string(actualStdinBytes))

	// source file exists?
	sourceFile := filepath.Join(sandbox.options.execDir, srcFileName)
	_, err = os.Stat(sourceFile)
	s.NoError(err)

	// source file correct?
	actualSourceBytes, err := ioutil.ReadFile(sourceFile)
	s.NoError(err)
	s.EqualValues(expectedSourceString, string(actualSourceBytes))
}

func (s *SandboxTestSuite) TestPrepareTmpDir_WithSourceFile() {
	var (
		err error
	)

	expectedSourceString := "f(g) => g()"
	expectedStdin := "h() => search('cauliflower')"
	srcFileName := "source.code"

	// create test source file
	err = ioutil.WriteFile(srcFileName, []byte(expectedSourceString), 0777)
	s.NoError(err)
	defer os.Remove(srcFileName)

	// setup test sandbox
	sandbox := &sandbox{
		ID:       "gobi",
		language: ExecutionDetails{SourceFile: srcFileName},
		code: Code{
			IsFile:         true,
			SourceFileName: srcFileName,
			Path:           "",
		},
		stdin: expectedStdin,
		options: options{
			execDir: "",
			path:    DataPath(),
		},
	}

	err = sandbox.PrepareTmpDir()
	s.NoError(err)

	// defer cleanup
	defer os.RemoveAll(sandbox.options.execDir)

	// input file exists?
	inputFile := filepath.Join(sandbox.options.execDir, STDIN_FILENAME)
	_, err = os.Stat(inputFile)
	s.NoError(err)

	// input file correct?
	actualStdinBytes, err := ioutil.ReadFile(inputFile)
	s.NoError(err)
	s.EqualValues(expectedStdin, string(actualStdinBytes))

	// source file exists?
	sourceFile := filepath.Join(sandbox.options.execDir, srcFileName)
	_, err = os.Stat(sourceFile)
	s.NoError(err)

	// source file correct?
	actualSourceBytes, err := ioutil.ReadFile(sourceFile)
	s.NoError(err)
	s.EqualValues(expectedSourceString, string(actualSourceBytes))
}

func (s *SandboxTestSuite) TestPrepareTmpDir_WithResourceFiles() {
	var (
		err error
	)

	expectedSourceString := "f(g) => g()"
	expectedStdin := "h() => search('squid?')"
	srcFileName := "source.code"
	resourceFiles := map[string]string{
		"resource_one.csv": "one,golden,mole",
		"resource_two.csv": "two,pickled,newts",
	}

	// create test source file
	err = ioutil.WriteFile(srcFileName, []byte(expectedSourceString), 0777)
	s.NoError(err)
	defer os.Remove(srcFileName)

	// create test resource files
	for resourceFileName, resourceString := range resourceFiles {
		err = ioutil.WriteFile(resourceFileName, []byte(resourceString), 0777)
		s.NoError(err)
		defer os.Remove(resourceFileName)
	}

	// setup test sandbox
	sandbox := &sandbox{
		ID:       "kalahari",
		language: ExecutionDetails{SourceFile: srcFileName},
		code: Code{
			IsFile:            true,
			SourceFileName:    srcFileName,
			ResourceFileNames: []string{"resource_one.csv", "resource_two.csv"},
			Path:              "",
		},
		stdin: expectedStdin,
		options: options{
			execDir: "",
			path:    DataPath(),
		},
	}

	err = sandbox.PrepareTmpDir()
	s.NoError(err)

	// defer cleanup
	defer os.RemoveAll(sandbox.options.execDir)

	// input file exists?
	inputFile := filepath.Join(sandbox.options.execDir, STDIN_FILENAME)
	_, err = os.Stat(inputFile)
	s.NoError(err)

	// input file correct?
	actualStdinBytes, err := ioutil.ReadFile(inputFile)
	s.NoError(err)
	s.EqualValues(expectedStdin, string(actualStdinBytes))

	// source file exists?
	sourceFile := filepath.Join(sandbox.options.execDir, srcFileName)
	_, err = os.Stat(sourceFile)
	s.NoError(err)

	// source file correct?
	actualSourceBytes, err := ioutil.ReadFile(sourceFile)
	s.NoError(err)
	s.EqualValues(expectedSourceString, string(actualSourceBytes))

	// are resource files existing and correct?
	for resourceFileName, expectedResourceString := range resourceFiles {
		resourceFilePath := filepath.Join(sandbox.options.execDir, resourceFileName)
		_, err = os.Stat(resourceFilePath)
		s.NoError(err)

		actualResourceBytes, err := ioutil.ReadFile(resourceFilePath)
		s.NoError(err)
		s.EqualValues(expectedResourceString, string(actualResourceBytes))

	}
}

func (s *SandboxTestSuite) TestRewriteUserFiles() {
	var (
		err error
	)

	originalSourceString := "do(something smart)"
	expectedSourceString := "f(g) => g()"
	srcFileName := "source.code"
	originalResourceFiles := map[string]string{
		"resource_one.csv": "one,golden,mole",
		"resource_two.csv": "two,pickled,newts",
	}
	expectedResourceFiles := map[string]string{
		"resource_one.csv": "five,purring,lions",
		"resource_two.csv": "eighty,sauted,sonnets",
	}

	// create test source file
	err = ioutil.WriteFile(srcFileName, []byte(originalSourceString), 0777)
	s.NoError(err)
	defer os.Remove(srcFileName)

	// create test resource files
	for resourceFileName, resourceString := range originalResourceFiles {
		err = ioutil.WriteFile(resourceFileName, []byte(resourceString), 0777)
		s.NoError(err)
		defer os.Remove(resourceFileName)
	}

	// setup test sandbox
	sandbox := &sandbox{
		ID:       "kalahari",
		language: ExecutionDetails{SourceFile: srcFileName},
		code: Code{
			IsFile:            true,
			SourceFileName:    srcFileName,
			ResourceFileNames: []string{"resource_one.csv", "resource_two.csv"},
			Path:              "",
		},
		options: options{
			execDir: "",
			path:    DataPath(),
		},
	}

	err = sandbox.PrepareTmpDir()
	s.NoError(err)

	// defer cleanup
	defer os.RemoveAll(sandbox.options.execDir)

	// simulate execution by rewriting the files in the tmp directory
	// rewrite test source file
	err = ioutil.WriteFile(
		filepath.Join(sandbox.options.execDir, srcFileName),
		[]byte(expectedSourceString),
		0777,
	)
	s.NoError(err)
	// rewrite test resource files
	for resourceFileName, resourceString := range expectedResourceFiles {
		err = ioutil.WriteFile(
			filepath.Join(sandbox.options.execDir, resourceFileName),
			[]byte(resourceString),
			0777,
		)
		s.NoError(err)
	}

	err = sandbox.rewriteUserFiles()
	s.NoError(err)

	// source file exists?
	_, err = os.Stat(srcFileName)
	s.NoError(err)

	// source file correctly modified?
	actualSourceBytes, err := ioutil.ReadFile(srcFileName)
	s.NoError(err)
	s.EqualValues(expectedSourceString, string(actualSourceBytes))

	// are resource files existing and correctly modified?
	for resourceFileName, expectedResourceString := range expectedResourceFiles {
		_, err = os.Stat(resourceFileName)
		s.NoError(err)

		actualResourceBytes, err := ioutil.ReadFile(resourceFileName)
		s.NoError(err)
		s.EqualValues(expectedResourceString, string(actualResourceBytes))

	}
}

// somehow we want to test if the process being run is constrained s.t. it
// can't write to other parts of the file system or something. but this should be
// the whole point of using docker so maybe not...
