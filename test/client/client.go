// Copyright 2017 Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// An example implementation of a client.

package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"google.golang.org/grpc"

	"crypto/tls"

	"github.com/golang/sync/errgroup"
	pb "istio.io/manager/test/grpcecho"
)

var (
	count   int
	timeout time.Duration

	url       string
	headerKey string
	headerVal string
)

func init() {
	flag.IntVar(&count, "count", 1, "Number of times to make the request")
	flag.DurationVar(&timeout, "timeout", 60*time.Second, "Request timeout")
	flag.StringVar(&url, "url", "", "Specify URL")
	flag.StringVar(&headerKey, "key", "", "Header key")
	flag.StringVar(&headerVal, "val", "", "Header value")
}

func makeHTTPRequest(client *http.Client) func(i int) func() error {
	return func(i int) func() error {
		return func() error {
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				return err
			}

			log.Printf("[%d] Url=%s\n", i, url)
			if headerKey != "" {
				req.Header.Add(headerKey, headerVal)
				log.Printf("[%d] Header=%s:%s\n", i, headerKey, headerVal)
			}

			resp, err := client.Do(req)
			if err != nil {
				return err
			}

			log.Printf("[%d] StatusCode=%d\n", i, resp.StatusCode)

			data, err := ioutil.ReadAll(resp.Body)
			defer func() {
				if err = resp.Body.Close(); err != nil {
					log.Printf("[%d error] %s\n", i, err)
				}
			}()

			if err != nil {
				return err
			}

			for _, line := range strings.Split(string(data), "\n") {
				if line != "" {
					log.Printf("[%d body] %s\n", i, line)
				}
			}

			return nil
		}
	}
}

func makeGRPCRequest(client pb.EchoTestServiceClient) func(i int) func() error {
	return func(i int) func() error {
		return func() error {
			log.Printf("[%d] grpcecho.Echo\n", i)
			r, err := client.Echo(context.Background(),
				&pb.EchoRequest{Message: fmt.Sprintf("request #%d", i)},
			)
			if err != nil {
				log.Printf("[%d error] %v\n", i, err)
				return err
			}
			for _, line := range strings.Split(string(r.GetMessage()), "\n") {
				if line != "" {
					log.Printf("[%d body] %s\n", i, line)
				}
			}
			return nil
		}
	}
}

func main() {
	flag.Parse()
	var f func(int) func() error
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		/* #nosec */
		client := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
			Timeout: timeout,
		}
		f = makeHTTPRequest(client)
	} else if strings.HasPrefix(url, "grpc://") {
		address := url[len("grpc://"):]
		conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(timeout))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		client := pb.NewEchoTestServiceClient(conn)
		f = makeGRPCRequest(client)
		defer conn.Close()
	} else {
		log.Fatalf("Unrecognized protocol %q", url)
	}

	g, _ := errgroup.WithContext(context.Background())
	for i := 0; i < count; i++ {
		g.Go(f(i))
	}
	if err := g.Wait(); err != nil {
		log.Printf("Error %s\n", err)
	} else {
		log.Println("All requests succeeded")
	}
}
