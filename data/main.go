package main

import (
	"flag"
	"fmt"
	"neo4j-csv/data/script"
	"strings"
	"sync"
)

var (
	model  string
	dir    string
	delete bool
)

const (
	User      = ("user")
	Product   = ("product")
	Purchased = ("purchased")
	LiveIn    = ("livein")
	IP        = ("ip")
	Referee   = ("referee")
)

func main() {

	flag.StringVar(&model, "model", "", "导入类型")
	flag.StringVar(&dir, "dir", "", "文件")
	flag.BoolVar(&delete, "delete", false, "导入前是否清理数据")
	flag.Parse()
	files := strings.Split(dir, ",")

	var wg sync.WaitGroup
	wg.Add(len(files))
	for _, file := range files {
		if strings.HasSuffix(file, ".csv") {
			go func(f string) {
				switch model {
				case User:
					script.NewImportEntityUser().OpData(f, delete)
					wg.Done()
				case Product:
					script.NewImportEntityProduct().OpData(f, delete)
					wg.Done()
				case Purchased:
					script.NewImportRelationshipPurchased().OpData(f, delete)
					wg.Done()
				case IP:
					script.NewImportRelationshipIP().OpData(f, delete)
					wg.Done()
				case LiveIn:
					script.NewImportRelationshipLiveIn().OpData(f, delete)
					wg.Done()
				case Referee:
					script.NewImportRelationshipReferee().OpData(f, delete)
					wg.Done()
				}
				fmt.Println("用户协程", "执行完毕")
			}(file)
		} else {
			wg.Done()
		}
	}
	// 等待所有协程执行完毕
	wg.Wait()
}
