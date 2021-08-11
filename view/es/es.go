package es

import (
	"context"
	"encoding/json"
	"errors"

	"time"

	"github.com/olivere/elastic/v7"
	"github.com/xxjwxc/public/mylog"
	//"github.com/olivere/elastic"
	//"gopkg.in/olivere/elastic.v5"
)

// MyElastic elastic tools (操作工具)
type MyElastic struct {
	client *elastic.Client
	ops    options
}

// New 新建es
func New(opts ...Option) (*MyElastic, error) {
	es := &MyElastic{
		ops: options{ // default option
			retries: 2,
			timeout: 4 * time.Second,
			ctx:     context.Background(),
		},
	}
	for _, o := range opts {
		o.apply(&es.ops)
	}

	err := es.newClient()
	if err != nil {
		return nil, err
	}
	es.CreateIndex(mapping)

	return es, nil
}

//
// func newClient(url string) MyElastic {
// 	var es MyElastic
// 	es.Ctx = context.Background()
// 	es.Client, es.Err = elastic.NewClient(elastic.SetURL(url))
// 	if es.Err != nil {
// 		log.Println(es.Err)
// 		//mylog.Error(es.Err)
// 		//panic(es.Err)
// 	}

// 	return es
// }

// GetClient 获取es client
func (es *MyElastic) GetClient() *elastic.Client {
	return es.client
}

// GetClient 动态设置option
func (es *MyElastic) WithOption(opts ...Option) *MyElastic {
	for _, o := range opts {
		o.apply(&es.ops)
	}

	return es
}

// CreateIndex 创建索引（相当于数据库）.mapping 如果为空("")则表示不创建模型
func (es *MyElastic) CreateIndex(mapping string) error {
	exists, err := es.client.IndexExists(es.ops.indexName).Do(es.ops.ctx)
	if err != nil {
		mylog.Error(err)
		return err
	}

	if !exists {
		var re *elastic.IndicesCreateResult
		if len(mapping) == 0 {
			re, err = es.client.CreateIndex(es.ops.indexName).Do(es.ops.ctx)
		} else {
			re, err = es.client.CreateIndex(es.ops.indexName).BodyString(mapping).Do(es.ops.ctx)
		}

		if err != nil {
			mylog.Error(err)
			return err
		}

		if re.Acknowledged {
			return nil
		}
	}

	return nil
}

// SortQuery 排序查询,返回json数据集合
func (es *MyElastic) SortQuery(builder []elastic.Sorter, query []elastic.Query) ([]string, error) {
	searchResult := es.client.Search().Index(es.ops.indexName)

	if len(builder) > 0 {
		for _, v := range builder {
			searchResult = searchResult.SortBy(v)
		}
	}
	if len(query) > 0 {
		for _, v := range query {
			searchResult = searchResult.Query(v)
		}
	}
	es_result, err := searchResult.Do(es.ops.ctx) // execute
	if err != nil {
		mylog.Error(err)
		return nil, err
	}
	//log.Println("Found a total of %d entity\n", es_result.TotalHits())

	if len(es_result.Hits.Hits) > 0 {
		var result []string
		//log.Println("Found a total of %d entity\n", searchResult.Hits.TotalHits)
		for _, hit := range es_result.Hits.Hits {
			result = append(result, string(hit.Source))
		}
		return result, nil
	}

	return nil, nil
}

// SortQueryReturnHits  排序查询  返回原始Hit(builder：排序 agg：聚合 类似group_by sum,query：查询)
func (es *MyElastic) SortQueryReturnHits(from, size int, builder []elastic.Sorter, query []elastic.Query) ([]*elastic.SearchHit, error) {
	searchResult := es.client.Search().Index(es.ops.indexName)
	if len(builder) > 0 {
		for _, v := range builder {
			searchResult = searchResult.SortBy(v)
		}
	}
	if len(query) > 0 {
		for _, v := range query {
			searchResult = searchResult.Query(v)
		}
	}
	if size > 0 {
		searchResult = searchResult.From(from)
		searchResult = searchResult.Size(size)
	}
	esResult, err := searchResult.Do(es.ops.ctx) // execute
	if err != nil {
		mylog.Error(err)
		return nil, err
	}

	//	log.Println("wwwwww", es_result.Aggregations)
	if len(esResult.Hits.Hits) > 0 {
		return esResult.Hits.Hits, nil
	}

	return []*elastic.SearchHit{}, nil
}

