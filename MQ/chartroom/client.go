package main

import (
	"bufio"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"os"
	"time"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main(){
	con,err := amqp.Dial("amqp://helin:123@localhost:5672/learnTest")
	failOnError(err,"connect failed")
	defer con.Close()
	//-------------- 创建一个 channel -------------------------
	ch ,err := con.Channel()
	failOnError(err,"create channel failed")
	defer ch.Close()

	//创建一个协程用于发送信息
	go func(){
		for;;{
			reader := bufio.NewReader(os.Stdin)
			str, _ := reader.ReadString('\n')
			err := ch.Publish(
				"chatroom",
				"",
				false,
				false,
				amqp.Publishing{	//可以在此处指定发送人的一系列信息
					ContentType:     "text/message",
					Timestamp:       time.Now(),
					Body:            []byte(str),
				})
			failOnError(err,"send failed")
		}
	}()

	//先要订阅聊天室的交换机
	//创建一个队列用来绑定
	q,err := ch.QueueDeclare(
		"client1",    // name
		false, // durable
		false, // delete when usused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err,"create queue failed")
	//将队列跟交换机绑定
	err = ch.QueueBind(
		q.Name,
		"",
		"chatroom",
		false,
		nil)
	failOnError(err,"bind err")

	//开始接收交换机中的消息
	cm,err := ch.Consume(q.Name,"C1",false,false,false,false,nil)
	failOnError(err,"create consumer failed")
	forever := make(chan bool)
	//cm中包含了发送者的所有信息
	go func(){
		for d := range cm{
			fmt.Print(d.Timestamp," : ",string(d.Body))
			d.Ack(false)
		}
	}()
	<-forever
}
