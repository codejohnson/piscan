# pirepeat
pirepeat is a tool to look for repetition sequences of digits on big files of pi

___
## Usage:

go run pirepeat.go -i:inputFile [-o:outputFile] [-v] [-s:startPosition] [-bMB:bufferSize] [-r:minimumRepetitions]

___

## Arguments explained:
**-i:**:INPUTFILE. The path to data filename</br>
**-o:**:OUTPUTFILE.  the output file name for the final statistics. If the file exist, result is appended. Default is the name of the input file with sufix ""-data-rep.txt"</br>
**-v** VERBOSE MODE. Allows you to se the progress every gigabyte of file processed. Default is Verbose</br>
**-s:**:STARTPOSITION. Is the physical start position in the file. Normallally, first file of pi have the integer part 3, followed by a dot(.). You must ommit this using -s:2</br>
**-r:**:MINIMUMREPETITIONS. Minimum size of repetition sequences to search and save in result file.</br>
**-bMB:**:BUFFERSIZE. Is the size of the memory buffer. Default is 1GB. Some old machines with low main memory can use little buffers. </br>

___

## samples:

    NO VERBOSE
    START IN POSITION 2 (first position on disk is zero)
    INPUT FILE IS /Volumes/Data/Pi/pi_dec_1t_01.txt

go run pirepeat.go -i:/Volumes/Data/Pi/pi_dec_1t_01.txt -s:2
___

    VERBOSE
    START IN POSITION 2 (first position on disk is zero)
    INPUT FILE IS /Volumes/Data/Pi/pi_dec_1t_01.txt
    OUTPUT FILE IS dataresult.txt
go run pirepeat.go -i:/Volumes/Data/Pi/pi_dec_1t_01.txt -o:dataresult.txt -s:2 -v


___

## How verbose option -v shows the information:
Jorges-iMac:pirepeat jorgejohnson$ go run pirepeat.go -i:/Volumes/Data/Pi/pi_dec_1t_03.txt -v -bMB:512 -r:10</br>
</br>
verbose is On</br>
analysing file '/Volumes/Data/Pi/pi_dec_1t_03.txt'</br>
out file name is '/Volumes/Data/Pi/pi_dec_1t_03.txt-data-rep.txt' (if exist, results will be appended).</br>
starting from position = 0</br>
mínimum repetitions = 10</br>
buffer size is  0.5GB</br>
9 (10) :...1374995177 9999999999 16599490435...</br>
5 (10) :...1625189876 5555555555 21941533653...</br>
8 (10) :...3199712313 8888888888 24586628067...</br>
9 (10) :...5362301114 9999999999 53502424972...</br>
1 (10) :...4310812926 1111111111 60473612322...</br>
3Gb aprox proccessed.</br>
4 (10) :...0954251575 4444444444 78657261384...</br>
2 (10) :...7408858649 2222222222 09412510873...</br>
