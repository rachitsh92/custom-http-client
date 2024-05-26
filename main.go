package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/rachitsh92/better-http-client/client"
	req "github.com/rachitsh92/better-http-client/request"
	res "github.com/rachitsh92/better-http-client/response"
)

func main() {
	url := flag.String("url", "defaut", "url to fetch the result from")
	method := flag.String("X", "GET", "HTTP method to use")
	headers := flag.String("H", "", "HTTP headers")
	data := flag.String("d", "", "HTTP POST data")
	timeout := flag.Int("timeout", 10, "Tiomeout for the request in seconds")
	retries := flag.Int("retries", 3, "Number of retries on failure")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "something went wrong\n available flags :\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *url == "" {
		fmt.Println("url not found.")
		flag.Usage()
		os.Exit(1)
	}

	fmt.Printf("initiating contact with %s\n", *url)

	// creating a custom http client
	clientConfig := client.Config{
		Timeout: *timeout,
		Retries: *retries,
	}
	httpClient := client.NewClient(clientConfig)

	//make request
	request := req.NewReq(*method, *url, *headers, *data)

	resp, err := request.DoRequest(httpClient)
	if err != nil {
		log.Fatal("request failed: %v", err)
	}

	fmt.Printf("Response Body: %v", string(res.FormattedResponse(resp)))

}
