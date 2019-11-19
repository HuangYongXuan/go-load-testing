package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	c = flag.Int("c", 0, "Plz input client quantity")
	t = flag.Int("t", 0, "Plz input times quantity")
	u = flag.String("u", "", "Plz input url")
)

var (
	total   = 0.0
	about   = 0.0
	success = 0.0
	failure = 0.0
	useTime = 0.0
	index   = 0
)

var wg sync.WaitGroup

func run(num int) {

	defer wg.Done()

	no := 0.0
	ok := 0.0

	for i := 0; i < num; i++ {
		start := time.Now()
		resp, err := http.Get(*u)
		index++

		if err != nil {
			no += 1
			log.Printf("%d	%s	请求耗时：%.4fs	Error:%s\n", index, *u, time.Since(start).Seconds(), err.Error())
			continue
		}
		//log.Printf("#%d	%s	请求耗时：%.4fs	HttpCode:%s\n", index, *u, time.Since(start).Seconds(), resp.Status)

		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			no += 1
			continue
		}

		ok += 1
		continue
	}

	success += ok
	failure += no
	total += float64(num)

}

func main() {

	startTime := time.Now().UnixNano()

	flag.Parse()

	if *c == 0 || *t == 0 || *u == "" {
		flag.PrintDefaults()
		return
	}

	for i := 0; i < *c; i++ {
		wg.Add(1)
		go run(*t)
	}

	wg.Wait()
	endTime := time.Now().UnixNano()
	_time := float64(endTime-startTime) / 1e9
	count := total / _time

	fmt.Println("应请求总数:    ", (*c)*(*t))
	fmt.Println("实际请求总数:  ", total)
	fmt.Println("成功:          ", success)
	fmt.Println("失败:          ", failure)
	fmt.Println("成功率:        ", fmt.Sprintf("%.2f", (success/total)*100.0), "%")
	fmt.Println("总耗时:        ", fmt.Sprintf("%.4f", _time), "s")
	fmt.Println("每秒请求数：   ", fmt.Sprintf("%.0f", count))
}
