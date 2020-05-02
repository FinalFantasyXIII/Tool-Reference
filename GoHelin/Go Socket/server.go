package GoSocket

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func Server(){
	fd ,err := net.Listen("tcp","192.168.31.220:1234")
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
