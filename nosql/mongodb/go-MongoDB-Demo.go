package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"sync"
	"time"
)

//-----------------mgo---------------------
type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	Job  string `json:"job"`
}

func DbInsert() {
	session, err := mgo.Dial("")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("helin").C("person")
	err = c.Insert(&Person{"helin", 26, "programer"},
		&Person{"heqin", 14, "student"},
		&Person{"hehan", 13, "student"})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("insert over")
	time.Sleep(time.Second * 5)
	var res []Person
	err = c.Find(bson.M{}).All(&res)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
}
func DbFind() {
	session, err := mgo.Dial("")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("helin").C("person")
	var res []Person
	err = c.Find(bson.M{"age": bson.M{"$gt": 13}}).All(&res) //db.person.find({"age" : {$gt:13}})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
}
func DbUpdate() {
	session, err := mgo.Dial("")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("helin").C("person")
	err = c.Update(bson.M{"name": "hehan"}, bson.M{"$set": bson.M{"age": 12}})
	if err != nil {
		fmt.Println(err)
	}
}

//-----------------official driver-----------
type Trainer struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	City string `json:"city"`
}

func MgInsert() {
	clientOption := options.Client().ApplyURI("mongodb://localhost:27017")
	con, err := mongo.Connect(context.TODO(), clientOption)
	if err != nil {
		fmt.Println(err)
		return
	}
	c := con.Database("helin").Collection("test")
	ash := Trainer{"Ash", 10, "Pallet Town"}
	misty := Trainer{"Misty", 10, "Cerulean City"}
	brock := Trainer{"Brock", 15, "Pewter City"}
	data := []interface{}{ash, misty, brock}
	res, err := c.InsertMany(context.TODO(), data)
	if err != nil {
		fmt.Println(res.InsertedIDs)
		return
	}

}
func MgFind() {
	clientOption := options.Client().ApplyURI("mongodb://localhost:27017")
	con, err := mongo.Connect(context.TODO(), clientOption)
	if err != nil {
		fmt.Println(err)
		return
	}
	c := con.Database("helin").Collection("test")
	//-------------------find------------------
	findOptions := options.Find()
	findOptions.SetLimit(3)
	var result []Trainer
	r, err := c.Find(context.TODO(), bson.M{"age": bson.M{"$gt": 10}}, findOptions)
	if err != nil {
		fmt.Println(err)
		return
	}
	for r.Next(context.TODO()) {
		var tmp Trainer
		err := r.Decode(&tmp)
		if err != nil {
			fmt.Println(err)
			return
		}
		result = append(result, tmp)
	}
	fmt.Println(result)
}
func MgUpdate() {
	clientOption := options.Client().ApplyURI("mongodb://localhost:27017")
	con, err := mongo.Connect(context.TODO(), clientOption)
	if err != nil {
		fmt.Println(err)
		return
	}
	c := con.Database("helin").Collection("test")
	result, err := c.UpdateOne(context.TODO(), bson.M{"name": "Ash"}, bson.M{"$inc": bson.M{"age": 1}})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(result)
}

//-----------------数据采集-----------------------
type HData struct {
	Title string `json:"title"`
	Pic   string `json:"pic"`
	Year  int    `json:"year"`
	Month int    `json:"month"`
	Day   int    `json:"day"`
	Des   string `json:"des"`
	Lunar string `json:"lunar"`
}

func UrlLink(month int, url string) []string {
	var days int
	var urls []string
	if month == 2 {
		days = 29
	} else if month == 4 && month == 6 && month == 9 && month == 11 {
		days = 30
	} else {
		days = 31
	}
	for i := 1; i <= days; i++ {
		urls = append(urls, fmt.Sprintf(url, month, i))
	}
	return urls
}
func DownloadData(urls []string) []HData {
	var ret []HData
	for _, s := range urls {
		result := new(HData)
		r, _ := http.Get(s)
		defer r.Body.Close()
		json.NewDecoder(r.Body).Decode(result)
		fmt.Println(result)
		ret = append(ret, *result)
	}
	return ret
}
func StoreData(data []HData, collectionName string, conn *mgo.Session) error {
	c := conn.DB("HistoryOfToday").C(collectionName)
	for _, d := range data {
		err := c.Insert(d)
		if err != nil {
			return err
		}
	}
	return nil
}
func WorkFunc(month int, url string, conn *mgo.Session, wg *sync.WaitGroup) {
	defer wg.Done()
	defer conn.Close()
	urls := UrlLink(month, url)   //拼接url
	ressult := DownloadData(urls) //请求拿数据
	var collectionName string
	switch month {
	case 1:
		collectionName = "January"
		break
	case 2:
		collectionName = "February"
		break
	case 3:
		collectionName = "March"
		break
	case 4:
		collectionName = "April"
		break
	case 5:
		collectionName = "May"
		break
	case 6:
		collectionName = "June"
		break
	case 7:
		collectionName = "July"
		break
	case 8:
		collectionName = "August"
		break
	case 9:
		collectionName = "September"
		break
	case 10:
		collectionName = "October"
		break
	case 11:
		collectionName = "November"
		break
	case 12:
		collectionName = "December"
		break
	default:
		panic("err month")
	}
	err := StoreData(ressult, collectionName, conn) //将数据存储到本地
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	var wg sync.WaitGroup
	var conAarry []*mgo.Session
	for i := 0; i < 12; i++ {
		con, err := mgo.Dial("")
		if err != nil {
			return
		}
		conAarry = append(conAarry, con)
	}
	url := "http://api.juheapi.com/japi/toh?v=1.0&month=%d&day=%d&key=c25c742a5f7f7b6444ba30b8a2734dff"
	for i := 1; i <= 12; i++ {
		wg.Add(1)
		go WorkFunc(i, url, conAarry[i-1], &wg)
	}
	wg.Wait()
}
