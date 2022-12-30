package tests

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

const port = "8899"

func TestInteroperability(t *testing.T) {
	// webrpc server <=> generated client
	{
		version := "v0.9.1"
		webrpcGen := fmt.Sprintf("./bin/webrpc-gen@%s", version)
		webrpcTest := fmt.Sprintf("./bin/webrpc-test@%s", version)

		{ // Download webrpc-gen binary.
			if _, err := os.Stat(webrpcGen); err != nil && os.IsNotExist(err) {
				binaryURL := fmt.Sprintf("https://github.com/webrpc/webrpc/releases/download/%s/webrpc-gen.%s-%s", version, runtime.GOOS, runtime.GOARCH)
				if out, err := exec.Command("curl", "-o", webrpcGen, "-LJO", binaryURL).CombinedOutput(); err != nil {
					t.Fatalf("%v:\n%s", err, string(out))
				}
				if out, err := exec.Command("chmod", "+x", webrpcGen).CombinedOutput(); err != nil {
					t.Fatalf("%v:\n%s", err, string(out))
				}
			}
		}

		{ // Download webrpc-test binary.
			if _, err := os.Stat(webrpcTest); err != nil && os.IsNotExist(err) {
				binaryURL := fmt.Sprintf("https://github.com/webrpc/webrpc/releases/download/%s/webrpc-test.%s-%s", version, runtime.GOOS, runtime.GOARCH)
				if out, err := exec.Command("curl", "-o", webrpcTest, "-LJO", binaryURL).CombinedOutput(); err != nil {
					t.Fatalf("%v:\n%s", err, string(out))
				}
				if out, err := exec.Command("chmod", "+x", webrpcTest).CombinedOutput(); err != nil {
					t.Fatalf("%v:\n%s", err, string(out))
				}
			}
		}

		if out, err := exec.Command("bash", "-c", fmt.Sprintf("%s -print-schema > api.ridl", webrpcTest)).CombinedOutput(); err != nil {
			t.Fatalf("%v:\n%s", err, string(out))
		}

		if out, err := exec.Command("bash", "-c", fmt.Sprintf("%s -schema=./api.ridl -target=golang -pkg=tests -client -server -out=./test.gen.go", webrpcGen)).CombinedOutput(); err != nil {
			t.Fatalf("%v:\n%s", err, string(out))
		}

		// Run server.
		server := exec.Command(webrpcTest, "-server", fmt.Sprintf("-port=%v", port))
		if err := server.Start(); err != nil {
			t.Fatal(err)
		}
		defer server.Process.Kill()

		// Wait for server to be ready.
		if out, err := exec.Command("bash", "-c", fmt.Sprintf("until nc -z localhost %v; do sleep 0.1; done;", port)).CombinedOutput(); err != nil {
			t.Fatalf("%v:\n%s", err, string(out))
		}

		err := RunTests(fmt.Sprintf("http://localhost:%v", port))
		assert.NoError(t, err)
	}
}
