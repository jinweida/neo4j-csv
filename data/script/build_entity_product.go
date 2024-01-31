package script

import (
	"fmt"
	"os"
	"time"

	"github.com/gocarina/gocsv"
)

type EntityProduct struct {
	BaseEntity
	PRODCT_CODE  string `csv:"PRODUCT_CODE"`
	PRODUCT_NAME string `csv:"PRODUCT_NAME"`
	AMOUNT       string `csv:"AMOUNT"`
}

func NewImportEntityProduct() *EntityProduct {
	return &EntityProduct{}
}

func (self *EntityProduct) generateCQL(path string) error {
	start := time.Now()
	tasks := make([]chan error, 0)
	// 使用 encoding/csv 读取 CSV 文件
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	model := []*EntityProduct{}

	if err := gocsv.UnmarshalFile(f, &model); err != nil { // Load clients from file
		panic(err)
	}

	cqlQueries := []string{}
	for _, row := range model {

		// 构建 CQL 语句
		cql := fmt.Sprintf(`MERGE (p0:Product {productID: "%s"})
			ON CREATE SET p0.productName="%s",p0.amount="%s"
			ON MATCH SET p0.productName="%s",p0.amounnt="%s"`,
			row.PRODCT_CODE,
			row.PRODUCT_NAME, row.AMOUNT,
			row.PRODUCT_NAME, row.AMOUNT)

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

func (this *EntityProduct) OpData(path string, isDeleted bool) error {
	if isDeleted {
		err := this.clearData("MATCH (n:Product) DELETE n")
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
