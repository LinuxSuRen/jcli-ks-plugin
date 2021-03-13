package main

import (
	inner "github.com/linuxsuren/jcli-ks-plugin/cmd"
	"os"
)

func main() {
	cmd := inner.NewKSPlugin()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
