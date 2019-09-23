package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type Node struct {
	value interface{}
	key   string
	next  *Node
	last  *Node
}

type List struct {
	len  int
	head *Node
	tail *Node
}

type LRUCache struct {
	rList List
	rMap  map[string]*Node
	cap   int
}

func (c *LRUCache) Get(key string) interface{} {
	if n, ok := c.rMap[key]; ok {
		c.rList.SetToHead(n)
		return n.value
	}
	return nil
}

//it adapt node in the list and node not in the list
func (l *List) SetToHead(n *Node) {
	if n == l.head {
		return
	}

	//remove n
	last := n.last
	next := n.next
	if last != nil {
		last.next = next
	}
	if next != nil {
		next.last = last
	} else if l.tail == n {
		l.tail = n.last
	}

	//n to head
	lhead := l.head
	if lhead != nil {
		lhead.last = n
		n.next = lhead
	} else {
		n.next = nil
	}
	l.head = n
	n.last = nil

	//set to tail if tail is nil
	if l.tail == nil {
		l.tail = l.head
	}
}

func (c *LRUCache) Put(key string, value interface{}) {
	if c.cap <= 0 {
		return
	}
	if n, ok := c.rMap[key]; ok {
		c.rList.SetToHead(n)
		//set value
		n.value = value
	} else if c.rList.len == c.cap {
		okey := c.rList.tail.key
		delete(c.rMap, okey)
		c.rMap[key] = c.rList.tail
		c.rList.SetToHead(c.rList.tail)
		c.rList.head.value = value
		c.rList.head.key = key
	} else {
		n := Node{value: value, key: key}
		c.rList.SetToHead(&n)
		c.rList.len++
		c.rMap[key] = &n
	}
}

func create() LRUCache {
	return LRUCache{rMap: make(map[string]*Node, 1), cap: 6}
}

func (l List) printAll() {
	for n := l.head; n != nil; n = n.next {
		fmt.Println("node:", n.value)
	}
}

func main() {
	lruc := create()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	t := time.Now()
	i := 0
	ni := 0
	for time.Since(t) < time.Second {
		rn := r.Int()
		rn %= 12
		lruc.Put(strconv.Itoa(rn), strconv.Itoa(rn))
		rn = r.Int()
		rn %= 12
		g := lruc.Get(strconv.Itoa(rn))
		if g == nil {
			ni++
		}
		i++
	}
	fmt.Println(time.Since(t))
	fmt.Println("i:", i)
	fmt.Println("ni:", ni)
	fmt.Println(lruc.Get("1"))
	lruc.Put("1", 1)
	fmt.Println(lruc.Get("1"))
	lruc.Put("2", 2)
	fmt.Println(lruc.Get("2"))
	lruc.Put("1", 3)
	fmt.Println(lruc.Get("1"))
	lruc.Put("4", 4)
	fmt.Println(lruc.Get("4"))
	lruc.Put("5", 5)
	fmt.Println(lruc.Get("5"))
	fmt.Println(lruc.Get("2"))
	a := lruc.Get("2")

	fmt.Println("=============")
	fmt.Println(lruc.cap)
	fmt.Println(lruc.rList.len)
	fmt.Println(lruc.rList.head)
	fmt.Println(lruc.rList.tail)
	lruc.rList.printAll()
}
