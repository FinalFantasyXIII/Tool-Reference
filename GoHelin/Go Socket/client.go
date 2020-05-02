package GoSocket

import (
	"fmt"
	"net"
	"os"
	"time"
)

func client(){
	c,err := net.Dial("tcp","192.168.31.220:1234")
	if err != nil{
		fmt.Println(err)
		os.Exit(-1)
	}

	for ;;{
		c.Write([]byte("hello,world&"))
		time.Sleep(time.Second*5)
	}
}
