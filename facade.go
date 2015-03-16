package facade

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/mattn/go-colorable"
	"io"
	"os"
	"os/exec"
	"strings"
)

func init() {
	log.SetOutput(colorable.NewColorableStdout())
}

type Facade struct {
	Env map[string]string
}

var builtins = map[string]func(){
	"help": help,
}

func help() {
	fmt.Printf("Usage: $ %s sub-command args...\n", me())
}

func (f *Facade) Run() {
	var subCommand string
	if len(os.Args) > 1 {
		subCommand = os.Args[1]
	}

	if subCommand == "" {
		help()
	} else {
		if b := builtins[subCommand]; b != nil {
			b()
		} else {
			f.dispatch(subCommand)
		}
	}

	os.Exit(0)
}

func Run() {
	f := &Facade{}
	f.Run()
}

func (f *Facade) dispatch(subCommand string) {
	cmd := exec.Command(fmt.Sprintf("%s-%s", me(), subCommand), os.Args[2:]...)
	if f.Env != nil {
		newenv := os.Environ()
		for k, v := range f.Env {
			newenv = append(newenv, fmt.Sprintf("%s=%s", k, v))
		}
		cmd.Env = newenv
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fatal(err.Error())
		os.Exit(1)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		fatal(err.Error())
		os.Exit(1)
	}
	if err = cmd.Start(); err != nil {
		fatal(err.Error())
		os.Exit(1)
	}

	go readFrom(stdout, info)
	go readFrom(stderr, fatal)

	if err = cmd.Wait(); err != nil {
		fatal(err.Error())
		os.Exit(1)
	}
}

func me() string {
	chunks := strings.Split(os.Args[0], string(os.PathSeparator))
	return chunks[len(chunks)-1]
}

func readFrom(in io.ReadCloser, logger func(string)) {
	b := make([]byte, 1024)
	for {
		n, err := in.Read(b)
		if n > 0 {
			if err != nil {
				if err == io.EOF {
					in.Close()
				}
			}
			logger(string(b[:n]))
		} else {
			in.Close()
		}
	}
}

func info(s string) {
	for _, line := range stringToLines(s) {
		log.Info(line)
	}
}

func fatal(s string) {
	for _, line := range stringToLines(s) {
		log.Error(line)
	}
}

func stringToLines(s string) []string {
	return strings.Split(strings.TrimRight(s, "\n"), "\n")
}
