package docker

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var update = flag.Bool("update", false, "update .golden files")

// Test_Docker_Auth_Encode_Error verifies that the required secrets must be
// given in order to properly encode the docker config JSON. The encoded string
// is used to template into e.g. pull secrets.
//
//     go test ./... -run Test_Docker_Auth_Encode_Error
//
func Test_Docker_Auth_Encode_Error(t *testing.T) {
	testCases := []struct {
		secrets map[string]string
		matcher func(error) bool
	}{
		// Case 0 ensures that the docker password must be given in order to
		// properly compute the docker config JSON.
		{
			secrets: map[string]string{
				"docker.registry": "<registry>",
				"docker.username": "<username>",
			},
			matcher: IsExecutionFailed,
		},
		// Case 1 ensures that the docker registry must be given in order to
		// properly compute the docker config JSON.
		{
			secrets: map[string]string{
				"docker.password": "<password>",
				"docker.username": "<username>",
			},
			matcher: IsExecutionFailed,
		},
		// Case 2 ensures that the docker username must be given in order to
		// properly compute the docker config JSON.
		{
			secrets: map[string]string{
				"docker.password": "<password>",
				"docker.registry": "<registry>",
			},
			matcher: IsExecutionFailed,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var err error

			var a AuthEncoder
			{
				c := AuthConfig{
					Secrets: tc.secrets,
				}

				a, err = NewAuth(c)
				if err != nil {
					t.Fatal(err)
				}
			}

			_, err = a.Encode()
			if !tc.matcher(err) {
				t.Fatal("matcher must match error")
			}
		})
	}
}

// Test_Docker_Auth_Encode_Golden verifies the proper computation of the base64
// encoded docker config JSON. The encoded string is used to template into e.g.
// pull secrets.
//
//     go test ./... -run Test_Docker_Auth_Encode_Golden -update
//
func Test_Docker_Auth_Encode_Golden(t *testing.T) {
	testCases := []struct {
		secrets map[string]string
		encoded string
	}{
		// Case 0 ensures that the docker config JSON can be computed given the
		// minimal amount of secret data.
		{
			secrets: map[string]string{
				"docker.password": "<password>",
				"docker.registry": "<registry>",
				"docker.username": "<username>",
			},
		},
		// Case 1 ensures that the docker config JSON can be computed given more
		// secret data than necessary.
		{
			secrets: map[string]string{
				"aws.access.id":     "<id>",
				"aws.access.secret": "<secret>",
				"docker.password":   "<foo>",
				"docker.registry":   "<bar>",
				"docker.username":   "<baz>",
			},
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var err error

			var a AuthEncoder
			{
				c := AuthConfig{
					Secrets: tc.secrets,
				}

				a, err = NewAuth(c)
				if err != nil {
					t.Fatal(err)
				}
			}

			var actual []byte
			{
				enc, err := a.Encode()
				if err != nil {
					t.Fatal(err)
				}
				actual = []byte(fmt.Sprintf("%s\n", enc))
			}

			p := filepath.Join("testdata/auth", fileName(i))
			if *update {
				err := ioutil.WriteFile(p, actual, 0644) // nolint:gosec
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
