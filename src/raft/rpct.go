// package main

// import (
// 	crand "crypto/rand"
// 	"encoding/base64"
// 	"fmt"
// 	"labrpc"
// 	"sync"
// 	"time"
// )

// // 这个例子学到的知识：
// // 1. 在6.824给的例子中RPC的调用框架顺序
// // 2. RPC中不能修改结构体的值，只能读，返回结果后自己进行修改

// var group Group

// func randstring(n int) string {
// 	b := make([]byte, 2*n)
// 	crand.Read(b)
// 	s := base64.URLEncoding.EncodeToString(b)
// 	return s[0:n]
// }

// type Node struct {
// 	peers []*labrpc.ClientEnd

// 	value int
// 	mu    sync.Mutex
// }

// type Group struct {
// 	net     *labrpc.Network
// 	endname [][]string
// 	node    []Node
// }

// type RequestArgs struct {
// 	IncreaseNumber int
// }

// type ReplyArgs struct {
// 	Status bool
// }

// func (nd *Node) Add(requestargs *RequestArgs, replyargs *ReplyArgs) {
// 	nd.mu.Lock()
// 	defer nd.mu.Unlock()

// 	fmt.Println("in Add RPC")
// 	// fmt.Println(requestargs.IncreaseNumber)
// 	nd.value = nd.value + requestargs.IncreaseNumber
// 	replyargs.Status = true
// 	// RPC无法改变结构体中的数据
// 	// {[0xc000068100 0xc000068120 0xc000068140] 3 {1 0}}
// 	// {[0xc000068100 0xc000068120 0xc000068140] 0 {0 0}}
// 	// fmt.Println("address: ")
// 	// fmt.Println(nd)
// 	// fmt.Print(group.node[0])
// 	return
// }

// func (nd *Node) sendAdd(index int, requestargs *RequestArgs, replyargs *ReplyArgs) bool {
// 	fmt.Println("In send")
// 	fmt.Println(requestargs.IncreaseNumber)
// 	ok := nd.peers[index].Call("Node.Add", requestargs, replyargs)
// 	// fmt.Println(replyargs.Status)
// 	// fmt.Println("prepare to out send")
// 	return ok
// }

// func main() {

// 	//fmt.Println("233")
// 	//node := make([]Node, 3)

// 	group.endname = make([][]string, 3)
// 	group.node = make([]Node, 3)
// 	group.net = labrpc.MakeNetwork()
// 	for i := 0; i < 3; i++ {
// 		group.endname[i] = make([]string, 3)
// 	}
// 	for i := 0; i < 3; i++ {
// 		for j := 0; j < 3; j++ {
// 			group.endname[i][j] = randstring(20)
// 			fmt.Println(group.endname[i][j])
// 		}

// 	}

// 	for i := 0; i < 3; i++ {
// 		node := &Node{}
// 		var ends []*labrpc.ClientEnd
// 		ends = make([]*labrpc.ClientEnd, 3)
// 		for j := 0; j < 3; j++ {
// 			ends[j] = group.net.MakeEnd(group.endname[i][j])
// 			group.net.Connect(group.endname[i][j], j)
// 		}
// 		//添加服务
// 		node.peers = ends
// 		group.node[i] = *node
// 		svc := labrpc.MakeService(node)
// 		srv := labrpc.MakeServer()
// 		srv.AddService(svc)
// 		group.net.AddServer(i, srv)
// 	}

// 	//connect everyone
// 	for i := 0; i < 3; i++ {
// 		for j := 0; j < 3; j++ {
// 			endname := group.endname[i][j]
// 			group.net.Enable(endname, true)
// 		}
// 	}

// 	for i := 0; i < 3; i++ {
// 		go func(index int) {
// 			var arg RequestArgs
// 			arg.IncreaseNumber = 3
// 			var reply ReplyArgs
// 			fmt.Println(group.node[index].sendAdd(index, &arg, &reply))
// 		}(i)
// 	}
// 	//go fmt.Println(node[1].sendAdd(1, &arg, &reply))

// 	//停止一段时间，然后将相关的结果输出
// 	time.Sleep(time.Millisecond * 3000)
// 	for i := 0; i < 3; i++ {
// 		fmt.Println(group.node[i].value, group.node[i])
// 	}

// 	return
// }

package raft
