package tests

import (
	"context"
	"flag"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/webrpc/gen-golang/tests/client"
	"github.com/webrpc/gen-golang/tests/server"
)

var (
	serverFlag = flag.Bool("server", true, "run server")
	portFlag   = flag.Int("port", 9889, "run server at port")
	// -serverTimeout, not to be confused with -timeout flag provided by go test
	serverTimeoutFlag = flag.Duration("serverTimeout", 2*time.Second, "exit after given timeout")

	clientFlag = flag.Bool("client", true, "run client tests")
	urlFlag    = flag.String("url", "http://localhost:9889", "run client tests against given server URL")
)

func TestInteroperability(t *testing.T) {
	if *serverFlag {
		bind := fmt.Sprintf("0.0.0.0:%v", *portFlag)

		fmt.Printf("Running generated test server at %v\n", bind)
		srv, err := server.RunTestServer(bind, *serverTimeoutFlag)
		assert.NoError(t, err)

		if *clientFlag {
			// Close server after we finish running client tests.
			defer srv.Close()
		}

		defer srv.Wait()
	}

	if *clientFlag {
		fmt.Println("Running generated client tests against", *urlFlag)

		err := client.RunTests(context.Background(), *urlFlag)
		assert.NoError(t, err)
	}
}
