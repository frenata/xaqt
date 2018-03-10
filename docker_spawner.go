package compilebox

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

type Sandbox struct {
	language Language
	code     string
	stdin    string
	options  SandboxOptions
}

type SandboxOptions struct {
	folder  string
	path    string
	vm_name string
	timeout time.Duration
}

func DefaultSandboxOptions() SandboxOptions {
	pwd, _ := os.Getwd()

	tmp := ""
	if runtime.GOOS == "darwin" {
		tmp = "/tmp"
	}

	return SandboxOptions{tmp, pwd, "virtual_machine", time.Second * 5}
}

func NewSandbox(l Language, code, stdin string, options SandboxOptions) *Sandbox {
	box := Sandbox{l, code, stdin, options}

	return &box
}

func (s *Sandbox) Run() (string, error) {
	s.prepare()
	return s.execute()
}

func (s *Sandbox) prepare() {
	tmpFolder, err := ioutil.TempDir(s.options.folder, "docker-test")
	if err != nil {
		log.Fatal(err)
	}

	if err := os.Chmod(tmpFolder, 0777); err != nil {
		log.Fatal(err)
	}

	// record tmpdir for easy deletion
	s.options.folder = tmpFolder

	err = s.copyPayload()
	if err != nil {
		log.Fatal(err)
	}

	//log.Println(s.code)
	err = ioutil.WriteFile(tmpFolder+"/"+s.language.SourceFile, []byte(s.code), 0777)
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println(ioutil.ReadFile(tmpFolder + "/" + s.language.SourceFile))

	// write a file for stdin
	log.Printf("writing inputfile: %v", s.stdin)
	err = ioutil.WriteFile(tmpFolder+"/inputFile", []byte(s.stdin), 0777)
	if err != nil {
		log.Fatal(err)
	}
	// log msg
}

func (s *Sandbox) execute() (string, error) {
	defer os.RemoveAll(s.options.folder)

	compiler := s.language.Compiler
	filename := s.language.SourceFile
	optionalExecutable := s.language.OptionalExecutable
	flags := s.language.CompilerFlags

	dockerCommand := s.options.path + "/DockerTimeout.sh"

	args := []string{fmt.Sprintf("%s", s.options.timeout), "-u", "mysql", "-i", "-t", "--volume=" + s.options.folder + ":/usercode", s.options.vm_name, "/usercode/script.sh", compiler, filename, optionalExecutable, flags}

	done := make(chan error)

	go spawnDocker(dockerCommand, args, done)

	select {
	case res := <-done:
		_ = res

		log.Printf("Docker returns: %v", res)
		errorBytes, err := ioutil.ReadFile(s.options.folder + "/errors")
		if err != nil {
			log.Println("Missing error file")
		}

		outputBytes, err := ioutil.ReadFile(s.options.folder + "/completed")
		if err != nil {
			log.Println("Missing completed file")
		}

		log.Printf("Completed File: \n%s", string(outputBytes))
		// TODO: handle file io errors

		if len(errorBytes) > 0 {
			err = fmt.Errorf("compile error: %s", errorBytes)
		}

		return string(outputBytes), err
	case <-time.After(s.options.timeout):
		log.Println("timed out")
		return "", fmt.Errorf("Timed out")
	}
}

func spawnDocker(dockerCommand string, args []string, done chan error) {
	cmd := exec.Command(dockerCommand, args...)
	bytes, err := cmd.CombinedOutput()
	_ = bytes
	log.Printf("Docker stdout: %v", string(bytes))
	done <- err
}

func (s Sandbox) copyPayload() error {
	source := filepath.Join(s.options.path, "Payload")
	dest := filepath.Join(s.options.folder)

	directory, err := os.Open(source)
	if err != nil {
		return err
	}

	files, err := directory.Readdir(-1)
	if err != nil {
		return err
	}

	for _, file := range files {
		// read the file
		destfile := dest + "/" + file.Name()
		sourcefile := source + "/" + file.Name()
		bytes, err := ioutil.ReadFile(sourcefile)
		if err != nil {
			return err
		}

		// write the file to tmp
		err = ioutil.WriteFile(destfile, bytes, 0777)
		if err != nil {
			return err
		}
	}

	return nil
}
