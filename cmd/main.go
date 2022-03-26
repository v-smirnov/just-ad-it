package main

import (
	"flag"
	"fmt"

	"github.com/v-smirnov/just-ad-it/internal/app/infrastructure"
	"github.com/v-smirnov/just-ad-it/internal/app/service"
)

const defaultParallel = 10

func main() {
	var numParallel int
	flag.IntVar(&numParallel, "parallel", defaultParallel, "number of parallel requests")
	flag.Parse()

	urlList := flag.Args()

	if len(urlList) == 0 {
		fmt.Println("Please provide one or more url to send requests")
		return
	}

	requestService := service.NewService(infrastructure.NewClient(), int8(numParallel))

	for _, result := range requestService.DoRequests(urlList) {
		fmt.Println(result)
	}
}
