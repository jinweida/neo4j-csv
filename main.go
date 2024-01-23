package main

import (
	"flag"
	"fmt"
	"strings"
	"sync"

	"fclink.cn/neo4j-tool/data"
)

var (
	model  string
	dir    string
	delete bool
)

func main() {
	flag.StringVar(&model, "model", "", "导入类型")
	flag.StringVar(&dir, "dir", "", "文件")
	flag.BoolVar(&delete, "delete", false, "导入前是否清理数据")
	flag.Parse()
	files := strings.Split(dir, ",")

	fmt.Println(files)
	var wg sync.WaitGroup
	wg.Add(len(files))
	for _, file := range files {
		fmt.Println(file)
		if strings.HasSuffix(file, ".csv") {
			go func() {
				if model == "user" {
					// 执行任务
					data.NewImportEntityUser().OpData(file, delete)
					wg.Done()
				}
				if model == "purchased" {
					// 执行任务
					data.NewImportRelationshipPurchased().OpData(file, delete)
					wg.Done()
				}
				fmt.Println("用户协程", "执行完毕")
			}()
		} else {
			wg.Done()
		}
	}
	// 等待所有协程执行完毕
	wg.Wait()

	fmt.Println("主线程执行完毕")
}
