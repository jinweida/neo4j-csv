package script

import (
	"fmt"
	"os"
	"time"

	"github.com/gocarina/gocsv"
)

type EntityUser struct {
	BaseEntity
	ID               string `csv:"ID"`
	NAME             string `csv:"NAME"`
	MOBILE           string `csv:"mobile"`
	SEX              string `csv:"ACCOUNT"`
	GMT_CREATE       string `csv:"GMT_CREATE"`
	USER_GROUP       string `csv:"USER_GROUP"`
	EDUCATION        string `csv:"EDUCATION"`
	ADDRESS          string `csv:"ADDRESS"`
	CERTIFICATE_CODE string `csv:"CERTIFICATE_CODE"`
}

func NewImportEntityUser() *EntityUser {
	return &EntityUser{}
}
func (this *EntityUser) getBirthday(idCard string) string {
	birthDay := ""
	if idCard != "" && len(idCard) == 18 {
		year := idCard[6:10]
		month := idCard[10:12]
		day := idCard[12:14]
		birthDay = year + "-" + month + "-" + day
	}

	return birthDay
}

func (self *EntityUser) generateCQL(path string) error {
	start := time.Now()
	tasks := make([]chan error, 0)
	// 使用 encoding/csv 读取 CSV 文件
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	users := []*EntityUser{}

	if err := gocsv.UnmarshalFile(f, &users); err != nil { // Load clients from file
		panic(err)
	}

	cqlQueries := []string{}
	for _, row := range users {

		// 构建 CQL 语句
		cql := fmt.Sprintf(`CREATE (user%s:User {name: "%s", userID: "%s", gender: "%s",userGroup: "%s",
					mobile: "%s",signUpDate: "%s", education: "%s", birthDate: "%s",address: "%s"})`,
			row.ID, row.NAME, row.ID, row.SEX, row.USER_GROUP, row.MOBILE, row.GMT_CREATE, row.EDUCATION, self.getBirthday(row.CERTIFICATE_CODE),
			row.ADDRESS)

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

func (this *EntityUser) OpData(path string, isDeleted bool) error {
	if isDeleted {
		err := this.clearData("MATCH (n:User) DELETE n")
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
