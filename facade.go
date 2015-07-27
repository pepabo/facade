package facade

import (
	"bufio"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/mattn/go-colorable"
	"gopkg.in/pipe.v2"
	"io"
	"os"
	"os/exec"
	"sort"
	"strings"
	"sync"
)

func init() {
	log.SetOutput(colorable.NewColorableStdout())
}

type Facade struct {
	Env map[string]string
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

var builtins = map[string]func(){
	"help": help,
	"list": list,
}

func help() {
	fmt.Printf("Usage: $ %s sub-command args...\n", me())
}

func list() {
	prefix := fmt.Sprintf("%s-", me())

	results := make(chan os.FileInfo)
	stop := make(chan struct{})

	var wg sync.WaitGroup
	for _, p := range strings.Split(os.Getenv("PATH"), ":") {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			d, err := os.Open(p)
			if err != nil {
				return
			}
			defer d.Close()

			files, err := d.Readdir(-1)
			if err != nil {
				return
			}

			<-stop

			for _, f := range files {
				if !f.IsDir() && strings.HasPrefix(f.Name(), prefix) {
					results <- f
				}
			}
		}(p)
	}

	go func() {
		wg.Wait()
		close(results)
	}()
	close(stop)

	var commands []string
	for f := range results {
		commands = append(commands, strings.TrimPrefix(f.Name(), prefix))
	}

	if len(commands) > 0 {
		sort.Strings(commands)
		for _, c := range commands {
			println(c)
		}
	} else {
		println("No sub commands found.")
	}
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

	stdin, err := cmd.StdinPipe()
	if err != nil {
		fatal(err.Error())
		os.Exit(1)
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

	go func() {
		if err = cmd.Wait(); err != nil {
			fatal(err.Error())
			os.Exit(1)
		} else {
			os.Exit(0)
		}
	}()

	for {
		p := pipe.Line(
			pipe.Read(os.Stdin),
			pipe.Write(stdin),
		)
		if err := pipe.Run(p); err != nil {
			fatal(err.Error())
			os.Exit(1)
		}
		if s, err := pipe.Output(p); err != nil {
			fatal(string(s))
			fatal(err.Error())
			os.Exit(1)
		}
	}
}

func me() string {
	chunks := strings.Split(os.Args[0], string(os.PathSeparator))
	return chunks[len(chunks)-1]
}

func readFrom(in io.ReadCloser, logger func(string)) {
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		logger(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		in.Close()
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
