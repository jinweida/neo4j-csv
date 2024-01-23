package utils

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

var instance *Neo4jProvider

type Neo4jProvider struct {
	provider neo4j.DriverWithContext
	ctx      context.Context
}

func (this *Neo4jProvider) init() {
	ctx := context.Background()
	// URI examples: "neo4j://localhost", "neo4j+s://xxx.databases.neo4j.io"
	dbUri := "bolt://10.18.23.161:8687"
	dbUser := "neo4j"
	dbPassword := "fclink.123"
	driver, err := neo4j.NewDriverWithContext(
		dbUri,
		neo4j.BasicAuth(dbUser, dbPassword, ""))

	this.provider = driver
	this.ctx = ctx
	//defer driver.Close(ctx)

	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		panic(err)
	}

}

func NewNeo4jProvider() *Neo4jProvider {
	if instance == nil {
		instance = &Neo4jProvider{}
		instance.init()
	}
	return instance
}

func (this *Neo4jProvider) ExecQuery(cql string, paramter map[string]any) *neo4j.EagerResult {
	result, err := neo4j.ExecuteQuery(this.ctx, this.provider, cql,
		paramter, neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"),
	)
	if err != nil {
		panic(err)
	}
	return result
}
func (self *Neo4jProvider) BatchExecWrite(cqls []string) error {
	session := self.provider.NewSession(self.ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})

	defer session.Close(self.ctx)
	for _, cql := range cqls {

		session.ExecuteWrite(self.ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			result, _ := tx.Run(self.ctx, cql, nil)
			return result, nil
		})
	}
	return nil
}
