package routers

import (
	"fmt"
	"neo4j-csv/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Graph struct{}
type models struct {
	Title      string      `json:"title"`
	Label      string      `json:"label"`
	Entity     string      `json:"entity"`
	Properties interface{} `json:"properties"`
	ID         string      `json:"id"`
}
type relationship struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Label string `json:"label"`
}
type targetEntity struct {
	Nodes         []models       `json:"nodes"`
	Relationships []relationship `json:"relationships"`
}

func (r *Graph) SliceContains(nodes []models, id string) bool {
	for _, v := range nodes {
		if v.ID == id {
			return true
		}
	}
	return false
}
func (r *Graph) SliceContainsRelation(rels []relationship, sid, eid string) bool {
	for _, v := range rels {
		if v.From == sid && v.To == eid {
			return true
		}
	}
	return false
}

func (r *Graph) nodesTextProperties(key string) string {
	nodesTextProperties := map[string]string{
		"User":    "name",
		"Product": "productName",
		"Module":  "moduleName",
		"Address": "address",
		"IP":      "ip",
	}

	if textProperty, ok := nodesTextProperties[key]; ok {
		return textProperty
	} else {
		return key
	}
}
func (g *Graph) TargetAllNeighbours(c *gin.Context) {
	id := c.DefaultQuery("id", "")
	hasID := c.DefaultQuery("has_id", "")
	rel := c.DefaultQuery("rel", "")
	limit := c.DefaultQuery("limit", "25")
	results := []string{}
	for _, item := range strings.Split(hasID, ",") {
		result := strings.Split(item, ":")
		results = append(results, string(result[2]))
	}
	cql := fmt.Sprintf(`
		MATCH (a) WHERE id(a) = $id
		WITH a, size([(a)-[$rel]-() | 1]) AS allNeighboursCount
		MATCH p = (a)-[$rel]-(o) WHERE NOT id(o) IN $results
		RETURN p
		ORDER BY id(o)
		LIMIT $limit`, id, rel, rel, results, limit,
	)
	fmt.Println(cql)
	paramter := make(map[string]any)
	result, err := utils.NewNeo4jProvider().ExecQuery(cql, paramter)
	if err != nil {
		// 如果执行查询有错误，返回错误响应
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	data, err := g.getRelation(result)

	c.JSON(http.StatusOK, data)

}

// @Summary 目标图谱搜索邻居节点
// @Description 目标图谱搜索邻居节点
// @Accept  json
// @Produce json
// @Param	id query string false "当前节点"
// @Param	has_id query string false "已经链接的关系"
// @Param	rel query string false "搜索的关系类型"
// @Success 200 {object} targetEntity
// @Graph /target_graph [get]
func (r *Graph) Targetgraph(c *gin.Context) {

	rel := c.DefaultQuery("rel", "")
	start := c.DefaultQuery("start", "1")
	end := c.DefaultQuery("end", "1")
	field := c.DefaultQuery("field", "1")
	keyword := c.DefaultQuery("keyword", "1")
	limit := c.DefaultQuery("limit", "25")

	rels := fmt.Sprintf("*%s..%s", start, end)
	if rel != "" {
		rels = fmt.Sprintf(":%s*%s..%s", rel, start, end)
	}
	cql := fmt.Sprintf("MATCH p=(left)-[%s]-(right) RETURN p LIMIT %s", rels, limit)
	if field != "" && keyword != "" {
		if field == "name" {
			field = "userID"
		}
		if rel != "" {
			rel = ":" + rel
		}
		if field == "userID" || field == "mobile" {
			condition := fmt.Sprintf("Where s0.%s CONTAINS '%s'", field, keyword)
			cql = fmt.Sprintf(`
                MATCH (s0:User %s)-[r1%s]->(s1)
                WITH s0, collect(s1)[..%s] as s1,collect(distinct r1) as r1
                UNWIND s1 as s2 
                MATCH (s2)-[r2%s]-(s3:User) 
                WITH s0,s2,collect(s3)[..%s] as s3,r1,collect(distinct r2) as r2
                UNWIND s3 as s4
                RETURN [s0,s2,s4] as p,r1,r2
            `, condition, rel, start, rel, end)
		}
	}

	fmt.Println(cql)
	paramter := make(map[string]any)
	result, err := utils.NewNeo4jProvider().ExecQuery(cql, paramter)

	if err != nil {
		// 如果执行查询有错误，返回错误响应
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var data interface{}
	if strings.Contains(cql, "UNWIND") {
		data, _ = r.getNodeAndRelation(result)
	} else {
		data, _ = r.getRelation(result)
	}

	c.JSON(http.StatusOK, data)
}

func (r *Graph) getRelation(result *neo4j.EagerResult) (*targetEntity, error) {
	nodes := []models{}
	relationships := []relationship{}

	for _, record := range result.Records {

		// entity, _ := record.Get("p")
		// for _, v := range entity.([]interface{}) {
		// 	node := v.(neo4j.Node)

		// 	title := r.nodesTextProperties(node.Labels[0])
		// 	model := models{
		// 		Properties: node.Props,
		// 		Title:      node.Props[title].(string),
		// 		ID:         node.ElementId,
		// 		Entity:     node.Labels[0],
		// 		Label:      node.Labels[0],
		// 	}
		// 	if !r.SliceContains(nodes, model.ID) {
		// 		nodes = append(nodes, model)
		// 	}
		// }

		r1, _ := record.Get("p")
		// fmt.Printf("%+v", reflect.ValueOf(r1))
		path := r1.(neo4j.Path)
		for _, node := range path.Nodes {
			title := r.nodesTextProperties(node.Labels[0])
			model := models{
				Properties: node.Props,
				Title:      node.Props[title].(string),
				ID:         node.ElementId,
				Entity:     node.Labels[0],
				Label:      node.Labels[0],
			}
			if !r.SliceContains(nodes, model.ID) {
				nodes = append(nodes, model)
			}
		}
		for _, rel := range path.Relationships {

			sid := rel.StartElementId
			eid := rel.EndElementId

			if !r.SliceContainsRelation(relationships, sid, eid) {
				relationships = append(relationships, relationship{
					From: sid, To: eid, Label: rel.Type,
				})
			}
		}

	}
	graph := &targetEntity{
		Nodes: nodes, Relationships: relationships,
	}
	return graph, nil
}

func (r *Graph) getNodeAndRelation(result *neo4j.EagerResult) (*targetEntity, error) {
	nodes := []models{}
	relationships := []relationship{}

	for _, record := range result.Records {

		entity, _ := record.Get("p")
		for _, v := range entity.([]interface{}) {
			node := v.(neo4j.Node)

			title := r.nodesTextProperties(node.Labels[0])
			model := models{
				Properties: node.Props,
				Title:      node.Props[title].(string),
				ID:         node.ElementId,
				Entity:     node.Labels[0],
				Label:      node.Labels[0],
			}
			if !r.SliceContains(nodes, model.ID) {
				nodes = append(nodes, model)
			}
		}

		r1, _ := record.Get("r1")
		for _, v := range r1.([]interface{}) {
			rel := v.(neo4j.Relationship)
			sid := rel.StartElementId
			eid := rel.EndElementId

			if !r.SliceContainsRelation(relationships, sid, eid) {
				relationships = append(relationships, relationship{
					From: sid, To: eid, Label: rel.Type,
				})
			}

		}
		r2, _ := record.Get("r1")
		for _, v := range r2.([]interface{}) {
			rel := v.(neo4j.Relationship)
			sid := rel.StartElementId
			eid := rel.EndElementId

			if !r.SliceContainsRelation(relationships, sid, eid) {
				relationships = append(relationships, relationship{
					From: sid, To: eid, Label: rel.Type,
				})
			}

		}

	}
	graph := &targetEntity{
		Nodes: nodes, Relationships: relationships,
	}
	return graph, nil
}
