package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sgw_user_migrate/sgw"
	"time"
)

func main() {
	//read from flags
	syncGatewayURLPointer := flag.String("url", "http://localhost:4985", "Sync Gateway URL")
	dbPointer := flag.String("db", "offline_reads", "Database name")

	flag.Parse()

	syncGatewayURL := *syncGatewayURLPointer
	db := *dbPointer
	//open a file
	file, err := os.Open("files/sample.json")
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	//open an error file if it doesn't exist
	errorFile, err := os.OpenFile("files/error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		panic(err)
	}
	defer func(errorFile *os.File) {
		err := errorFile.Close()
		if err != nil {
			panic(err)
		}
	}(errorFile)

	//create a new decoder and decode the file into a struct
	decoder := json.NewDecoder(file)
	var queryResult []sgw.CouchbaseSgwUser
	err = decoder.Decode(&queryResult)
	if err != nil {
		panic(err)
	}

	jobNums := len(queryResult)
	workerNums := 10

	results := make(chan string, jobNums)
	jobs := make(chan int, jobNums)

	now := time.Now()

	service := sgw.NewService(syncGatewayURL, db, errorFile)

	//loop through the queryResult and print the values
	for i := 0; i < workerNums; i++ {
		go service.WriteToSyncGateway(&queryResult, results, jobs, i)
	}

	for i := 0; i < jobNums; i++ {
		jobs <- i
	}

	close(jobs)

	for i := 0; i < len(queryResult); i++ {
		fmt.Println(<-results)
	}

	println("Time taken: ", time.Since(now))
}
