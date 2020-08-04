package main

import (
	"fmt"
	"hash/crc32"
	"sort"
	"strconv"
)

type Int64Slice []int64

func (s Int64Slice) Len() int {
	return len(s)
}

func (s Int64Slice) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s Int64Slice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type Hash func(data []byte) uint32

type Map struct {
	hash     Hash
	replicas int               // 复制因子,虚拟节点数
	keys     Int64Slice       // 已排序的节点哈希切片
	hashMap  map[int64]string // 节点哈希和KEY的map，键是哈希值，值是节点Key
}

func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int64]string),
	}
	// 默认使用CRC32算法
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

func (m *Map) IsEmpty() bool {
	return len(m.keys) == 0
}

// Add 方法用来添加缓存节点，参数为节点key，比如使用IP
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		// 结合复制因子计算所有虚拟节点的hash值，并存入m.keys中，同时在m.hashMap中保存哈希值和key的映射
		for i := 0; i < m.replicas; i++ {
			hash := int64(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
	// 对所有虚拟节点的哈希值进行排序，方便之后进行二分查找
	sort.Sort(m.keys)
}

// Get 方法根据给定的对象获取最靠近它的那个节点key
func (m *Map) Get(key string) string {
	if m.IsEmpty() {
		return ""
	}

	hash := m.hash([]byte(key))

	// 通过二分查找获取最优节点，第一个节点hash值大于对象hash值的就是最优节点
	idx := sort.Search(len(m.keys), func(i int) bool { return m.keys[i] >= int64(hash) })

	// 如果查找结果大于节点哈希数组的最大索引，表示此时该对象哈希值位于最后一个节点之后，那么放入第一个节点中
	if idx == len(m.keys) {
		idx = 0
	}

	return m.hashMap[m.keys[idx]]
}

//删除一个节点，并删除其复制节点
func (m *Map) Delete (key string){
	for i := 0; i < m.replicas; i++ {
		hash := int64(m.hash([]byte(strconv.Itoa(i) + key)))
		for index, v := range m.keys{
			if v == hash{
				m.keys[index] = -1
			}
		}
		delete(m.hashMap,hash)
	}
	tmp := make(Int64Slice,0)
	for _ ,value := range m.keys{
		if value != -1{
			tmp = append(tmp,value)
		}
	}
	m.keys = tmp
}



func main(){
	m := New(50,nil)
	m.Add("l","x","m","o")

	for i:=0; i<30; i++{
		fmt.Print(m.Get(strconv.Itoa(i))," ")
	}
	fmt.Println()

	m.Delete("m")
	for i:=0; i<30; i++{
		fmt.Print(m.Get(strconv.Itoa(i))," ")
	}
	fmt.Println()
}
