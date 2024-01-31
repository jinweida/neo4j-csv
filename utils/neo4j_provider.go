package utils

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

var instance *Neo4jProvider

type Neo4jProvider struct {
	provider neo4j.DriverWithContext
	ctx      context.Context
}

func init() {
	ctx := context.Background()
	err := godotenv.Load() // look for .env file in current directory
	if err != nil {
		panic(err)
	}
	// URI examples: "neo4j://localhost", "neo4j+s://xxx.databases.neo4j.io"
	dbUri := os.Getenv("NEO4J_URI")
	dbUser := os.Getenv("NEO4J_USER")
	dbPassword := os.Getenv("NEO4J_PASSWORD")
	driver, err := neo4j.NewDriverWithContext(
		dbUri,
		neo4j.BasicAuth(dbUser, dbPassword, ""))

	// defer driver.Close(ctx)
	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	instance = &Neo4jProvider{
		provider: driver,
		ctx:      ctx,
	}

}

func NewNeo4jProvider() *Neo4jProvider {
	return instance
}

func (this *Neo4jProvider) ExecQuery(cql string, paramter map[string]any) (*neo4j.EagerResult, error) {
	result, err := neo4j.ExecuteQuery(this.ctx, this.provider, cql,
		paramter, neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"),
	)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (this *Neo4jProvider) BatchExecWrite(cqls []string) error {
	session := this.provider.NewSession(this.ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})

	defer session.Close(this.ctx)
	for _, cql := range cqls {

		session.ExecuteWrite(this.ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			result, _ := tx.Run(this.ctx, cql, nil)
			return result, nil
		})
	}
	return nil
}
