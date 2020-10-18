package docker

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var update = flag.Bool("update", false, "update .golden files")

// Test_Docker_Auth TODO
//
//     go test ./... -run Test_Docker_Auth -update
//
func Test_Docker_Auth(t *testing.T) {
	testCases := []struct {
		input  map[string]string
		output string
	}{
		// Case 0 TODO
		{
			input: map[string]string{
				"aws.access.id":     "<id>",
				"aws.access.secret": "<secret>",
				"docker.password":   "<password>",
				"docker.registry":   "<registry>",
				"docker.username":   "<username>",
			},
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var err error

			var a AuthEncoder
			{
				c := AuthConfig{
					Bytes: mustBytes(tc.input),
				}

				a, err = NewAuth(c)
				if err != nil {
					t.Fatal(err)
				}
			}

			actual, err := a.Encode()
			if err != nil {
				t.Fatal(err)
			}

			p := filepath.Join("testdata/auth", fileName(i))
			if *update {
				err := ioutil.WriteFile(p, []byte(actual), 0644) // nolint:gosec
				if err != nil {
					t.Fatal(err)
				}
			}

			expected, err := ioutil.ReadFile(p)
			if err != nil {
				t.Fatal(err)
			}

			if !bytes.Equal(expected, []byte(actual)) {
				t.Fatalf("\n\n%s\n", cmp.Diff(string(actual), string(expected)))
			}
		})
	}
}

func fileName(i int) string {
	return "case-" + strconv.Itoa(i) + ".golden"
}

func mustBytes(v interface{}) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	return b
}
