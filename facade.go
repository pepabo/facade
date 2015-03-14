package facade

import (
	"bytes"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/mattn/go-colorable"
	"os"
	"os/exec"
	"strings"
)

func init() {
	log.SetOutput(colorable.NewColorableStdout())
}

func Run() {
	chunks := strings.Split(os.Args[0], string(os.PathSeparator))
	me := chunks[len(chunks)-1]
	sub := os.Args[1]
	full := fmt.Sprintf("%s-%s", me, sub)

	var buf bytes.Buffer
	cmd := exec.Command(full, os.Args[2:]...)
	cmd.Stdout = &buf

	err := cmd.Run()
	if err != nil {
		fatal(err.Error())
		return
	}

	info(buf.String())
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
