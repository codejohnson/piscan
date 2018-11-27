package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func showFileSegment(inFileName string, from int64, size int64) error {
	f, err := os.Open(inFileName)
	if err != nil {
		return err
	}
	defer f.Close()
	f.Seek(from, 0)
	bufferedReader := bufio.NewReader(f)
	byteSlice := make([]byte, size)
	_, err = bufferedReader.Read(byteSlice)
	fmt.Printf("%q", byteSlice)
	return nil
}

func resetTerminarColors() {
	esc := "\u001b"
	reset := "[0m"
	print(esc + reset)
}

func getParamValue(arg string) (value string, present bool) {
	for _, commadArgument := range os.Args {
		if strings.HasPrefix(commadArgument, arg) {
			values := strings.Split(commadArgument, ":")
			if len(values) == 2 {
				value = values[1]
				present = true
				return
			}
		}
	}
	return
}

func getCommandLineArguments() (inputFileName string, from string, size string, err error) {
	present := false
	if inputFileName, present = getParamValue("-i"); !present {
		err = fmt.Errorf("error: data file is required")
		return
	}
	if from, present = getParamValue("-from"); !present {
		from = "0"
	}
	if size, present = getParamValue("-size"); !present {
		size = "255"
	}
	return
}

func main() {
	resetTerminarColors()
	inputFileName, from, size, err := getCommandLineArguments()
	if err != nil {
		println(err.Error())
		return
	}
	ifrom := 0
	if ifrom, err = strconv.Atoi(from); err != nil {
		println(err.Error())
		return
	}

	isize := 0
	if isize, err = strconv.Atoi(size); err != nil {
		println(err.Error())
		return
	}
	println(from, size)
	if err = showFileSegment(inputFileName, int64(ifrom), int64(isize)); err != nil {
		println(err.Error())
		return
	}
	return
}
