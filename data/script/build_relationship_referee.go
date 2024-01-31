package script

import (
	"fmt"
	"os"
	"time"

	"github.com/gocarina/gocsv"
)

type RelationshipReferee struct {
	BaseEntity
	RACCOUNT string `csv:"RACCOUNT"`
	FACCOUNT string `csv:"FACCOUNT"`
	REF      string `csv:"REF"`
}

func NewImportRelationshipReferee() *RelationshipReferee {
	return &RelationshipReferee{}
}

func (self *RelationshipReferee) generateCQL(path string) error {
	start := time.Now()
	tasks := make([]chan error, 0)
	// 使用 encoding/csv 读取 CSV 文件
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	models := []*RelationshipReferee{}

	if err := gocsv.UnmarshalFile(f, &models); err != nil { // Load clients from file
		panic(err)
	}

	cqlQueries := []string{}
	for _, row := range models {
		// 构建 CQL 语句
		cql := fmt.Sprintf(`MATCH (p:User{} WHERE p.userID='%s'),(q:User WHERE q.userID='%s') 
		MERGE (p)-[rel:REFEREE{refdate:"%s"}]-(q)`, row.RACCOUNT, row.FACCOUNT, row.REF)
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

func (this *RelationshipReferee) OpData(path string, isDeleted bool) error {
	if isDeleted {
		err := this.clearData("MATCH (m:User)-[r:REFEREE]->(n:User) DELETE r")
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
