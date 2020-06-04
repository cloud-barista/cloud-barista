#!/bin/bash

for var in "$@"
do
    ./full_test.sh "$var" &
done


