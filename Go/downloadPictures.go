package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"
)

func main(){
	fmt.Println("*******************FBI WARNNING*******************")
	fmt.Println("此软件所下载的内容只供个人学习交流使用，禁止作商业用途！")
	fmt.Println("否则后果自负！")
	fmt.Println("*******************   USAGE   ********************")
	fmt.Println("第一步输入你想要下载到本地的路径，例如: D:/ 或 D:/pictures/")
	fmt.Println("第二步输入你想要下载的开始番号，按回车结束。例如: 输入10000 按回车结束")
	fmt.Println("第三步输入你想要下载的结束番号，按回车结束。例如: 输入10050 按回车结束")
	fmt.Println("上述操作会下载50组套图资源,当界面不再有任何信息输出时表示资源下载完毕，按任意键结束程序")
	fmt.Println("*******************   TIPS    ********************")
	fmt.Println("每组套图磁盘占用大致在20MB~30MB间，请考虑自身磁盘资源下载")
	fmt.Println("程序根据每100个资源为一组任务，请考虑自身计算机硬件资源合理使用")

	var wg sync.WaitGroup
	var localDir string
	var begin ,end int
	fmt.Println("请输入你想要下载到本地的路径:")
	fmt.Scanln(&localDir)
	fmt.Println("请输入你想要下载的开始番号:")
	fmt.Scanln(&begin)
	fmt.Println("请输入你想要下载的结束番号:")
	fmt.Scanln(&end)

	if begin > end || begin <0 || end <0 {
		fmt.Println("你输入的起止番号不符合规则！")
		os.Exit(-1)
	}
	for i := begin;i < end;{
		j := i
		i = i+100
		if i > end{
			i = end
		}
		go DownLoad(localDir,j,i,&wg)
	}
	//DownLoad("D:/pics/",10000,10003,&wg)
	wg.Wait()
	fmt.Scanln(&localDir)
}

func DownLoad(localDir string ,begin,end int,wg *sync.WaitGroup){
	wg.Add(1)
	url := "https://mtl.gzhuibei.com/images/img/"
	for index := begin; index <= end; index++ {
		picDir := strconv.Itoa(index) + "/"
		if err := os.Mkdir(localDir+picDir, os.ModePerm); err != nil{
			wg.Done()
			return
		}
		for i := 1; i < 500; i++ {
			link := url + picDir + strconv.Itoa(i) + ".jpg"
			fmt.Println(link)
			resp, _ := http.Get(link)
			if resp.StatusCode != 200 {
				fmt.Println(index)
				break
			}
			body, _ := ioutil.ReadAll(resp.Body)

			file, err := os.Create(localDir + picDir + strconv.Itoa(i) + ".jpg")
			if err != nil {
				panic(err)
			}
			io.Copy(file, bytes.NewReader(body))
			file.Close()
		}
	}
	wg.Done()
}