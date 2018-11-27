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

const gigabyte int = 1024 * 1024 * 1024

type counter struct {
	summary        [10]uint64
	bytesProcessed int64
	inFileName     string
	outName        string
	verbose        bool
	startOn        int64
}

func (c *counter) showCount() {
	esc := "\u001b"
	reset := "[0m"
	escSec := [10]string{esc + "[31m", esc + "[32m", esc + "[33m", esc + "[34m", esc + "[35m", esc + "[36m", esc + "[37m", esc + "[31m", esc + "[32m", esc + "[33m"}
	for i := 0; i <= 9; i++ {
		fmt.Printf("%s%12d", escSec[i]+"["+strconv.Itoa(i)+"]="+esc+reset, c.summary[i])
		if i != 9 {
			print(";")
		}
	}
}

func (c *counter) countDigits(bytes []byte, numBytes int) {
	for i := 0; i < numBytes; i++ {
		c.summary[bytes[i]-48]++
	}
	if c.verbose {
		c.showCount()
	}
}

func (c *counter) slideDataFile(bufferSize int) (int, error) {
	// Open file and create a buffered reader on top
	f, err := os.Open(c.inFileName)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	if c.startOn > 0 {
		f.Seek(c.startOn, 0)
	}
	bufferedReader := bufio.NewReader(f)
	byteSlice := make([]byte, bufferSize)
	var tbufferSize int
	for i := 1; ; i++ {
		readStart := time.Now()
		numBytesRead, err := bufferedReader.Read(byteSlice)
		readElapsed := time.Since(readStart)
		tbufferSize += numBytesRead
		procStart := time.Now()
		c.countDigits(byteSlice, numBytesRead)
		procElapsed := time.Since(procStart)
		if c.verbose {
			esc := "\u001b"
			reset := "[0m"
			print(esc + reset)
			fmt.Printf(esc+"[33m"+"\n(%f Gb total proccessed). ReadTime:%s ; ProcTime:%s\n", float64(tbufferSize)/float64(gigabyte), readElapsed, procElapsed)
		}
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

func (c *counter) saveStats(filename string) {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	for i := 0; i <= 9; i++ {
		digit := strconv.Itoa(i)
		total := strconv.FormatUint(c.summary[i], 10)
		f.WriteString(digit + ";" + total + "\n")
	}
}

func resetTerminarColors() {
	esc := "\u001b"
	reset := "[0m"
	print(esc + reset)
}

func doCount(ifile string, ofile string, startOnByte int64, verbose bool) error {
	bufferSize := gigabyte
	var stats counter
	stats.inFileName = ifile
	stats.outName = ofile
	stats.startOn = startOnByte
	stats.verbose = verbose
	bytesProcessed, err := stats.slideDataFile(bufferSize)
	if err != nil {
		log.Fatal(err)
	}
	if stats.verbose {
		esc := "\u001b"
		reset := "[0m"
		println(esc + reset)
		escSec := [10]string{esc + "[31m", esc + "[32m", esc + "[33m", esc + "[34m", esc + "[35m", esc + "[36m", esc + "[37m", esc + "[31m", esc + "[32m", esc + "[33m"}
		println(esc + reset)
		for i := 0; i <= 9; i++ {
			fmt.Printf("%s%11.5f%%", escSec[i]+"["+strconv.Itoa(i)+"]="+esc+reset, float64(stats.summary[i])/float64(bytesProcessed)*100.0)
			if i != 9 {
				print(", ")
			}
		}
	}
	stats.saveStats("digitCountStat.txt")
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
		}
	}
	return
}

func getCommandLineArguments() (inputFileName string, outputFileName string, startOn int, verboseOn bool, err error) {
	paramValue := ""
	present := false
	if inputFileName, present = getParamValue("-i"); !present {
		err = fmt.Errorf("error: data file is required")
		return
	}
	if outputFileName, present = getParamValue("-o"); !present {
		outputFileName = "output.txt"
	}
	if _, present = getParamValue("-v"); !present {
		verboseOn = false
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
	inputFileName, outputFileName, startOn, verbose, err := getCommandLineArguments()
	if err != nil {
		println(err.Error())
		return
	}
	if verbose {
		println("verbose is On")
		println("analysing file '" + inputFileName + "'")
		println("out file name is '" + outputFileName + "' (if exist, results will be appended).")
	}
	if err := doCount(inputFileName, outputFileName, int64(startOn), verbose); err != nil {
		println("ERROR: ", err.Error)
		return
	}
	return
}