// Add 添加记录,覆盖添加
func (es *MyElastic) Add(data interface{}, id ...string) (err error) {
	// Index a tweet (using JSON serialization)
	if len(id) > 0 {
		_, err = es.client.Index().
			Index(es.ops.indexName).
			Type(es.ops.typeName).
			Id(id[0]).
			BodyJson(data).
			Do(es.ops.ctx)
	} else {
		_, err = es.client.Index().
			Index(es.ops.indexName).
			Type(es.ops.typeName).
			BodyJson(data).
			Do(es.ops.ctx)
	}

	if err != nil {
		mylog.Error(err)
		return err
	}
	_, err = es.client.Flush().Index(es.ops.indexName).Do(es.ops.ctx)
	if err != nil {
		mylog.Error(err)
		return err
	}
	return nil
}

// BulkAdd 批量新增
func (es *MyElastic) BulkAdd(data []interface{}, id ...string) (err error) {
	// Index a tweet (using JSON serialization)
	bulkRequest := es.client.Bulk()
	if len(id) > 0 {
		for _, doc := range data {
			esRequest := elastic.NewBulkIndexRequest().
				Index(es.ops.indexName).
				Type(es.ops.typeName).
				Id(id[0]).Doc(doc)
			bulkRequest = bulkRequest.Add(esRequest)
		}
	} else {
		for _, doc := range data {
			esRequest := elastic.NewBulkIndexRequest().
				Index(es.ops.indexName).
				Type(es.ops.typeName).
				Doc(doc)
			bulkRequest = bulkRequest.Add(esRequest)
		}
	}

	_, err = bulkRequest.Do(es.ops.ctx)
	if err != nil {
		mylog.Error(err)
		return err
	}
	_, err = es.client.Flush().Index(es.ops.indexName).Do(es.ops.ctx)
	if err != nil {
		mylog.Error(err)
		return err
	}
	return nil

}

// SearchMap 添加记录,覆盖添加 index_name type_name query interface{} //查询条件 out *[]Param //查询结果
func (es *MyElastic) SearchMap(query interface{}) (out []map[string]interface{}, err error) {
	es_search := es.client.Search()
	if len(es.ops.typeName) > 0 {
		es_search = es_search.Type(es.ops.typeName)
	}
	if len(es.ops.indexName) > 0 {
		es_search = es_search.Index(es.ops.indexName)
	}

	esResult, err := es_search.Source(query).Do(es.ops.ctx)
	if err != nil {
		mylog.Error(err)
		return nil, err
	}
	if esResult.Hits == nil {
		return nil, errors.New("expected SearchResult.Hits != nil; got nil")
	}
	for _, hit := range esResult.Hits.Hits {
		tmp := make(map[string]interface{})
		err := json.Unmarshal(hit.Source, &tmp)
		if err != nil {
			mylog.Error(err)
		} else {
			out = append(out, tmp)
		}
	}

	return out, nil
}

// Search 自定义搜索结果
func (es *MyElastic) Search(query interface{}, f func(e []byte) error) (total int64, err error) {

	es_search := es.client.Search()
	if len(es.ops.typeName) > 0 {
		es_search = es_search.Type(es.ops.typeName)
	}
	if len(es.ops.indexName) > 0 {
		es_search = es_search.Index(es.ops.indexName)
	}

	esResult, err := es_search.Source(query).Do(es.ops.ctx)
	if err != nil {
		mylog.Error(err)
		return 0, err
	}
	if esResult.Hits == nil {
		return 0, errors.New("expected SearchResult.Hits != nil; got nil")
	}

	for _, hit := range esResult.Hits.Hits {
		f(hit.Source) //如果 inmprt github.com/olivere/elastic 需要去掉 *
	}

	return esResult.TotalHits(), nil
}
