package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
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

	//创建一个交换机
	err = ch.ExchangeDeclare(
		"chatroom",
		"fanout",
		true,
		false,
		false,
		false,
		nil)
	failOnError(err,"create exchange failed")

	//开始接收交换机中收到的信息并缓存
	//先创建一个接收消息的queue
	q,err := ch.QueueDeclare(
		"server",    // name
		true, // durable
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
	cm,err := ch.Consume(q.Name,"S",false,false,false,false,nil)
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