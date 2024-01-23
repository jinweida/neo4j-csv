package data

import (
	"fmt"
	"github.com/gocarina/gocsv"
	"os"
	"time"
)

type RelationshipVisited struct {
	BaseEntity
	MODULE       string `csv:"MODULE"`
	USER_ID      string `csv:"USER_ID"`
	VISITDATE    string `csv:"VISITDATE"`
	LOGIN_IP     string `csv:"LOGIN_IP"`
	PLATFORM     string `csv:"PLATFORM"`
	DEVICE_MODAL string `csv:"DEVICE_MODAL"`
}

func NewImportRelationshipVisited() *RelationshipPurchased {
	return &RelationshipPurchased{}
}

func (self *RelationshipVisited) generateCQL(path string) error {
	start := time.Now()
	tasks := make([]chan error, 0)
	// 使用 encoding/csv 读取 CSV 文件
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	visited := []*RelationshipVisited{}

	if err := gocsv.UnmarshalFile(f, &visited); err != nil { // Load clients from file
		panic(err)
	}

	cqlQueries := []string{}
	for _, row := range visited {
		// 构建 CQL 语句
		cql := fmt.Sprintf(` MATCH(p:User),(q:Product) WHERE p.userID='%s' and q.productID='%s' CREATE (p)-[rel:PURCHASED{purchaseDate:'%s',amount:%s}]->(q) `,
			row.USER_ID, row.USER_ID, row.USER_ID, row.USER_ID)
		cqlQueries = append(cqlQueries, cql)

		if len(cqlQueries)%100 == 0 {
			// 使用协程异步执行 CQL 查询
			tasks = append(tasks, self.executeCQLInGoroutine(cqlQueries))
			cqlQueries = []string{}
		}
	}

	tasks = append(tasks, self.executeCQLInGoroutine(cqlQueries))
	// 等待所有协程完成
	for _, task := range tasks {
		err := <-task
		if err != nil {
			return err
		}
	}
	//
	elapsed := time.Since(start)
	fmt.Printf("task=%d, 程序运行时间为 %s\n", len(tasks), elapsed)
	return nil
}

func (this *RelationshipVisited) OpData(path string, isDeleted bool) error {
	if isDeleted {
		err := this.clearData("MATCH (p:User)-[r:VISITED]-(q:Module) DELETE r")
		if err != nil {
			return err
		}
	}
	err := this.generateCQL(path)
	if err != nil {
		// 处理错误
		// ...
	}

	return nil
}
