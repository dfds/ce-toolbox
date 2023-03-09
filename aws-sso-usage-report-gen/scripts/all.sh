#!/bin/sh
go run ../main.go general -i $1 | qsv table
go run ../main.go user-country -i $1 | qsv table
go run ../main.go user-activity -i $1 | qsv sort -s totalSignIns -R -N | qsv table
