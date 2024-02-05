package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func main() {
	var home string

	flag.StringVar(&home, "home", "/etc/shsw", "ShieldSweep home directory (where shsw.json is located)")
	flag.Parse()

	sock := filepath.Join(home, "shsw.sock")
	_, err := os.Stat(sock)
	if err != nil {
		panic(errors.Wrapf(err, "Error checking for socket file: %v\n", sock))
	}

	httpc := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", sock)
			},
		},
	}

	command := flag.Args()

	switch command[0] {
	case "run":
		response, err := httpc.Get("http://unix/run")
		if err != nil {
			panic(err)
		}
		fmt.Println(response.Body)
	case "status":
		response, err := httpc.Get("http://unix/status")
		if err != nil {
			panic(err)
		}
		fmt.Println(response.Body)
	}
}
