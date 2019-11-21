package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

var (
	c = flag.Int("c", 0, "请输入请求数量")
	t = flag.Int("t", 0, "请输入线程数量")
	u = flag.String("u", "", "请输入测试地址")
)

type Result struct {
	PreTotal int
	Total    int64
	Success  int64
	Failure  int64
	Error    int64
	Index    int64
}

var result Result

var startTime = time.Now().UnixNano()
var endTime int64
var wg sync.WaitGroup

func main() {

	flag.Parse()

	startTime = time.Now().UnixNano()
	if *c == 0 || *t == 0 || *u == "" {
		flag.PrintDefaults()
		return
	}

	for i := 0; i < *t; i++ {
		wg.Add(1)
		go run(&i, *c)
	}

	wg.Wait()
	endTime = time.Now().UnixNano()
	printInfo()
}

func run(cIndex *int, num int) {
	for i := 0; i < num; i++ {
		atomic.AddInt64(&result.Total, 1)
		request(*u, &i, cIndex)
	}
	defer wg.Done()
}

func request(url string, i, c *int) {
	atomic.AddInt64(&result.Index, 1)

	resp, err := http.Get(url)
	if err != nil {
		atomic.AddInt64(&result.Error, 1)
		log.Println(err, i, c)
	} else if resp.StatusCode != 200 {
		atomic.AddInt64(&result.Failure, 1)
		defer resp.Body.Close()
	} else {
		atomic.AddInt64(&result.Success, 1)
		defer resp.Body.Close()
	}

}

func printInfo() {
	_time := float64(endTime-startTime) / 1e9
	count := float64(result.Total) / _time
	rate := float64(result.Success) / float64(result.Total) * 100.0

	fmt.Println("-----------------------------------------------")
	fmt.Println("应请求总数:    ", (*c)*(*t))
	fmt.Println("实际请求总数:  ", result.Total)
	fmt.Println("成功:          ", result.Success)
	fmt.Println("失败:          ", result.Failure)
	fmt.Println("错误:          ", result.Error)
	//fmt.Println("丢失:          ", result.Total-(result.Success+result.Failure+result.Error))
	fmt.Println("成功率:        ", fmt.Sprintf("%.2f%s", rate, "%"))
	fmt.Println("总耗时:        ", fmt.Sprintf("%.4fs", _time))
	fmt.Println("每秒请求数：   ", fmt.Sprintf("%.0f", count))
}
