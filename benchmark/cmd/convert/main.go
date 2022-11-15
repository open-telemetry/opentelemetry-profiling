package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"otelprofiling/anonymizer"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jzelinskie/must"
)

func main() {
	filepaths, err := filepath.Glob("./profiles/src/*")

	if err != nil {
		panic(err)
	}

	for _, path := range filepaths {
		if strings.HasSuffix(path, ".collapsed") {
			inputFile := must.NotError(os.Open(path))
			outPath := strings.Replace(path, "src", "intermediary", 1)
			outputFile := must.NotError(os.Create(outPath))
			scanner := bufio.NewScanner(inputFile)
			a := anonymizer.New()
			for scanner.Scan() {
				str := scanner.Text()
				arr := strings.Split(str, " ")

				stacktrace := arr[0]
				array := strings.Split(stacktrace, ";")
				for i := 0; i < len(array); i++ {
					array[i] = a.Anonymize(array[i])
				}
				stacktrace = strings.Join(array, ";")
				line := fmt.Sprintf("%s %v", stacktrace, must.NotError(strconv.Atoi(arr[1])))
				outputFile.Write([]byte(line))
			}

			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}
		} else {
			panic("Unexpected file extension")
		}
	}
}
