package req

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/rachitsh92/better-http-client/client"
	res "github.com/rachitsh92/better-http-client/response"
)

type Requester interface {
	DoRequest(httpClient *client.HTTPClient) (res.Response, error)
}

type Request struct {
	method  string
	url     string
	headers string
	data    string
}

func NewReq(method, url, headers, data string) Requester {
	return &Request{
		method:  method,
		url:     url,
		headers: headers,
		data:    data,
	}
}

func (r *Request) DoRequest(httpClient *client.HTTPClient) (res.Response, error) {
	// Creating request
	log.Printf("inside DoRequest")

	req, err := http.NewRequest(r.method, r.url, nil)
	if err != nil {
		return res.Response{}, fmt.Errorf("error creating request: %v", err)
	}

	// Setting headers
	if r.headers != "" {
		for _, h := range splitHeaders(r.headers) {
			header := splitEachHeader(h)
			req.Header.Add(header[0], header[1])
		}
	}

	// Setting data
	if (r.method == "POST" || r.method == "PUT") && r.data != "" {
		req.Body = io.NopCloser(strings.NewReader(r.data))
	}

	log.Printf("Request: %s %s", req.Method, req.URL)
	
	for i := 0; i < httpClient.Retries; i++ {
		start := time.Now()
		resp, err := httpClient.Client.Do(req)
		duration := time.Since(start)
		log.Printf("Request completed in %v\n", duration)

		if err == nil {
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatalf("Error reading response body: %v", err)
			}

			var result map[string]interface{}
			if err := json.Unmarshal(body, &result); err != nil {
				log.Fatalf("Error unmarshalling response body: %v", err)
			}
			log.Printf("Status Code: %d",resp.StatusCode)
			return result, nil
		}

		log.Printf("request failed: %v. retrying(%d/%d)...", err, i+1, httpClient.Retries)
		time.Sleep(2 * time.Second)
	}

	return res.Response{}, fmt.Errorf("request failed after %d retries: %v", httpClient.Retries, err)
}


func splitHeaders(headers string) []string {
	return strings.Split(headers, ";")
}

func splitEachHeader(header string) []string {
	headerParts := strings.SplitN(header, ":", 2)
	if len(headerParts) != 2 {
		log.Fatal("Invalid header format supplied: %s", header)
	}
	return headerParts
}
