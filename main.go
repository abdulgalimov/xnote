package main

import (
	"github.com/labstack/gommon/log"
	"xnote"
)

func main() {
	//token := "265574930:AAFJh4-gQaJHc3Iw-J1Fo1GA-kXh21hTyj4"
	err := xnote.Create()
	if err != nil {
		log.Panic(err)
	}
	xnote.Start()

}

/*

import (
	"fmt"
	"time"
)

func main() {
	c1 := make(chan string)
	c2 := make(chan string)

	go func() {
		time.Sleep(time.Second * 2)
		c1 <- "data 1"
	}()
	go func() {
		time.Sleep(time.Second * 3)
		c2 <- "data 2"
	}()

	var data1 string
	var data2 string
	func() {
		var completeNum = 0
		for completeNum < 2 {
			select {
			case msg1 := <- c1:
				fmt.Println(msg1)
				completeNum++
				data1 = msg1
			case msg2 := <- c2:
				fmt.Println(msg2)
				completeNum++
				data2 = msg2
			}
		}
	}()

	fmt.Println("exit", data1, data2)
}
*/
