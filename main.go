package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)
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
func main(){

}

//--------------------------------------------------
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
//--------------------------------------------------
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
//--------------------------------------------------
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
//--------------------------------------------------
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