# picount
picount is a tool to get a count on the digits of pi in large repositories.
this command line util allows you to get your own count of the digits of pi.

___
## Usage:

go run picount.go -i:inputFile [-o:outputFile] [-v] [-s:startPosition]

___

## Arguments explained:
**-i**:INPUTFILE. The path to filename</br>
**-o**:OUTPUTFILE.  the output file name for the final statistics. If the file exist, result is appended</br>
**-v** VERBOSE MODE. Allows you to se the progress every gigabyte of file processed. 
Default is Verbose</br>
**-s**:STARTPOSITION. Is the physical start position in the file. Normallally, first file of pi have the integer part 3, followed by a dot(.). You must ommit this using -s:2
___

## samples:

    NO VERBOSE
    START IN POSITION 2 (first position on disk is zero)
    INPUT FILE IS /Volumes/Data/Pi/pi_dec_1t_01.txt

go run picount.go -i:/Volumes/Data/Pi/pi_dec_1t_01.txt -s:2
___

    VERBOSE
    START IN POSITION 2 (first position on disk is zero)
    INPUT FILE IS /Volumes/Data/Pi/pi_dec_1t_01.txt
    OUTPUT FILE IS dataresult.txt
go run picount.go -i:/Volumes/Data/Pi/pi_dec_1t_01.txt -o:dataresult.txt -s:2 -v


___

## How Verbose Show the information:
Every gigabyte of data proccessed, you see the following information for the ten digits

[0]=   751587464;[1]=   751654381;[2]=   751634627;[3]=   751612377;[4]=   751567578;[5]=   751621242;[6]=   751615969;[7]=   751629511;[8]=   751622137;[9]=   751647482
(7.000000 Gb total proccessed). ReadTime:11.180175194s ; ProcTime:2.710680122s

the final line shows the number of gigabytes advanced, and the times of read and processing
