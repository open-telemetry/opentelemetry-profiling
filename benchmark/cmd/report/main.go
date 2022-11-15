package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"otelprofiling/implementations/collapsed"
	"otelprofiling/implementations/pprof"

	"github.com/jzelinskie/must"
)

type Implementation interface {
	Name() string
	Append(stacktrace []string, value int)
	Serialize(w io.Writer) error
}

func main() {
	filepaths, err := filepath.Glob("./profiles/intermediary/*")

	if err != nil {
		panic(err)
	}

	for _, path := range filepaths {

		implementations := []Implementation{
			collapsed.New(),
			pprof.New(),
		}

		for _, implementation := range implementations {
			file := must.NotError(os.Open(path))
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				str := scanner.Text()
				arr := strings.Split(str, " ")
				stacktrace := strings.Split(arr[0], ";")
				implementation.Append(stacktrace, must.NotError(strconv.Atoi(arr[1])))
			}

			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}

			buf := &bytes.Buffer{}
			implementation.Serialize(buf)

			basename := filepath.Base(path)
			res := buf.Bytes()

			dir := filepath.Join("profiles", "out", implementation.Name())
			os.MkdirAll(dir, 0755)
			ioutil.WriteFile(filepath.Join(dir, filepath.Base(path)), res, 0644)

			fmt.Printf("Implementation: %s\n", implementation.Name())
			fmt.Printf("Filename: %s\n", basename)
			fmt.Printf("Size of output: %d\n", len(res))
			fmt.Print("\n")
		}
	}
}
