package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func main(){
	c,err := net.Dial("tcp","192.168.1.102:1234")
	if err != nil{
		fmt.Println(err)
		os.Exit(-1)
	}

	for ;;{
		c.Write([]byte("hello,world&"))
		time.Sleep(time.Second*5)
	}
}
