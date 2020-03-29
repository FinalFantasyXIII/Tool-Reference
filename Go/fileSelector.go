package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

type Counter struct {
	info []string
}
func SearchFile(path string, c *Counter){
	dirInfo,err := ioutil.ReadDir(path)
	if err != nil{
		fmt.Println(err)
		return
	}
	for _,v := range dirInfo{
		if v.IsDir() == true{
			SearchFile(path + "/" + v.Name(), c)
		}else{
			c.info = append(c.info,path+"/"+v.Name())
		}
	}
}


func Controller (path string) []string{
	c := &Counter{[]string{}}
	SearchFile(path,c)
	return c.info
}

func SortFileInfo(fInfo []string) {
	picture := []string{".jpg",".png",".gif",".bmp",".tif",".pcx"}
	video	:= []string{".mp4",".rm",".rmvb",".avi",".mkv",".wmv"}
	audio 	:= []string{".mp3",".wav",".ncm"}
	//打开3个文件，分别负责记录磁盘中所有视频，音频，图片资源位置
	wg := &sync.WaitGroup{}
	go Store("pics.txt",fInfo,wg,picture)
	go Store("video.txt",fInfo,wg,video)
	go Store("audio.txt",fInfo,wg,audio)
	wg.Wait()
}

func Store(fName string,fInfo []string ,wg *sync.WaitGroup,Ft []string){
	wg.Add(1)
	fd,err := os.OpenFile(fName,os.O_CREATE|os.O_TRUNC, 6)
	if err != nil{
		fmt.Println(fName,"打开失败")
		wg.Done()
		os.Exit(-1)
	}
	defer fd.Close()

	for _,vaule := range fInfo{
		for _,v := range Ft{
			if strings.HasSuffix(vaule,v) == false{
				continue
			}
			newValue := strings.ReplaceAll(vaule,"/","\\")
			fd.WriteString(newValue + "\n")
		}
	}
	wg.Done()
}

func main(){
	var path string
	fmt.Scanln(&path)
	info := Controller(path)

	fmt.Println(len(info))
	SortFileInfo(info)

	fmt.Scanln(&path)
}
