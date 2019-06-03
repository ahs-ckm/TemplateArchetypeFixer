#!/usr/bin/bash

#FILES="$1"
FILES=./templates-di/*
for f in $FILES
do
 echo "Processing $f..."
 ./ArchetypeFixer "https://ahsckm.ca" 1175.115.70 20180614 am9uLmJlZWJ5OlBhNTV3b3Jk "$f"
done
