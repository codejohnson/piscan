package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const gigabyte int = 1000 * 1000 * 1000

type repetitions struct {
	counts         [15][10]int
	bytesProcessed int64
	inFileNames    string
	outName        string
	countFileName  string
	verbose        bool
	startOn        int64
	curDiskPtrRef  int64
	diskHits       int64
	minRepetitions int
	maxRepetitions int
	bufferSize     int
	restart        bool
}

func persistInit(filename string) {
	var _, err = os.Stat(filename)
	// create file if not exists
	if !os.IsNotExist(err) {
		var err = os.Remove(filename)
		if err != nil {
			panic(err)
		}
	} else {
		print("\n-file '", filename, "' do not exist.")
	}
	_, err = os.OpenFile(filename, os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	} else {
		print("\n-output file '", filename, "' was created.")
	}
}

func moveToFilePosition(f *os.File, restartPosition int64) {
	if restartPosition == 0 {
		print("\n-scanning from the beginning...")
	} else {
		print("\n-scanning from position $restartPosition...")
		f.Seek(restartPosition, 0)
		print("\n-restarted!")
	}
}

func (r *repetitions) displayRepetition(ifile string, digit byte, repetitions int, buffer []byte, bufferPosition int) {
	esc := "\u001b"
	reset := "[0m"
	fmt.Printf("\n%s > %c (%d) :...", ifile, digit, repetitions)
	from := bufferPosition - 10
	to := bufferPosition + repetitions + 10
	for i := from; i <= len(buffer) && i <= to; i++ {
		if i == bufferPosition || i == to-10 {
			print(esc + "[35m")
		}
		if i == to-10 {
			print(esc + reset)
		}
		fmt.Printf("%c", buffer[i])
	}
	print("... position > ", r.curDiskPtrRef+(int64)(bufferPosition))
}

