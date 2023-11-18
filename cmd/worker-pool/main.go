package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"worker-pool/internal/http"
)

type TaskParams struct {
	c  chan int
	wc int
}

type Response struct {
	Comments *[]Comment
}

type Comment struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

func main() {

	var T, N, M int
	flag.IntVar(&T, "T", 3, "number of tasks")
	flag.IntVar(&N, "N", 10, "number of threads")
	flag.IntVar(&M, "M", 4, "max threads per task")
	flag.Parse()

	//считаем сколько мы можем запустить воркеров для каждой залачи
	//с учетом максимального количества потоков
	wc := 0
	var taskInfo = make([]TaskParams, 0, T)
	for i := 0; i < T; i++ {
		ch := make(chan int, 1000)

		wt := 0
		if wc+M <= N {
			wc += M
			wt = M
		} else {
			wt = N - wc
			wc += wt
		}

		tp := TaskParams{ch, wt}
		taskInfo = append(taskInfo, tp)
	}

	wg := sync.WaitGroup{}
	for t, info := range taskInfo {
		t := t

		//заполняем для каждой задачи свой канал рандомными данными
		wg.Add(1)
		go func(ch chan<- int) {
			defer wg.Done()
			for i := 0; i < 15; i++ {
				ch <- rand.Intn(300) + 1
			}
			close(ch)
		}(info.c)

		//запускаем для каждой задачи n воркеров для обработки данных из канала задачи
		wg.Add(info.wc)
		for i := 0; i < info.wc; i++ {
			i := i

			//запуск воркера
			go func(ch <-chan int) {
				defer wg.Done()

				//имитируем работу воркера
				for {
					ss := rand.Intn(5) + 1
					time.Sleep(time.Duration(ss) * time.Second)
					limit, ok := <-ch
					if !ok {
						return
					}

					url := "https://dummyjson.com/comments?limit=" + strconv.Itoa(limit)
					r, err := http.Get(url)
					if err != nil {
						log.Fatalln(err)
					}

					resp := Response{}
					err = json.Unmarshal(r, &resp)
					if err != nil {
						log.Fatal(err)
					}

					fmt.Printf("Task: %d; Worker: %d; Response len items: %d \n", t+1, i+1, len(*resp.Comments))
				}

			}(info.c)
		}
	}

	wg.Wait()

}
