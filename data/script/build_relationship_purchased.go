package script

import (
	"fmt"
	"os"
	"time"

	"github.com/gocarina/gocsv"
)

type RelationshipPurchased struct {
	BaseEntity
	PRODUCT_ID string `csv:"PRODUCT_ID"`
	ACCOUNT    string `csv:"ACCOUNT"`
	UNIT       string `csv:"UNIT"`
}

func NewImportRelationshipPurchased() *RelationshipPurchased {
	return &RelationshipPurchased{}
}

func (self *RelationshipPurchased) generateCQL(path string) error {
	start := time.Now()
	tasks := make([]chan error, 0)
	// 使用 encoding/csv 读取 CSV 文件
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	purchaseds := []*RelationshipPurchased{}

	if err := gocsv.UnmarshalFile(f, &purchaseds); err != nil { // Load clients from file
		panic(err)
	}

	cqlQueries := []string{}
	for _, row := range purchaseds {
		// 构建 CQL 语句
		cql := fmt.Sprintf(`MATCH(p:User{userID:'%s'}),(q:Product{productID:'%s'}) 
        	MERGE (p)-[rel:PURCHASED{amount:%s}]->(q)
		`, row.ACCOUNT, row.PRODUCT_ID, row.UNIT)
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

func (this *RelationshipPurchased) OpData(path string, isDeleted bool) error {
	if isDeleted {
		err := this.clearData("MATCH (:User)-[r:PURCHASED]->(:Product) DELETE r")
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
