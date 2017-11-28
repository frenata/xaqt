package testbox

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
	//rand := strconv.Itoa(rand.Intn(1000))
	pwd, _ := os.Getwd()

	tmp := ""
	if runtime.GOOS == "darwin" {
		tmp = "/tmp"
	}

	return SandboxOptions{tmp, pwd, "virtual_machine", time.Second * 20}
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

	if err := os.Chmod(tmpFolder, 0755); err != nil {
		log.Fatal(err)
	}

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
	err = ioutil.WriteFile(tmpFolder+"/inputFile", []byte(s.stdin), 0777)
	if err != nil {
		log.Fatal(err)
	}
	// log msg
}

func (s *Sandbox) execute() (string, error) {
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
		log.Printf("Docker returns: %v", res)
		errorBytes, err := ioutil.ReadFile(s.options.folder + "/errors")
		bytes, err := ioutil.ReadFile(s.options.folder + "/completed")
		// TODO: handle file io errors
		if len(errorBytes) > 0 {
			bytes, err = errorBytes, fmt.Errorf("compile error")
		}

		return string(bytes), err
	case <-time.After(time.Second * s.options.timeout): // TODO: use timeout
		// TODO clean up temp folders spawnDocker
		log.Println("timed out")
		return "", fmt.Errorf("Timed out")
	}
}

func spawnDocker(dockerCommand string, args []string, done chan error) {
	cmd := exec.Command(dockerCommand, args...)
	bytes, err := cmd.CombinedOutput()
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
