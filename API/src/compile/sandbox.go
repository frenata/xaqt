package compile

import (
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
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
	timeout int
}

func DefaultSandboxOptions() SandboxOptions {
	rand := strconv.Itoa(rand.Intn(1000))
	pwd, _ := os.Getwd()
	return SandboxOptions{"temp/" + rand, pwd, "virtual_machine", 20}
}

func NewSandbox(l Language, code, stdin string, options SandboxOptions) *Sandbox {
	box := Sandbox{l, code, stdin, options}

	return &box
}

func (s *Sandbox) Run() {
	s.prepare()
	s.execute()
}

func (s *Sandbox) prepare() {
	tmpFolder, err := ioutil.TempDir("", "docker-test")
	if err != nil {
		log.Fatal(err)
	}

	if err := os.Chmod(tmpFolder, 0777); err != nil {
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

func (s *Sandbox) execute() {
	compiler := s.language.Compiler
	filename := s.language.SourceFile
	optionalExecutable := s.language.OptionalExecutable
	flags := s.language.CompilerFlags

	dockerCommand := s.options.path + "/DockerTimeout.sh"

	args := []string{strconv.Itoa(s.options.timeout) + "s", "-u", "mysql", "-i", "-t", "--volume=" + s.options.folder + ":/usercode", s.options.vm_name, "/usercode/script.sh", compiler, filename, optionalExecutable, flags}

	done := make(chan error)

	go spawnDocker(dockerCommand, args, done)

	select {
	case res := <-done:
		//log.Println(res)
		_ = res
		bytes, err := ioutil.ReadFile(s.options.folder + "/completed")
		log.Println(string(bytes), err)
	case <-time.After(time.Second * 10): // TODO: use timeout
		// clean up spawnDocker
		log.Println("timed out")
	}
}

func spawnDocker(dockerCommand string, args []string, done chan error) {
	cmd := exec.Command(dockerCommand, args...)
	err := cmd.Run()
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
