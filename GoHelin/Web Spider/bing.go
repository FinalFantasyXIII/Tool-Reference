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
)


func GetBYPictures(localDir string,index int){
	var wg sync.WaitGroup
	dom := "https://bing.ioliu.cn/?p="
	for i:=1;i<=index;i++{
		url := dom + strconv.Itoa(i)
		rsp,err := http.Get(url)
		if err != nil{
			fmt.Println(err)
			continue
		}
		if rsp.StatusCode != 200{
			fmt.Println(rsp.StatusCode)
			continue
		}
		//处理html
		doc ,_ := goquery.NewDocumentFromReader(rsp.Body)
		var result []string
		doc.Find("a").Each(func(i int,selection *goquery.Selection) {
			s,_:= selection.Attr("href")
			if strings.Contains(s,"?force=download") {
				result = append(result,s)
			}
		})
		rsp.Body.Close()
		go Catch(&wg ,result,localDir)
	}
	wg.Wait()
}


func Catch(wg *sync.WaitGroup,urls []string,localDir string){
	front := "https://bing.ioliu.cn"
	wg.Add(1)
	for _,s := range urls{
		name := strings.Split(strings.Split(s,"/")[2],"_")[0]
		link := front+s
		fmt.Println("link : ",link)
		resp, err := http.Get(link)
		if err != nil{
			resp.Body.Close()
			continue
		}
		if resp.StatusCode != 200 {
			resp.Body.Close()
			continue
		}

		body, _ := ioutil.ReadAll(resp.Body)
		file, _:= os.Create(localDir + name + ".jpg")
		io.Copy(file, bytes.NewReader(body))

		resp.Body.Close()
		file.Close()
	}
	wg.Done()
}


func BingRun(){
	fmt.Println("输入你要下载的磁盘位置,例如 : D:/pics/")
	var path string
	var index int
	fmt.Scanln(&path)
	fmt.Println("输入你要下载的资源页数,例如 : 5代表下载1-5页")
	fmt.Scanln(&index)
	GetBYPictures(path,index)
}
