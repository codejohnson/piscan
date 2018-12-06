package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

const gigabyte int = 1024 * 1024 * 1024

type repetitions struct {
	summary        [10]uint64
	bytesProcessed int64
	inFileName     string
	outName        string
	verbose        bool
	startOn        int64
	curDiskPtrRef  int64
	minRepetitions int
	bufferSize     int
}

func persistInit(filename string) {
	var _, err = os.Stat(filename)
	// create file if not exists
	if !os.IsNotExist(err) {
		var err = os.Remove(filename)
		if err != nil {
			panic(err)
		}
	}
	_, err = os.OpenFile(filename, os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
}

func moveToFilePosition(f *os.File, restartPosition int64) {
	if restartPosition == 0 {
		println("scanning from the beginning...")
	} else {
		println("scanning from position $restartPosition...")
		f.Seek(restartPosition, 0)
		println("restarted!")
	}
}

func (r *repetitions) displayRepetition(digit byte, repetitions int, buffer []byte, bufferPosition int) {
	fmt.Printf("\n%c (%d) :...", digit, repetitions)
	from := bufferPosition - 10
	to := bufferPosition + repetitions + 10
	for i := from; i <= len(buffer) && i <= to; i++ {
		if i == bufferPosition || i == to-10 {
			print(" ")
		}
		fmt.Printf("%c", buffer[i])
	}
	print("...")
}

func (r *repetitions) saveRepetition(digit byte, repetitions int, bufferPosition int) {
	f, err := os.OpenFile(r.outName, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.WriteString(fmt.Sprintf("%c;%d;%d\n", digit, repetitions, r.curDiskPtrRef+(int64)(bufferPosition)))
}

func (r *repetitions) countRepetitions(buffer []byte, numBytes int) int {
	var i int
	for i = 0; i < numBytes; i++ {
		j := i + 1
		for j < numBytes && buffer[i] == buffer[j] {
			j++
		}
		j--
		if j == numBytes {
			return i //buffer is complete. use next buffer from last valid position to avoid lost of posible repetition
		}
		if j-i+1 >= r.minRepetitions {
			repetitions := (int)(j - i + 1)
			r.saveRepetition(buffer[i], repetitions, i)
			if r.verbose {
				r.displayRepetition(buffer[i], repetitions, buffer, i)
			}
		}
	}
	return numBytes //all bytes exhausted.
}

func (r *repetitions) slideDataFile() (int, error) {
	// Open file and create a buffered reader on top
	f, err := os.Open(r.inFileName)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	if r.startOn == 0 {
		persistInit(r.outName)
	} else {
		f.Seek(r.startOn, 0)
	}
	verbosePass := 5
	r.curDiskPtrRef = r.startOn
	bufferedReader := bufio.NewReader(f)
	buffer := make([]byte, r.bufferSize)
	var tbufferSize int
	for i := 1; ; i++ {
		numBytesRead, err := bufferedReader.Read(buffer)
		tbufferSize += numBytesRead
		effectiveBytesProcessed := (int64)(r.countRepetitions(buffer, numBytesRead))
		r.curDiskPtrRef += effectiveBytesProcessed
		f.Seek(r.curDiskPtrRef, 0) //aunque a leer el puntero cambia, es mejor reposicionar el puntero con los bytes efectivamente procesados
		if r.verbose && verbosePass == 5 {
			verbosePass = 0
			esc := "\u001b"
			reset := "[0m"
			print(esc + reset)
			fmt.Printf(esc+"[33m"+"\n%6.2f Gb proccessed. ", float32(tbufferSize)/float32(gigabyte))
			print(esc + reset)
		}
		verbosePass++
		if err == io.EOF {
			defer f.Close()
			return tbufferSize, nil
		}
		if err != nil {
			log.Fatal(err)
			defer f.Close()
			return tbufferSize, err
		}
	}
}

func resetTerminarColors() {
	esc := "\u001b"
	reset := "[0m"
	print(esc + reset)
}

func doScanForRepetitions(ifile string, ofile string, bufferSize int, minRepetitions int, startOnByte int64, verbose bool) error {
	var repStruct repetitions
	repStruct.inFileName = ifile
	repStruct.outName = ofile
	repStruct.startOn = startOnByte
	repStruct.verbose = verbose
	repStruct.minRepetitions = minRepetitions
	repStruct.bufferSize = bufferSize
	bytesProcessed, err := repStruct.slideDataFile()
	if err != nil {
		log.Fatal(err)
	}
	if repStruct.verbose {
		println("job done. Total digits analized=", bytesProcessed)
		println("output file was ", repStruct.outName)
	}
	return nil
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
			return "", true
		}
	}
	return "", false
}

func getCommandLineArguments() (inputFileName string, outputFileName string, bufferSize int, minRepetition int, startOn int, verboseOn bool, err error) {
	paramValue := ""
	present := false
	if inputFileName, present = getParamValue("-i"); !present {
		err = fmt.Errorf("error: data file is required")
		return
	}
	if outputFileName, present = getParamValue("-o"); !present {
		outputFileName = inputFileName + "-data-rep.txt"
	}
	if paramValue, present := getParamValue("-bMB"); !present {
		bufferSize = 1024 * 1024 * 1024 //default buffer is 1GB
	} else {
		if bufferSize, err = strconv.Atoi(paramValue); err != nil {
			err = fmt.Errorf("error: buffsize in MB is incorrect")
			return
		}
		bufferSize *= 1024 * 1024 //default buffer is 1GB
	}
	if paramValue, present = getParamValue("-r"); !present {
		minRepetition = 8
	} else {
		if minRepetition, err = strconv.Atoi(paramValue); err != nil {
			err = fmt.Errorf("error: the value for minimum of repetitions is invalid")
			return
		}
	}
	if _, present = getParamValue("-v"); present {
		verboseOn = true
	}
	if paramValue, present = getParamValue("-s"); !present {
		startOn = 0
	} else {
		if startOn, err = strconv.Atoi(paramValue); err != nil {
			err = fmt.Errorf("error: the value for star position is invalid")
			return
		}
	}
	return
}

func main() {
	resetTerminarColors()
	inputFileName, outputFileName, bufferSize, minRepetitions, startOn, verbose, err := getCommandLineArguments()
	if err != nil {
		println(err.Error())
		return
	}
	if verbose {
		println("verbose is On")
		println("analysing file '" + inputFileName + "'")
		println("out file name is '" + outputFileName + "' (if exist, results will be appended).")
		fmt.Printf("\nstarting from position = %d", startOn)
		fmt.Printf("\nmínimum repetitions = %d ", minRepetitions)
		fmt.Printf("\nbuffer size is %4.1fGB", (float32)(bufferSize)/1024.0/1024.0/1024.0)
	}
	if err := doScanForRepetitions(inputFileName, outputFileName, bufferSize, minRepetitions, int64(startOn), verbose); err != nil {
		println("ERROR: ", err.Error)
		return
	}
	return
}
