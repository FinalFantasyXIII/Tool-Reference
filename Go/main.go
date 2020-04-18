package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)


//用来获取html文件
func test(){
	rsp,_ := http.Get("http://www.mgsdaigou9.com/390jac-040.html")
	if rsp.StatusCode != 200{
		os.Exit(-1)
	}

	f,_ := os.Create("rr.html")
	body, _ := ioutil.ReadAll(rsp.Body)
	io.Copy(f, bytes.NewReader(body))
	rsp.Body.Close()
}
func main(){
	Excute()
}

func Excute(){
	f,_ := os.Create("details.txt")
	encoder := json.NewEncoder(f)
	link := "http://www.mgsdaigou9.com/luxu/page/1"
	urls,_ := GetSinglePageLink(link)
	for _,v := range urls{
		fmt.Println(v)
		info,_ := AnalysisSinglePage(v)
		if err := encoder.Encode(info);err != nil{
			fmt.Println(err)
		}
	}
}

type Info struct{
	Code		string		`json:"code"`
	Number 		int			`json:"number"`
	Tag			string		`json:"tag"`
	Pics		[]string	`json:"pics"`
}

//link = http://www.mgsdaigou9.com/259luxu-1263.html"
func AnalysisSinglePage(link string)  (Info,error){
	rsp,err := http.Get(link)
	if err != nil{
		fmt.Println(err)
		return Info{},err
	}
	defer rsp.Body.Close()

	doc ,_ := goquery.NewDocumentFromReader(rsp.Body)

	//获取标题
	//title := doc.Find("h1").Text()

	//获取图片url
	var pics []string
	doc.Find("img").Each(func(i int,selection *goquery.Selection) {
		if s,b := selection.Attr("class");b == true && strings.Contains(s,"size-full") {
			picUrl,_ := selection.Attr("src")
			pics = append(pics,picUrl)
		}
	})

	code := strings.Split(strings.Split(link,"/")[3],".")[0]
	number,_ := strconv.Atoi( strings.Split(code,"-")[1])
	info := Info{
		Code: code,
		Number:number,
		Tag: strings.Split(code,"-")[0],
		Pics: pics,
	}
	return info,nil
}

//link = "http://www.mgsdaigou9.com/page/1"  20-url/page
func GetSinglePageLink(link string)([]string,error){
	var result []string
	rsp,_ := http.Get(link)
	if rsp.StatusCode != 200{
		return result,nil
	}

	doc ,_ := goquery.NewDocumentFromReader(rsp.Body)
	doc.Find("h2").Each(func(i int,selection *goquery.Selection) {
		s,_:= selection.ChildrenFiltered(`a`).Attr("href")
		result = append(result,s)
	})
	rsp.Body.Close()
	return result,nil
}

