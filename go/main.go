package main

import "github.com/katsuokaisao/sql-play/cmd"

func main() {
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
