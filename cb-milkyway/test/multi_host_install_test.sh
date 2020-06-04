#!/bin/bash

for var in "$@"
do
    ./full_test.sh "$var" install &
done


