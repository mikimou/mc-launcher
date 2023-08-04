package main

import (
	"github.com/hashicorp/go-getter/cmd/go-getter"
)

func main() {
	client := getter.Client{
		// This will prevent copying or writing files through symlinks
		DisableSymlinks: true,
	}
}
