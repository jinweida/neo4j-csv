package data

import "fclink.cn/neo4j-tool/utils"

type BaseEntity struct {
}

// 定义公共方法
func (this *BaseEntity) clearData(cql string) error {
	//cql := "MATCH (n:User) DELETE n"
	utils.NewNeo4jProvider().ExecQuery(cql, nil)
	return nil
}

func (this *BaseEntity) executeCQLInGoroutine(cqlQueries []string) chan error {
	errChan := make(chan error, 1)
	go func() {
		defer close(errChan)

		utils.NewNeo4jProvider().BatchExecWrite(cqlQueries)

		errChan <- nil
	}()
	return errChan
}
func (this *BaseEntity) generateCQL(path string) error {
	return nil
}
