package data

import (
	"fmt"
	"github.com/gocarina/gocsv"
	"os"
	"time"
)

type EntityUser struct {
	BaseEntity
	ID                   string `csv:"ID"`
	NAME                 string `csv:"NAME"`
	ACCOUNT              string `csv:"ACCOUNT"`
	TYPE                 string `csv:"TYPE"`
	SEX                  string `csv:"ACCOUNT"`
	RISK_LEVEL           string `csv:"RISK_LEVEL"`
	MEMBER_CODE          string `csv:"MEMBER_CODE"`
	TRADE_ACCOUNT        string `csv:"TRADE_ACCOUNT"`
	GMT_CREATE           string `csv:"GMT_CREATE"`
	USER_GROUP           string `csv:"USER_GROUP"`
	EDUCATION            string `csv:"EDUCATION"`
	ADDRESS              string `csv:"ADDRESS"`
	CERTIFICATE_CODE     string `csv:"CERTIFICATE_CODE"`
	MOBILE_PROVINCE      string `csv:"MOBILE_PROVINCE"`
	MOBILE_CITY          string `csv:"MOBILE_CITY"`
	CERTIFICATE_PROVINCE string `csv:"CERTIFICATE_PROVINCE"`
	CERTIFICATE_CITY     string `csv:"CERTIFICATE_CITY"`
}

func NewImportEntityUser() *EntityUser {
	return &EntityUser{}
}
func (this *EntityUser) getBirthday(types string, idCard string) string {
	birthDay := ""
	if types == "personal" && idCard != "" && len(idCard) == 18 {
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
		return err
	}

	users := []*EntityUser{}

	if err := gocsv.UnmarshalFile(f, &users); err != nil { // Load clients from file
		panic(err)
	}

	cqlQueries := []string{}
	for _, row := range users {
		Gender := "男"
		if row.SEX == "1" {
			Gender = "女"
		}
		// 构建 CQL 语句
		cql := fmt.Sprintf(`create (user%s:User {name: "%s", userID: "%s", account: "%s", 
						type: "%s", tradeAccount: "%s", memberCode: "%s", gender: "%s", riskLevel: "%s", 
						userGroup: "%s", signUpDate: "%s", education: "%s", birthDate: "%s", 
						address: "%s", mobile_province: "%s", mobile_city: "%s", certificate_province: "%s", 
						certificate_city: "%s"})`,
			row.ID, row.NAME, row.ID, row.ACCOUNT, row.TYPE, row.TRADE_ACCOUNT, row.MEMBER_CODE, Gender, row.RISK_LEVEL,
			row.USER_GROUP, row.GMT_CREATE, row.EDUCATION, self.getBirthday(row.TYPE, row.CERTIFICATE_CODE),
			row.ADDRESS, row.MOBILE_PROVINCE, row.MOBILE_CITY, row.CERTIFICATE_PROVINCE, row.CERTIFICATE_CITY)

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
