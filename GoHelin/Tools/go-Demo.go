package tools

import (
	list2 "container/list"
	"crypto/md5"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/go-gomail/gomail"
	"log"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
//gobook, ok := r.(map[string]interface{}) 判断类型匹配
//################################################################################
type Info struct {
	Name 	string 	`json:"name"`
	Age 	int		`json:"age"`
	Height	float64	`json:"height"`
	Weight 	float64	`json:"weight"`
}
//-----------------------测试Redis-----------------------
func testRedis(){
	myInfo := Info{"helin",25,163.0,58.5}
	conn , err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil{
		fmt.Println(err)
		return
	}
	defer conn.Close()
	str , err := json.Marshal(myInfo)
	if err != nil{
		fmt.Println(err)
		return
	}
	ret , err := conn.Do("SET","info",str)
	fmt.Println(ret)
	//-------------------------------------
	ret ,err = redis.String(conn.Do("GET","info"))
	if err != nil{
		fmt.Println(err)
		return
	}
	fmt.Println(ret)
	//--------------------------------------
	ret , err = redis.Values(conn.Do("HGETALL","family"))
	for _, v := range ret.([]interface{}) {
		fmt.Printf("%s\n", v)
	}

}
//-----------------------写json 文件---------------------
func writeJson(){
	w_file,err := os.Create("test.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	newEncode := json.NewEncoder(w_file)
	myInfo := Info{"helin",25,163.0,58.5}
	err = newEncode.Encode(myInfo)
	if err != nil{
		fmt.Println(err)
	}
}
//-----------------------读json文件--------------------------
func readJson(){
	r_file,err := os.Open("test.json")
	if err != nil{
		fmt.Println(err)
		return
	}
	newEncoder := json.NewDecoder(r_file)
	var pinfo Info
	err = newEncoder.Decode(&pinfo)
	if err != nil{
		fmt.Println(err)
		return
	}
	fmt.Println(pinfo.Name)
}
//-----------------------redis管道--------------------------
func redisPipe(){
	conn , err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil{
		fmt.Println(err)
		return
	}
	defer conn.Close()
	bak := "key_"
	for i := 0; i< 10000;i++{
		index := strconv.FormatInt(int64(i),10)
		key := bak + index
		conn.Send("set",key,index)
	}
	conn.Flush()
	for i := 0; i< 10000;i++{
		fmt.Println(redis.String(conn.Receive()))
	}
}
//################################################################################

//################################################################################
func main(){

}
//################################################################################
//-----------------------测试定时器,time.Since(time.now())可记录用时---------------------------
func TickerAndTimer(){
	ticker := time.NewTicker(time.Second)
	timer := time.NewTimer(time.Second*15)
	var count int
	for ; ;{
		select {
		case <-ticker.C:
			count++
			break
		case <-timer.C:
			goto STOP
		}
	}
STOP:
	fmt.Println("stop ...")
	fmt.Println(count)
}
//-----------------------多路复用select--------------------------
func testSelect(){
	cInt := make(chan int)
	cStr := make(chan string)
	cFloat := make(chan float64)

	go func(c chan int) {
		for ;; {
			c <- rand.Int()
			time.Sleep(time.Second * 2)
		}
	}(cInt)
	//-------------------------------
	go func(c chan string) {
		for ;;{
			sh := md5.Sum([]byte(strconv.FormatInt(int64(rand.Int()),10)))
			c <- fmt.Sprintf("%x",sh)
			time.Sleep(time.Second*4)
		}
	}(cStr)
	//-------------------------------
	go func(c chan float64) {
		for ;;{
			c <- rand.Float64()
			time.Sleep(time.Second*5)
		}
	}(cFloat)

	for;;{
		select {
		case <-cInt:
			fmt.Println("cInt : " , <-cInt)
			break
		case <-cStr:
			fmt.Println("cStr: " , <-cStr)
			break
		case <-cFloat:
			fmt.Println("cFloat : " , <-cFloat)
			break
		default:
			break
		}
	}
}
//-----------------------chan通道和waitgroup---------------------------
func testChan(){
	cInt := make(chan int)
	var w sync.WaitGroup
	w.Add(2)
	fmt.Println("start ...")
	go func(c chan int) {
		defer w.Done()
		time.Sleep(time.Second*5)
		c <- 10086
	}(cInt)

	go func(c chan int) {
		defer w.Done()
		fmt.Println(<- c)
	}(cInt)
	w.Wait()
}
//-----------------------go锁 sync.mutex---------------------------
var num int = 0
var wg sync.WaitGroup
func doTest(){
	fmt.Println(runtime.NumCPU())
	wg.Add(10)
	m := sync.Mutex{}
	for i:=0;i<10;i++{
		go dofun(uint(i),&m)
	}
	time.Sleep(time.Second)
	wg.Wait()
}
func dofun(num uint,mutex* sync.Mutex){
	defer wg.Done()
	mutex.Lock()
	num++
	time.Sleep(time.Second)
	fmt.Println(num)
	mutex.Unlock()
}
//################################################################################
//-----------------------go 字符类型和容器---------------------------
func goChar(){
	//go 字符类型有两种 byte 和 rune，byte = uint8 代表ascii；rune 代表一个utf-8字符 等价于 int32
	var c byte = 's'
	fmt.Printf("%c\n",c)
	var word rune = '锅'
	fmt.Printf("%c\n",word)
	//枚举用法
	const (
		Sunday  = iota
		Monday
		Tuesday
		Wednesday
		Thursday
		Friday
		Saturday
	)
	fmt.Println(Tuesday)
	//数组和切片的定义区别 ,切片中元素的删除需要用到切片和append
	var arry [10]int //数组
	var sarry []int  //切片
	for i :=0;i<10;i++{
		arry[i] = i
		sarry = append(sarry,i)
	}
	sarry = append(sarry[:0],sarry[1:]...) //或 sarry = sarry[1:]
	fmt.Println(arry)
	fmt.Println(sarry)
	//--------------go list--------------------------
	list := list2.New()
	list.PushBack(1)
	list.PushBack(10)
	list.PushBack(2)
	list.PushBack(5)
	fmt.Println(list.Front().Value,list.Back().Value)
	for i := list.Front();i != nil;i = i.Next(){
		fmt.Println(i.Value)
	}
	/*
		make 关键字的主要作用是创建切片、哈希表和 Channel 等内置的数据结构，
		而new 的主要作用是为类型申请一片内存空间，并返回指向这片内存的指针
	*/
}
//-----------------------go的一些特殊包------------------------------
func goPackage(){
	//-----------------sync包与锁:限制线程对变量的访问------------------
	//-----------------big包:对整数的高精度计算------------------------
	//-----------------image包:图像包制作GIF动画----------------------
	//-----------------regexp包:正则表达式---------------------------
	//-----------------os包：系统相关，包括信号文件创建等---------------
	//-----------------time包:时间相关,时间戳，定时器，计时器等----------
	//-----------------flag包：命令行参数解析--------------------------
	//-----------------go 发送电子邮件-------------------------------
	fromEmail := "1053206020@qq.com"
	toEmail := "1461671786@qq.com"
	m := gomail.NewMessage()
	m.SetAddressHeader("From", fromEmail, "贺林")
	m.SetAddressHeader("To", toEmail, "helin")
	m.SetHeader("Subject", "邮件测试")
	m.SetBody("text/html", "<h1>hello world</h1>")

	d := gomail.NewDialer("smtp.qq.com", 587, "1053206020@qq.com", "srusxeiybntsbcac")
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		fmt.Printf("***%s\n", err.Error())
	}
}
//################################################################################
//----------------------go http------------------------------------
func goHttpGet(){
	type Ip struct {
		Code	int		`json:"code"`
		Data	struct {
			Ip 		string	`json:"ip"`
			Country	string	`json:"country"`
			Area	string 	`json:"area"`
			Region 	string	`json:"region"`
			City 	string	`json:"city"`
			Isp		string 	`json:"isp"`
		}	`json:"data"`
	}
	result := new(Ip)
	ipAddr := "http://ip.taobao.com/service/getIpInfo.php?ip=" + "14.155.159.164"
	r,err := http.Get(ipAddr)
	failOnError(err,"request err")
	defer r.Body.Close()
	err = json.NewDecoder(r.Body).Decode(result)
	failOnError(err,"decoder err")
	fmt.Println(*result)
}
func goHttpPost(){

}
//################################################################################
