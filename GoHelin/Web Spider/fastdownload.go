package main

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
	"time"
)

func Start(){
	fmt.Println("*******************FBI WARNNING*******************")
	fmt.Println("此软件所下载的内容只供个人学习交流使用，禁止作商业用途！")
	fmt.Println("否则后果自负！")
	fmt.Println("*******************   USAGE   ********************")
	fmt.Println("第一步输入你想要下载到本地的路径，例如: D:/ 或 D:/pictures/")
	fmt.Println("第二步输入你想要下载的资源对象的URL，按回车结束。例如: https://www.meitulu.com/t/shiyijia-kittyjiang/")
	fmt.Println("上述操作会下载shiyijia-kittyjiang下所有的资源,当界面不再有任何信息输出时表示资源下载完毕，按任意键结束程序")
	fmt.Println("*******************   TIPS    ********************")
	fmt.Println("每组套图磁盘占用大致在1MB~30MB间，请考虑自身磁盘资源下载")
	fmt.Println("程序根据每100个资源为一组任务，请考虑自身计算机硬件资源合理使用")

	var localDir string
	var targetLink string
	fmt.Println("请输入你想要下载到本地的路径:")
	fmt.Scanln(&localDir)

	_,err := os.Stat(localDir)
	if err != nil{
		fmt.Println("目录不存在,程序将新建此目录")
		err:= os.MkdirAll(localDir, os.ModePerm)
		if err != nil{
			fmt.Println("目录创建失败,程序5s后自动退出...")
			time.Sleep(time.Second *5)
			os.Exit(-1)
		}
		fmt.Println("目录创建成功...")
	}

	fmt.Println("请输入你想要下载的资源对象的URL:")
	fmt.Scanln(&targetLink)

	results ,code  := GetPicsCode(targetLink)
	if code != 200{
		return
	}

	/*for i := begin;i < end;{
		j := i
		i = i+100
		if i > end{
			i = end
		}
		go MTL(localDir,j,i,&wg)
	}*/
	for _,value := range results{
		fmt.Println("----------------------------------------------------------------------------")
		fmt.Println(value)
		MTL(localDir,value)
	}
	//DownLoad("D:/pics/",10000,10003,&wg)
	fmt.Scanln(&localDir)
}

//---------------------取目标链接中的标题和图片地址-----------------------
type Result struct{
	Title string
	Link string
}