func (r *repetitions) saveRepetition(ifilename string, digit byte, repetitions int, bufferPosition int) {
	f, err := os.OpenFile(r.outName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.WriteString(fmt.Sprintf("%s;%c;%d;%d\n", ifilename, digit, repetitions, r.curDiskPtrRef+(int64)(bufferPosition)))
}

func (r *repetitions) saveCounts() {
	fptr, err := os.OpenFile(r.countFileName, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer fptr.Close()
	for f := 0; f < len(r.counts); f++ {
		fptr.WriteString(fmt.Sprintf("%d;", f))
		for c := 0; c <= 9; c++ {
			fptr.WriteString(fmt.Sprintf("%d", r.counts[f][c]))
			if c < 9 {
				fptr.WriteString(";")
			}
		}
		fptr.WriteString("\n")
	}
}

func (r *repetitions) showCounts() {
	for f := 0; f < len(r.counts); f++ {
		fmt.Printf("\n%d;", f)
		for c := 0; c <= 9; c++ {
			fmt.Printf("%d", r.counts[f][c])
			if c < 9 {
				fmt.Printf(";")
			}
		}
	}
}

func (r *repetitions) countRepetitions(ifilename string, buffer []byte, numBytes int) int {
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
		if r.minRepetitions <= j-i+1 && j-i+1 <= r.maxRepetitions {
			repetitions := (int)(j - i + 1)
			r.counts[repetitions][buffer[i]-48]++
			if r.outName != "" {
				r.saveRepetition(ifilename, buffer[i], repetitions, i)
			}
			if r.verbose {
				r.displayRepetition(ifilename, buffer[i], repetitions, buffer, i)
			}
		}
	}
	return numBytes //all bytes exhausted.
}

func (r *repetitions) slideDataFile(inputFileName string) (int64, error) {
	f, err := os.Open(inputFileName)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	if r.restart {
		if r.outName != "" {
			persistInit(r.outName)
		}
	}
	f.Seek(r.startOn, 0)
	r.diskHits++
	verbosePass := 1
	r.curDiskPtrRef = r.startOn
	bufferedReader := bufio.NewReader(f)
	buffer := make([]byte, r.bufferSize)
	for i := 1; ; i++ {
		numBytesRead, err := bufferedReader.Read(buffer)
		r.diskHits++
		effectiveBytesProcessed := (int64)(r.countRepetitions(inputFileName, buffer, numBytesRead))
		r.curDiskPtrRef += effectiveBytesProcessed
		f.Seek(r.curDiskPtrRef, 0) //aunque a leer el puntero cambia, es mejor reposicionar el puntero con los bytes efectivamente procesados
		if r.verbose && verbosePass == 5 {
			verbosePass = 0
			esc := "\u001b"
			reset := "[0m"
			print(esc + reset)
			print("\n"+esc+"[33m", r.curDiskPtrRef/(int64)(gigabyte), "Gb aprox proccessed.")
			print(esc + reset)
		}
		verbosePass++
		if err == io.EOF {
			f.Close()
			return r.curDiskPtrRef, nil
		}
		if err != nil {
			log.Fatal(err)
			f.Close()
			return r.curDiskPtrRef, err
		}
	}
}

func (r *repetitions) slideDataFiles() (int64, error) {
	files := strings.Split(r.inFileNames, ",")
	var procbytes int64
	for _, filename := range files {
		diskPointer, err := r.slideDataFile(filename)
		procbytes += diskPointer
		if err != nil {
			return procbytes, err
		}
	}
	if r.countFileName != "" {
		r.saveCounts()
	}
	if r.verbose {
		r.showCounts()
	}
	return procbytes, nil
}

func resetTerminarColors() {
	esc := "\u001b"
	reset := "[0m"
	print(esc + reset)
}

func doScanForRepetitions(ifile string, ofile string, countfile string, bufferSize int, minRepetitions int, maxRepetitions int, startOnByte int64, restart bool, verbose bool) error {
	var repStruct repetitions
	repStruct.inFileNames = ifile
	repStruct.outName = ofile
	repStruct.countFileName = countfile
	repStruct.startOn = startOnByte
	repStruct.verbose = verbose
	repStruct.minRepetitions = minRepetitions
	repStruct.maxRepetitions = maxRepetitions
	repStruct.bufferSize = bufferSize
	repStruct.restart = restart
	start := time.Now()
	bytesProcessed, err := repStruct.slideDataFiles()
	elapsed := time.Since(start)
	println("\nanalysis took %s", elapsed)

	if err != nil {
		log.Fatal(err)
	}
	if repStruct.verbose {
		println("-job done. Total digits analized=", bytesProcessed)
		println("-input files : ", repStruct.inFileNames)
		if repStruct.outName != "" {
			println("-output file is: ", repStruct.outName)
		} else {
			println("-no output file for repetitions defined.")
		}
		if repStruct.countFileName != "" {
			println("-count file is: ", repStruct.countFileName)
		} else {
			println("-no output file for final statistics defined.")
		}
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

func getCommandLineArguments() (inputFileName string, outputFileName string, countFileName string, bufferSize int, minRepetition int, maxRepetition int, startOn int, restart bool, verboseOn bool, err error) {
	paramValue := ""
	present := false
	if inputFileName, present = getParamValue("-i"); !present {
		err = fmt.Errorf("error: data file is required")
		return
	}
	println(":::::", inputFileName)
	outputFileName, present = getParamValue("-o")

	countFileName, present = getParamValue("-c")

	if paramValue, present := getParamValue("-bMB"); !present {
		bufferSize = gigabyte //default buffer is 1GB
	} else {
		if bufferSize, err = strconv.Atoi(paramValue); err != nil {
			err = fmt.Errorf("error: buffsize in MB is incorrect")
			return
		}
		bufferSize *= 1000 * 1000
	}
	if paramValue, present = getParamValue("-min"); !present {
		minRepetition = 5
	} else {
		if minRepetition, err = strconv.Atoi(paramValue); err != nil {
			err = fmt.Errorf("error: the value for minimum of repetitions is invalid")
			return
		}
	}
	if paramValue, present = getParamValue("-max"); !present {
		maxRepetition = 100
	} else {
		if maxRepetition, err = strconv.Atoi(paramValue); err != nil {
			err = fmt.Errorf("error: the value for maximum of repetitions is invalid")
			return
		}
	}
	if _, present = getParamValue("-v"); present {
		verboseOn = true
	}
	if _, present = getParamValue("-new"); present {
		restart = true
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
	inputFileName, outputFileName, countFileName, bufferSize, minRepetitions, maxRepetitions, startOn, restart, verbose, err := getCommandLineArguments()
	if err != nil {
		println(err.Error())
		return
	}
	if verbose {
		println("-verbose is On")
		println("-restart =", restart)
		println("-analysing data '" + inputFileName + "'")
		if outputFileName != "" {
			println("-out file name is '" + outputFileName + "' (if exist, results will be appended).")
		}
		fmt.Printf("-starting from position = %d", startOn)
		fmt.Printf("\n-minimum repetitions = %d ", minRepetitions)
		fmt.Printf("\n-maximum repetitions = %d ", maxRepetitions)
		fmt.Printf("\n-buffer size is %4.1fGB", (float32)(bufferSize)/1000.0/1000.0/1000.0)
	}
	if err := doScanForRepetitions(inputFileName, outputFileName, countFileName, bufferSize, minRepetitions, maxRepetitions, int64(startOn), restart, verbose); err != nil {
		print("-ERROR: ", err.Error)
		return
	}
	return
}
