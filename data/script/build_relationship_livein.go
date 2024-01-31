package script

import (
	"fmt"
	"os"
	"time"

	"github.com/gocarina/gocsv"
)

type RelationshipLiveIn struct {
	BaseEntity
	ID      string `csv:"ID"`
	Address string `csv:"Address"`
}

func NewImportRelationshipLiveIn() *RelationshipLiveIn {
	return &RelationshipLiveIn{}
}

func (self *RelationshipLiveIn) generateCQL(path string) error {
	start := time.Now()
	tasks := make([]chan error, 0)
	// 使用 encoding/csv 读取 CSV 文件
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	models := []*RelationshipLiveIn{}

	if err := gocsv.UnmarshalFile(f, &models); err != nil { // Load clients from file
		panic(err)
	}

	cqlQueries := []string{}
	for _, row := range models {
		// 构建 CQL 语句
		cql := fmt.Sprintf(`MATCH(p:User{userID:"%s"})
		MERGE (q:Address{address:"%s"}) 
		MERGE (p)-[rel:LIVE_IN]->(q)
				`, row.ID, row.Address)
		fmt.Println(cql)
		cqlQueries = append(cqlQueries, cql)

		if len(cqlQueries)%1000 == 0 {
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

func (this *RelationshipLiveIn) OpData(path string, isDeleted bool) error {
	if isDeleted {
		err := this.clearData("MATCH (p:User)-[r:LOGIN]-(q:IP) DELETE r")
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
