#!/bin/bash

in_file=input

while IFS='' read -r line; do
    I=0
    while [ $I -lt ${#line} ]; do

        I=$[ $I + 1 ]
        regex=$(echo $line | sed "s/[a-z]/./$I")

        RES=$(grep "$regex" $in_file | wc -l)

        if [ $RES -eq 2 ]; then
            echo
            echo $regex | sed 's/\.//'
            exit 0
        fi
    done

done < "$in_file"
