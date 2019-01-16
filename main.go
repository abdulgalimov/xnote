package main

import (
	"xnote"
)

func main() {
	err := xnote.Create()
	if err != nil {
		panic(err)
	}
	xnote.Start()
}