# pirepeat
pirepeat is a tool to look for repetition sequences of digits on big files of pi

___
## Usage:

go run pirepeat.go -i:inputFile [-o:outputFile] [-v] [-s:startPosition] [-bMB:bufferSize] [-r:minimumRepetitions]

___

## Arguments explained:
**-i:**:INPUTFILE. The path to data filename</br>
**-o:**:OUTPUTFILE.  the output file name for the final statistics. If the file exist, result is appended. Default is the name of the input file with sufix ""-data-rep.txt". -o option can be ommited, and no file output is generated.</br>
**-v** VERBOSE MODE. Allows you to se the progress every gigabyte of file processed. Default is Verbose</br>
**-s:**:STARTPOSITION. Is the physical start position in the file. Normallally, first file of pi have the integer part 3, followed by a dot(.). You must ommit this using -s:2</br>
**-min:**:MINIMUMREPETITIONS. Minimum size of repetition sequences to search and save in result file. Default value is 8</br>
**-max:**:MAXIMUMREPETITIONS. Maximum size of repetition sequences to search and save in result file. Default value is 100</br>
**-bMB:**:BUFFERSIZE. Is the size of the memory buffer. Default is 1GB. Some old machines with low main memory can use little buffers. </br>
**-new:**:RESTART. new file for statistics created. If file already exist, is removed before process start! </br>

___

## samples:

    NO VERBOSE
    START IN POSITION 2 (first position on disk is zero)
    INPUT FILE IS /Volumes/Data/Pi/pi_dec_1t_01.txt

go run pirepeat.go -i:/Volumes/Data/Pi/pi_dec_1t_01.txt -s:2
___

    VERBOSE
    START IN POSITION 0
    INPUT FILE IS /Volumes/Data/Pi/pi_dec_1t_01.txt
    OUTPUT FILE IS dataresult.txt
    MINIMUM REPETITION SEQUENCE OF 6 DIGITS 
    MAXIMUM REPETITION SEQUENCE OF 8 DIGITS 
go run pirepeat.go -i:/Volumes/Data/Pi/pi_dec_1t_01.txt -o:dataresult.txt -v -min:6 -max:9


___

## How verbose option -v shows the information:
Jorges-iMac:pirepeat jorgejohnson$ go run pirepeat.go -i:/Volumes/Data/Pi/pi_dec_1t_01.txt -s:2 -v  -min:10 -new -o:repet-file-pi01-min10.txt</br>
</br>
-verbose is On</br>
-restart = true</br>
-analysing file '/Volumes/Data/Pi/pi_dec_1t_01.txt'</br>
-out file name is 'repet-file-pi01-min10.txt' (if exist, results will be appended).</br>
-starting from position = 2</br>
-minimum repetitions = 10</br> 
-maximum repetitions = 100</br>
-buffer size is  1.0GB</br>
-file repet-file-pi01-min10.txtdo not exist.</br>
-output file repet-file-pi01-min10.txt was created.</br>
6 (10) :...4599643705666666666691436679246... position > 386980413</br>
8 (10) :...5455953577888888888821240339476... position > 3040319544</br>
1 (10) :...4930498886111111111172782238193... position > 3961184002</br>
3 (10) :...0298728519333333333386780471802... position > 4663739960</br>
4Gb aprox proccessed.</br>
5 (10) :...3459767569555555555573236183315... position > 7644991299</br>
3 (10) :...4576934174333333333396634111295... position > 8313901579</br>
0 (10) :...8009159279000000000029541559034... position > 8324296436</br>
