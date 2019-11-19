package main

import (
	"bytes"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"math/rand"
	"runtime"
	"strconv"
	"time"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
func GetGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}
//-------------------------------------------------------------------------------

func main(){
	setQueueLength()
}

//--------------------------模式1--------------------------------
func testMQ_1(){
	//-------------- 连接 MQ -------------------------
	con,err := amqp.Dial("amqp://helin:123@localhost:5672/learnTest")
	failOnError(err,"connect failed")
	defer con.Close()
	//-------------- 创建一个 channel -------------------------
	ch ,err := con.Channel()
	failOnError(err,"create channel failed")
	defer ch.Close()
	//-------------- 创建一个队列 ---------------
	q,err := ch.QueueDeclare(
		"BloodBorn",//消息名称
		false, //是否持久化
		false, //是否自动删除
		false, //是否排外
		false,
		nil) //队列中的消息什么时候会自动被删除
	failOnError(err,"QueueDeclare failed")

	body := "Farewell,good hunter.May you find your worth in the waking world!"
	//------------ 往创建的队列中发送消息 ----------
	err = ch.Publish(
		"",
		q.Name, //要发送的队列名称
		false,
		false,
		amqp.Publishing{
			ContentType:     "text/plain",
			Body:            []byte(body),
		})
	failOnError(err,"publish failed")
}
func consumeMQ_1(){
	//--------------------- 连接 MQ ---------------------------
	con,err := amqp.Dial("amqp://helin:123@localhost:5672/learnTest")
	failOnError(err,"connect to MQ failed")
	defer con.Close()
	//--------------------- 创建一个channel ----------------
	ch,err := con.Channel()
	failOnError(err,"create channel failed")
	defer ch.Close()
	//--------------------- 从MQ接受消息 -------------------
	//-------------------- 定义consumer 去消费队列中的消息 --------
	cm,err := ch.Consume(
		"BloodBorn", //要消费的队列名字
		"pony",// 消费者名字
		true,   // 自动应答，为true则MQ将消息发送出去久删除，为false则需要等待客户端应答后才把消息删除掉
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")
	//------------------ 打印cm中存取的消息-----------------
	forever := make(chan bool)
	go func() {
		for d := range cm {
			log.Printf("Received a message: %s", d.Body)
		}
	}()
	<-forever
}

//--------------------------模式2--------------------------------
func testMQ_2(){
	//-------------- 连接 MQ -------------------------
	con,err := amqp.Dial("amqp://helin:123@localhost:5672/learnTest")
	failOnError(err,"connect failed")
	defer con.Close()
	//-------------- 创建 channel -------------------------
	ch ,err := con.Channel()
	failOnError(err,"create channel failed")
	defer ch.Close()
	//-------------- 创建队列 --------------------
	q,err := ch.QueueDeclare(
		"BloodBorn",
		false, //是否持久化
		false, //是否自动删除
		false, //是否排外
		false,
		nil) //队列中的消息什么时候会自动被删除
	failOnError(err,"QueueDeclare failed")
	//------------- 往队列发100条消息 -------------
	for i := 0 ;i<100;i++{
		num := strconv.Itoa(rand.Int())
		e := ch.Publish(
			"",
			q.Name,
			false,
			false,
			amqp.Publishing{
				ContentType:"text/plain",
				Body: []byte(num),
			})
		failOnError(e,"publish failed")
	}
}
func consumeMQ_2(n uint64){
	//--------------------- 连接 MQ ---------------------------
	con,err := amqp.Dial("amqp://helin:123@localhost:5672/learnTest")
	failOnError(err,"connect to MQ failed")
	defer con.Close()
	//--------------------- 创建 channel ----------------
	ch,err := con.Channel()
	failOnError(err,"create channel failed")
	defer ch.Close()
	//-------------------- 定义consumer去消费 --------
	cm,err := ch.Consume(
		"BloodBorn", // queue
		"pony",// 消费者名字
		false,   // 设置为false，MQ就会等待客户端回应后再去处理这条消息
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")
	//------------------ get message from this channel:cm -----------------
	forever := make(chan bool)
	num  := time.Duration(n)
	go func() {
		for d := range cm {
			log.Printf("Received a message: %s : %d", d.Body, n)
			time.Sleep(time.Second * num)
			d.Ack(false)
		}
	}()
	<-forever
}

//-------------------------模式3---------------------------------
/*
	消费者可以根据consume函数的参数 auto-ack去设置是否回应MQ
	true：代表自动回应，MQ发送完消息给具体队列直接将消息删除
	false：代表需要等待客户端回应，客户端需要调用ack(false)将消息标记为已消费，
			否者MQ会永远保留这条消息，并发送给请求的连接，导致内存暴涨
	客户端可以使用channel.Qos()函数设置自己的消费能力，MQ推送数量
*/
func testMQ_3(){
	//连接到MQ
	conn,err := amqp.Dial("amqp://helin:123@localhost:5672/learnTest")
	failOnError(err,"connect failed")
	defer conn.Close()
	//创建一个channel
	ch,err := conn.Channel()
	failOnError(err,"create channel failed")
	defer ch.Close()
	//创建一个队列,一个持久化的队列
	q,err := ch.QueueDeclare(
		"dark soul",
		true, //持久化参数
		false,
		false,
		false,
		nil)
	failOnError(err,"create queue failed")
	//发送数据publish
	for i := 0;i<100;i++{
		num := strconv.Itoa(rand.Int())
		err = ch.Publish(
			"",
			q.Name,
			false,
			false,
			amqp.Publishing{
				ContentType:     "text/plain",
				DeliveryMode:    amqp.Persistent,//存取到磁盘
				Timestamp:       time.Time{},
				Body:            []byte(num),
			})
		failOnError(err,"send failed")
	}
}
func consumeMQ_3(n uint64){
	//连接到MQ
	conn,err := amqp.Dial("amqp://helin:123@localhost:5672/learnTest")
	failOnError(err,"connect failed")
	defer conn.Close()
	//创建一个channel
	ch,err := conn.Channel()
	failOnError(err,"create channel failed")
	defer conn.Close()
	//定义queue
	q, err := ch.QueueDeclare(
		"dark soul",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")
	//为channel设置每次接受数量
	err = ch.Qos(1,0,false)
	failOnError(err, "set Qos error")
	//定义consumer
	cm, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil)
	failOnError(err, "Failed to register a consumer")
	forever := make(chan bool)
	go func() {
		for d := range cm{
			fmt.Println(d.Body," ",d.Timestamp,":",n)
			num := time.Duration(n)
			time.Sleep(num*time.Second)
			err := d.Ack(false)
			failOnError(err, "Failed to register a consumer")
		}
	}()
	<-forever
}

//-------------------------模式4---------------------------------
/*
	发布/订阅 依靠交换机exchange，订阅者定义临时队列绑定到交换机，当发布者有新的
	消息发布时就会发到exchange中，exchange会把消息发送到所有绑定了本身的队列中。
*/
func testMQ_4(){
	//连接MQ
	con,err := amqp.Dial("amqp://helin:123@localhost:5672/learnTest")
	failOnError(err,"connect failed")
	defer con.Close()
	//创建一个channel
	ch,err := con.Channel()
	failOnError(err,"create channel failed")
	defer ch.Close()
	//创建一个exchange
	err = ch.ExchangeDeclare(
		"logs",
		"fanout",
		true,
		false,
		false,
		false,
		nil)
	failOnError(err,"create exchange failed")
	//往exchange发送消息
	for i:=0;i<100;i++{
		s := strconv.Itoa(rand.Int())
		err := ch.Publish(
			"logs",
			"",
			false,
			false,
			amqp.Publishing{
				ContentType:     "text/time",
				Timestamp:       time.Now(),
				Body:            []byte(s),
			})
		failOnError(err,"send failed")
	}
}
func consumeMQ_4(){
	//连接MQ
	con,err := amqp.Dial("amqp://helin:123@localhost:5672/learnTest")
	failOnError(err,"connect failed")
	defer con.Close()
	//创建channel
	ch,err := con.Channel()
	failOnError(err,"create channel failed")
	defer ch.Close()
	//创建exchange
	err = ch.ExchangeDeclare(
		"logs",
		"fanout",
		true,
		false,
		false,
		false,
		nil)
	failOnError(err,"create exchange failed")
	//创建队列
	q,err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when usused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err,"create queue failed")
	//绑定exchange 和 queue
	err = ch.QueueBind(
		q.Name,
		"",
		"logs",
		false,
		nil)
	failOnError(err,"bind err")
	//定义消费者
	cm,err := ch.Consume(q.Name,"tom",false,false,false,false,nil)
	failOnError(err,"create consumer failed")
	forever := make(chan bool)
	go func(){
		for d := range cm{
			fmt.Println(d.Body," : ",d.Timestamp)
			time.Sleep(time.Second)
			d.Ack(false)
		}
	}()
	<-forever
}

//------------------------模式5----------------------------------
/*
	路由模式通过在exchange中添加指定router-key来分发消息，
	当客户端连接绑定到具体的队列queue，带了router-key时，
	MQ就按照router-key的值将消息分类发送到指定绑定了router-key的队列中去
	客户端接着就可以调用consume去取这些消息
*/
func testMQ_5(){
	//连接MQ
	con,err := amqp.Dial("amqp://helin:123@localhost:5672/learnTest")
	failOnError(err,"connect failed")
	defer con.Close()
	//创建一个channel
	ch,err := con.Channel()
	failOnError(err,"create channel failed")
	defer ch.Close()
	//创建一个exchange，一个带路由的exchange
	err = ch.ExchangeDeclare("color","direct",true,false,false,false,nil)
	failOnError(err,"create exchange err")
	//往指exchange发送给指定路由的数据,publish的第二个参数
	for i:=0;i<20;i++{
		s := strconv.Itoa(rand.Int())
		value := amqp.Publishing{
			ContentType:     "text/time",
			Timestamp:       time.Now(),
			Body:            []byte(s),
		}
		err := ch.Publish("color", "red", false, false, value)
		failOnError(err,"send failed")
		time.Sleep(time.Second)
	}
}
func consumeMQ_5(){
	//连接MQ
	con,err := amqp.Dial("amqp://helin:123@localhost:5672/learnTest")
	failOnError(err,"connect failed")
	defer con.Close()
	//创建一个channel
	ch,err := con.Channel()
	failOnError(err,"create channel failed")
	defer ch.Close()
	//创建一个临时队列
	q,err := ch.QueueDeclare("",false,false,true,false,nil)
	failOnError(err,"create queue failed")
	//绑定到具体的exchange，并指明router-key 为 red
	err = ch.QueueBind(q.Name,"red","color",false,nil)
	failOnError(err,"bind failed")
	//定义consumer
	cm,err := ch.Consume(q.Name,"",false,false,false,false,nil)
	failOnError(err,"accept message err")
	forever := make(chan bool)
	dt := time.Duration(3)
	go func(){
		for d := range cm{
			fmt.Println(string(d.Body)," : ",d.Timestamp)
			d.Ack(false)
			time.Sleep(dt*time.Second)
		}
	}()
	<- forever
}

//------------------------模式6----------------------------------
func testMQ_6(){
	messages := map[string]string{"china.news":"猪肉涨价","china.edu":"大学生很多","china.people":"poor",
		"usa.news":"枪杀案","usa.edu":"expansive","usa.people":"mix"}
	//连接MQ
	con,err := amqp.Dial("amqp://helin:123@localhost:5672/learnTest")
	failOnError(err,"connect failed")
	defer con.Close()
	//创建channel
	ch,err := con.Channel()
	failOnError(err,"create channel failed")
	defer ch.Close()
	//创建一个topic exchange
	err = ch.ExchangeDeclare("news","topic",true,false,false,false,nil)
	failOnError(err,"exchange create failed")
	//往这个exchange发送消息
	for k,v := range messages{
		err := ch.Publish("news",k,false,false,
		amqp.Publishing{
			ContentType:     "text/news",
			Timestamp:       time.Now(),
			Body:            []byte(v),
		})
		failOnError(err,"publish failed")
		time.Sleep(time.Second*5)
	}
}
func consumeMQ_6_1(){
	con,err := amqp.Dial("amqp://helin:123@localhost:5672/learnTest")
	failOnError(err,"connect failed")
	defer con.Close()
	//创建channel
	ch,err := con.Channel()
	failOnError(err,"create channel failed")
	defer ch.Close()
	//创建一个临时队列
	q,err := ch.QueueDeclare("",false,false,true,false,nil)
	failOnError(err,"create queue failed")
	//绑定queue到具体exchange，并指明topic
	err = ch.QueueBind("","#.news","news",false,nil)
	failOnError(err,"bind err")
	err = ch.QueueBind("","#.edu","news",false,nil)
	failOnError(err,"bind err")
	//开始接受消息
	cm,err := ch.Consume(q.Name,"",true,false,false,false,nil)
	failOnError(err,"consume err")
	forever := make(chan bool)
	go func(){
		for d := range cm{
			fmt.Println(string(d.Body),":",d.Timestamp)
		}
	}()
	<-forever
}
func consumeMQ_6_2(){
	con,err := amqp.Dial("amqp://helin:123@localhost:5672/learnTest")
	failOnError(err,"connect failed")
	defer con.Close()
	//创建channel
	ch,err := con.Channel()
	failOnError(err,"create channel failed")
	defer ch.Close()
	//创建一个临时队列
	q,err := ch.QueueDeclare("",false,false,true,false,nil)
	failOnError(err,"create queue failed")
	//绑定queue到具体exchange，并指明topic
	err = ch.QueueBind("","usa.#","news",false,nil)
	failOnError(err,"bind err")
	err = ch.QueueBind("","china.#","news",false,nil)
	failOnError(err,"bind err")
	//开始接受消息
	cm,err := ch.Consume(q.Name,"",true,false,false,false,nil)
	failOnError(err,"consume err")
	forever := make(chan bool)
	go func(){
		for d := range cm{
			fmt.Println(string(d.Body),":",d.Timestamp)
		}
	}()
	<-forever
}

//-----------------------设置队列中消息的条数-----------------------
func setQueueLength(){
	//连接MQ
	con,err := amqp.Dial("amqp://helin:123@localhost:5672/learnTest")
	failOnError(err,"connect err")
	//创建channel
	ch,err := con.Channel()
	failOnError(err,"cereate channel err")
	//创建队列Queue
	config := map[string]interface{}{"x-max-length":30}
	q,err := ch.QueueDeclare("onSale",true,false,false,false,config)
	failOnError(err,"create queue err")
	//发送消息
	for i := 0 ;i<100;i++{
		num := strconv.Itoa(rand.Int())
		e := ch.Publish(
			"",
			q.Name,
			false,
			false,
			amqp.Publishing{
				ContentType:"text/plain",
				Body: []byte(num),
			})
		failOnError(e,"publish failed")
	}
}