func GetPicsCode(targetLink string) ([]Result,int){
	var results []Result
	rsp,err := http.Get(targetLink)
	if err != nil{
		return results,-1
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != 200{
		return results,rsp.StatusCode
	}
	doc ,_ := goquery.NewDocumentFromReader(rsp.Body)
	doc.Find("li").Each(func(i int,selection *goquery.Selection) {
		s,flag:= selection.ChildrenFiltered(`a`).ChildrenFiltered("img").Attr("src")
		title,_ := selection.ChildrenFiltered(`a`).ChildrenFiltered("img").Attr("alt")
		if flag{
			rt := Result{title,strings.Replace(s,"0.jpg","",1)}
			results = append(results,rt)
		}
	})
	return results,rsp.StatusCode
}


func MTL(localDir string,rt Result ){
	ret := strings.Contains(rt.Title,"/")
	if ret{
		rt.Title = strings.ReplaceAll(rt.Title,"/",".")
	}
	picDir := rt.Title + "/"
	fmt.Println(picDir)
	if err := os.Mkdir(localDir+picDir, os.ModePerm); err != nil{
		return
	}
	for i := 1; i < 500; i++ {
		link := rt.Link + strconv.Itoa(i) + ".jpg"
		resp, err := http.Get(link)
		if err != nil{
			resp.Body.Close()
			break
		}
		if resp.StatusCode != 200 {
			resp.Body.Close()
			break
		}
		body, _ := ioutil.ReadAll(resp.Body)

		file, err := os.Create(localDir + picDir + strconv.Itoa(i) + ".jpg")
		if err != nil {
			file.Close()
			continue
		}
		io.Copy(file, bytes.NewReader(body))
		file.Close()
		resp.Body.Close()
		fmt.Println(link," : 下载完成")
	}
}

func FastDownLoad(){
	indexs := []string{
		"https://www.lanvshen.com/x/93/",
		"https://www.lanvshen.com/x/93/index_1.html",
		"https://www.lanvshen.com/x/93/index_2.html",
		"https://www.lanvshen.com/x/93/index_3.html",
		//"https://www.lanvshen.com/x/4/index_4.html",
		//"https://www.lanvshen.com/x/4/index_5.html",
		//"https://www.lanvshen.com/x/4/index_6.html",
		//"https://www.lanvshen.com/x/4/index_7.html",
		//"https://www.lanvshen.com/x/4/index_8.html",
		//"https://www.lanvshen.com/x/4/index_9.html",
		//"https://www.lanvshen.com/x/4/index_10.html",
	}
	var source []string

	for _,value := range indexs{
		rsp,err := http.Get(value)
		if err != nil{
			return
		}
		defer rsp.Body.Close()

		if rsp.StatusCode != 200{
			return
		}
		doc ,_ := goquery.NewDocumentFromReader(rsp.Body)
		doc.Find("li").Each(func(i int,selection *goquery.Selection) {
			s,flag:= selection.ChildrenFiltered(`a`).ChildrenFiltered("img").Attr("src")
			if flag{
				source = append(source,strings.Replace(s,"0.jpg","",1))
			}
		})
		time.Sleep(time.Second * 3)
	}
	count := 1;//图片名称
	for _,v := range source{
		for i := 1; i < 100; i++ {
			link := v + strconv.Itoa(i) + ".jpg"
			resp, err := http.Get(link)
			if err != nil{
				resp.Body.Close()
				break
			}
			if resp.StatusCode != 200 {
				resp.Body.Close()
				break
			}
			body, _ := ioutil.ReadAll(resp.Body)

			file, err := os.Create("d:/picture/Young Gangan/" + strconv.Itoa(count) + ".jpg")
			if err != nil {
				file.Close()
				continue
			}
			io.Copy(file, bytes.NewReader(body))
			file.Close()
			resp.Body.Close()
			count++
		}
	}
}



func WpeNet(){
	indexs := []string{
		"https://www.lanvshen.com/x/33/",
		"https://www.lanvshen.com/x/33/index_1.html",
		"https://www.lanvshen.com/x/33/index_2.html",
		"https://www.lanvshen.com/x/33/index_3.html",
		"https://www.lanvshen.com/x/33/index_4.html",
		"https://www.lanvshen.com/x/33/index_5.html",
		//"https://www.lanvshen.com/x/4/index_6.html",
		//"https://www.lanvshen.com/x/4/index_7.html",
		//"https://www.lanvshen.com/x/4/index_8.html",
		//"https://www.lanvshen.com/x/4/index_9.html",
		//"https://www.lanvshen.com/x/4/index_10.html",
	}
	type Source struct {
		Url 	string
		Actress string
	}
	var source []Source

	for _,value := range indexs{
		rsp,err := http.Get(value)
		if err != nil{
			return
		}
		defer rsp.Body.Close()

		if rsp.StatusCode != 200{
			return
		}
		doc ,_ := goquery.NewDocumentFromReader(rsp.Body)
		doc.Find("li").Each(func(i int,selection *goquery.Selection) {
			s,flag:= selection.ChildrenFiltered(`a`).ChildrenFiltered("img").Attr("src")
			if flag{
				url := strings.Replace(s,"0.jpg","",1)
				contex := selection.Text()
				arry := strings.Split(contex,"\n")
				last := arry[3]
				actress := strings.ReplaceAll(last," ","")
				source = append(source,Source{url,actress})
			}
		})
		time.Sleep(time.Second * 3)
	}
	count := 1;//图片名称
	for _,v := range source{
		arry := strings.Split(v.Url,"/")
		path := "d:/picture/WPE-net/" + v.Actress + "/" + arry[len(arry) - 2] + "/"
		os.MkdirAll(path, os.ModePerm)
		for i := 1; i < 200; i++ {
			link := v.Url + strconv.Itoa(i) + ".jpg"
			resp, err := http.Get(link)
			if err != nil{
				resp.Body.Close()
				break
			}
			if resp.StatusCode != 200 {
				resp.Body.Close()
				break
			}
			body, _ := ioutil.ReadAll(resp.Body)

			file, err := os.Create(path + strconv.Itoa(count) + ".jpg")
			if err != nil {
				file.Close()
				continue
			}
			io.Copy(file, bytes.NewReader(body))
			file.Close()
			resp.Body.Close()
			count++
		}
	}
}



func main(){
	WpeNet()
}