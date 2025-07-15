// Your mission, should you choose to accept it:
//
// 1. Create a Github repository containing this code.
//
// 2. Write some tests for this service. You may modify the code to make it
// easier to test.
//
// 3. Create a Dockerfile that builds this service into a multiplatform
// amd64/arm64 Docker image.
//
// 4. Setup Github Actions (or your CI/CD provider of choice) so that:
//
//    - When a pull request is opened, the tests run.
//
//    - When a pull request is merged, a Docker image is built and pushed to ECR.
//
// 5. Be ready to answer questions about your work! We will ask you to walk us
// through how to run your Docker image.

package main

import (
	"log"
	"net/http"
	"os"

	"strongdm/handler"
)

func main() {
	bindAddr := os.Getenv("BIND_ADDR")
	log.Println("Listening on " + bindAddr)

	h := handler.New()
	log.Fatal(http.ListenAndServe(bindAddr, http.HandlerFunc(h.HandleRequest)))
}
