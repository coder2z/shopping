package utils

import (
	"errors"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

//声明新切片类型
type units []uint32

//返回切片长度
func (x units) Len() int {
	return len(x)
}

//比对两个数大小
func (x units) Less(i, j int) bool {
	return x[i] < x[j]
}

//切片中两个值的交换
func (x units) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

type ConsistentHashImp interface {
	Add(string)
	Remove(string)
	Get(string) (string, error)
}

//创建结构体保存一致性hash信息
type ConsistentHash struct {
	//hash环，key为哈希值，值存放节点的信息
	circle map[uint32]string
	//已经排序的节点hash切片
	sortedHashes units
	//虚拟节点个数，用来增加hash的平衡性
	VirtualNode int
	//map 读写锁
	sync.RWMutex
}

//创建一致性hash算法结构体，设置默认节点数量
func NewConsistent(nodeNum int) ConsistentHashImp {
	return &ConsistentHash{
		//初始化变量
		circle: make(map[uint32]string),
		//设置虚拟节点个数
		VirtualNode: nodeNum,
	}
}

//自动生成key值
func (c *ConsistentHash) generateKey(element string, index int) string {
	//副本key生成逻辑
	return element + strconv.Itoa(index)
}

//获取hash位置 计算key 在hash环中对应的位置
func (c *ConsistentHash) hashKey(key string) uint32 {
	//当长度不够填充
	if len(key) < 64 {
		//声明一个数组长度为64
		var tmpList [64]byte
		//拷贝数据到数组中
		copy(tmpList[:], key)
		//使用IEEE 多项式返回数据的CRC-32校验和
		return crc32.ChecksumIEEE(tmpList[:len(key)])
	}
	return crc32.ChecksumIEEE([]byte(key))
}

//更新排序，方便查找 因为后面我们使用的是sort.Search进行查找 sort.Search使用的是二分法进行查找，所以这里需要排序
func (c *ConsistentHash) updateSortedHashes() {
	hashes := c.sortedHashes[:0]
	//判断切片容量，是否过大，如果过大则重置
	if cap(c.sortedHashes)/(c.VirtualNode*4) > len(c.circle) {
		hashes = nil
	}

	//添加hashes
	for k := range c.circle {
		hashes = append(hashes, k)
	}

	//对所有节点hash值进行排序，
	//方便之后进行二分查找
	sort.Sort(hashes)
	//重新赋值
	c.sortedHashes = hashes
}

//向hash环中添加节点
func (c *ConsistentHash) Add(element string) {
	//加锁
	c.Lock()
	//解锁
	defer c.Unlock()
	c.add(element)
}

//添加节点
func (c *ConsistentHash) add(element string) {
	//循环虚拟节点，设置副本
	for i := 0; i < c.VirtualNode; i++ {
		//根据生成的节点添加到hash环中
		c.circle[c.hashKey(c.generateKey(element, i))] = element
	}
	//更新排序
	c.updateSortedHashes()
}

//删除节点
func (c *ConsistentHash) remove(element string) {
	for i := 0; i < c.VirtualNode; i++ {
		delete(c.circle, c.hashKey(c.generateKey(element, i)))
	}
	c.updateSortedHashes()
}

//删除一个节点
func (c *ConsistentHash) Remove(element string) {
	c.Lock()
	defer c.Unlock()
	c.remove(element)
}

//顺时针查找最近的节点
func (c *ConsistentHash) search(key uint32) int {
	//查找算法
	f := func(x int) bool {
		return c.sortedHashes[x] > key
	}
	//使用"二分查找"算法来搜索指定切片满足条件的最小值
	i := sort.Search(len(c.sortedHashes), f)
	//如果超出范围则设置i=0
	if i >= len(c.sortedHashes) {
		i = 0
	}
	return i
}

//根据数据标示获取最近的服务器节点信息
func (c *ConsistentHash) Get(name string) (string, error) {
	//添加锁
	c.RLock()
	//解锁
	defer c.RUnlock()
	//如果为零则返回错误
	if len(c.circle) == 0 {
		return "", errors.New("hash环没有数据")
	}
	//计算hash值
	key := c.hashKey(name)
	i := c.search(key)
	return c.circle[c.sortedHashes[i]], nil
}
