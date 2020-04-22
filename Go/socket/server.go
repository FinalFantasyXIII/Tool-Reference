package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func test(){
	tchan := make(chan int,100)

	for i:= 0;i<10;i++{
		tchan <- i
	}

	go func (){
		for ;;{
			i := <-tchan
			fmt.Println(i)
		}
	}()

	close(tchan)
	time.Sleep(time.Second * 5)
}
func main(){
	fd ,err := net.Listen("tcp","192.168.1.102:1234")
	if err != nil{
		fmt.Println(err)
		os.Exit(-1)
	}
	defer fd.Close()
	for ;; {
		new_fd,err := fd.Accept()
		if err != nil{
			fmt.Println(err)
			continue
		}
		arry := make([]byte,100)
		go func(c net.Conn,arry []byte){
			for ;;{
				c.Read(arry)
				str := string(arry)
				s := strings.Split(str,"&")[0]
				fmt.Println(s)
			}
		}(new_fd,arry)
	}
}
