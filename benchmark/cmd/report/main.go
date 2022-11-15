package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"otelprofiling/reference"

	"github.com/jzelinskie/must"
)

func main() {
	filepaths, err := filepath.Glob("./profiles/intermediary/*")

	if err != nil {
		panic(err)
	}

	for _, path := range filepaths {
		file := must.NotError(os.Open(path))

		p := reference.New()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			str := scanner.Text()
			arr := strings.Split(str, " ")
			p.Append(arr[0], must.NotError(strconv.Atoi(arr[1])))
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		buf := &bytes.Buffer{}
		p.Serialize(buf)

		basename := filepath.Base(path)
		size := buf.Bytes()

		fmt.Printf("Filename: %s\n", basename)
		fmt.Printf("Size of output: %d\n", len(size))
		fmt.Print("\n")
	}
}
