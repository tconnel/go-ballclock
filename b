#!/bin/sh
#go test github.com/emirpasic/gods/stacks/arraystack/arraystack_test.go
#go test github.com/emirpasic/gods/lists/arraylist/arraylist_test.go
go run main.go 30
go run main.go 45
go run main.go 30 --min=325
go build -buildmode=exe
