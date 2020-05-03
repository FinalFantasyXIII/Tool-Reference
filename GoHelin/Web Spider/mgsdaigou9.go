package webspider

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

/*
	这个文档是比较完全的爬虫项目
*/

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


type Info struct{
	Code		string		`json:"code"`
	Number 		int			`json:"number"`
	Tag			string		`json:"tag"`
	Pics		[]string	`json:"pics"`
}

func main(){
	dir,index := Input()
	Mgsdaigou9(dir,index)
	var i int
	fmt.Scanln(&i)
}

func Mgsdaigou9(dir string,index int){
	var wg sync.WaitGroup
	url  := "http://www.mgsdaigou9.com/page/"
	for i := 1;i <= index;i++{
		link := url + strconv.Itoa(i)
		fmt.Println(link)
		go Excute(link,dir,&wg)
	}
	wg.Wait()
}

//输入信息提示
func Input() (string,int){
	var dir string
	var index int
	fmt.Println("请输入你要下载到本地的目录按回车结束,例如: d:/pics/")
	fmt.Scanln(&dir)
	_,err := os.Stat(dir)
	if err != nil{
		fmt.Println("目录不存在,程序马上退出")
		time.Sleep(time.Second*3)
		os.Exit(-1)
	}
	fmt.Println("请输入你要下载页面数量按回车结束,例如: 5")
	fmt.Scanln(&index)
	if index < 1 {
		fmt.Println("你的输入有误,程序马上退出")
		time.Sleep(time.Second*3)
		os.Exit(-1)
	}
	return dir,index
}

//开多个携程进行抓取
func Excute(link ,dir string,wg *sync.WaitGroup){
	wg.Add(1)
	defer wg.Done()
	urls,code := GetSinglePageLink(link)
	if code != 200 {
		return
	}
	for _,v := range urls{
		info,code := AnalysisSinglePage(v)
		if code != 200{
			continue
		}
		SavePicture(info,dir)
	}
}

//下载并存储图片
func SavePicture(info Info,dir string){
	//先检查&创建种类目录
	path := dir + info.Tag + "/" + info.Code + "/"
	fmt.Println(path)
	os.MkdirAll(path, os.ModePerm)
	//开始下载存储
	for i,v := range info.Pics{
		resp, _ := http.Get(v)
		if resp.StatusCode != 200 {
			resp.Body.Close()
			continue
		}
		body, _ := ioutil.ReadAll(resp.Body)

		file, err := os.Create(path + strconv.Itoa(i+1) + ".jpg")
		if err != nil {
			continue
		}
		io.Copy(file, bytes.NewReader(body))
		file.Close()
	}
}

//解析单个页面 , http://www.mgsdaigou9.com/259luxu-1263.html"
func AnalysisSinglePage(link string)  (Info,int){
	rsp,err := http.Get(link)
	if err != nil{
		return Info{},-1
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != 200{
		return Info{},rsp.StatusCode
	}
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
	return info,rsp.StatusCode
}

//获取单个页面的url , "http://www.mgsdaigou9.com/page/1"  20-url/page
func GetSinglePageLink(link string)([]string,int){
	var result []string
	rsp,err := http.Get(link)
	if err != nil{
		return result,-1
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != 200{
		return result,rsp.StatusCode
	}

	doc ,_ := goquery.NewDocumentFromReader(rsp.Body)
	doc.Find("h2").Each(func(i int,selection *goquery.Selection) {
		s,_:= selection.ChildrenFiltered(`a`).Attr("href")
		result = append(result,s)
	})
	return result,rsp.StatusCode
}

