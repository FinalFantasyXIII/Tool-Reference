package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"golang.org/x/crypto/pbkdf2"
	"io"
	"os"
	"os/signal"
)

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(uint8(padding))}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	return origData, nil
}


//------------------check file if exists--------------------------
func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func EncodeToFile(src , dest string , key []byte){
	var err error
	//--------check src
	flag := checkFileIsExist(src)
	if !flag{
		fmt.Printf("%s 不存在,请确认\n",src)
		return
	}
	f1,err := os.Open(src)
	if err != nil{
		fmt.Printf("%s 打开失败,系统将退出\n",src)
		return
	}
	defer f1.Close()
	//--------check dest
	var f2 *os.File
	flag = checkFileIsExist(dest)
	if !flag{
		fmt.Printf("%s 不存在,系统将创建.....\n",dest)
		f2,err = os.Create(dest)
		if err != nil{
			fmt.Printf("%s 创建失败,系统将退出\n",dest)
			return
		}
	}else {
		f2 ,err = os.OpenFile(dest, os.O_WRONLY|os.O_TRUNC,0666)
		if err != nil{
			fmt.Printf("%s 打开失败,系统将退出\n",dest)
			return
		}
		fmt.Printf("%s 存在,系统将覆盖此文件\n" , dest)
	}
	defer f2.Close()

	//-------------start encrypt
	buff := make ([]byte , 256)
	tmp := make([]byte , 2)
	total := 0
	for {
		n , err := f1.Read(buff)
		if err != nil{
			if err == io.EOF {
				break
			}
			fmt.Println("读系统错误 , 文件前 %d 字节数据已被加密\n",total)
			return
		}
		buf1 , _ := AesEncrypt(buff[:n] , key)
		blen := uint16(len(buf1))
		binary.BigEndian.PutUint16(tmp,blen)
		finalBuff := make([]byte,2+blen)
		copy(finalBuff,tmp)
		copy(finalBuff[2:] , buf1)
		_ , err = f2.Write(finalBuff)
		if err != nil{
			fmt.Printf("写系统错误 , 文件前 % 字节数据已被加密\n" , total)
			return
		}
		total +=  int(blen)
	}
	fmt.Println("文件加密完毕.......")
}

func DecodeFromFile(src , dest string, key []byte){
	var err error
	//--------check src
	flag := checkFileIsExist(src)
	if !flag{
		fmt.Printf("%s 不存在,请确认\n",src)
		return
	}
	f1,err := os.Open(src)
	if err != nil{
		fmt.Printf("%s 打开失败,系统将退出\n",src)
		fmt.Println(err)
		return
	}
	defer f1.Close()
	//--------check dest
	var f2 *os.File
	flag = checkFileIsExist(dest)
	if !flag{
		fmt.Printf("%s 不存在,系统将创建.....\n",dest)
		f2,err = os.Create(dest)
		if err != nil{
			fmt.Printf("%s 创建失败,系统将退出\n",dest)
			return
		}
	}else {
		f2 ,err = os.OpenFile(dest, os.O_WRONLY|os.O_TRUNC,0666)
		if err != nil{
			fmt.Printf("%s 打开失败,系统将退出\n",dest)
			return
		}
		fmt.Printf("%s 存在,系统将覆盖此文件\n" , dest)
	}
	defer f2.Close()

	//------------start decode
	tmp := make([]byte,2)
	total := 0
	for   {
		_, err = io.ReadFull(f1, tmp)
		if err != nil{
			if err == io.EOF{
				break
			}
			fmt.Printf("读系统错误 , 文件前 %d 数据已被解密\n" , total)
			return
		}
		data_len := binary.BigEndian.Uint16(tmp)
		buff := make([]byte , data_len)

		_ , err = io.ReadFull(f1, buff)
		if err != nil{
			if err == io.EOF{
				break
			}
			fmt.Printf("读系统错误 , 文件前 %d 数据已被解密\n" , total)
			return
		}

		//decode
		newbuf , err := AesDecrypt(buff, key)
		if err != nil {
			fmt.Printf("解密失败 , 文件前 %d 数据已被解密\n" , total)
			return
		}

		_,err = f2.Write(newbuf)
		if err != nil{
			fmt.Printf("写系统错误 , 文件前 %d 数据已被解密\n" , total)
			return
		}

		total += 2
		total += int(data_len)
	}
	fmt.Println("文件解密完毕......")
}

func GenAESEncryptKey(key , salt string , length int) []byte{
	return  pbkdf2.Key([]byte(key), []byte(salt), 1024, length, sha1.New)
}

func ChoiceWhichToExcute(option int , src , dest string, key []byte){
	switch option {
	case 1:
		EncodeToFile(src,dest,key)
		break
	case 2:
		DecodeFromFile(src,dest,key)
		break
	default:
		fmt.Println("错误操作号 ......")
		break
	}
}

func main(){
	fmt.Println("-----------------Welcom-----------------")
	fmt.Println("-----------请输入 1(加密) 或 2(解密) , 回车确认------------")
	var choice int
again:
	fmt.Scanln(&choice)
	fmt.Println("-------> your input : " , choice)
	switch choice {
	case 1:
	case 2:
		break
	default:
		fmt.Println("请重新输入正确的操作号码")
		goto again
	}

	var src  string
	var dest string
	var key string
	var salt string
	if choice == 1{
		fmt.Println("请输入需要加密的文件目录 , 按回车确认")
	}else{
		fmt.Println("请输入需要解密的文件目录 , 按回车确认")
	}
	fmt.Scanln(&src)
	fmt.Println("-------> your input : " , src)
	if choice == 1 {
		fmt.Println("请输入加密后文件的存放目录 , 按回车确认")
	}else {
		fmt.Println("请输入解开密后文件的存放目录 , 按回车确认")
	}
	fmt.Scanln(&dest)
	fmt.Println("-------> your input : " , dest)
	fmt.Println("请输入加密key 要记住哦, 按回车确认")
	fmt.Scanln(&key)
	fmt.Println("-------> your input : " , key)
	fmt.Println("请输入加密盐值 要记住哦 嫌麻烦可以key相同, 按回车确认")
	fmt.Scanln(&salt)
	fmt.Println("-------> your input : " , salt)

	code := GenAESEncryptKey(key , salt, 32)

	ChoiceWhichToExcute(choice ,src, dest, code)

	waitForSignal()
}

func waitForSignal() os.Signal {
	signalChan := make(chan os.Signal, 1)
	defer close(signalChan)
	signal.Notify(signalChan, os.Kill, os.Interrupt)
	s := <-signalChan
	signal.Stop(signalChan)
	return s
}
