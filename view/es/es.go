package es

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sync"

	"time"
	//"github.com/olivere/elastic"
	"gopkg.in/olivere/elastic.v5"
)

var Url string
var Index string

//
func GetClient() *MyElastic {
	once.Do(func() {
		cli, err := NewClient(4*time.Second, 2)
		if err != nil {
			log.Println("EsClient create error :", err)
			once = sync.Once{} //重入
		} else {
			esc = new(MyElastic)
			esc.Client = cli
			esc.Ctx = context.Background()
			esc.CreateIndex(Index, mapping)
		}
	})

	return esc
}

//
func NewClient(timeout time.Duration, retries int) (*elastic.Client, error) {
	var cli *elastic.Client
	var err error

	httpClient := &http.Client{
		Timeout: timeout,
	}

	for i := 0; i < retry; i++ {
		cli, err = elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(Url), elastic.SetHttpClient(httpClient), elastic.SetMaxRetries(retries))
		if err == nil {
			break
		}
		log.Println("EsClient create failed ", err, ", tried %d time(s)", i+1)
		time.Sleep(100 * time.Millisecond)
	}

	return cli, err
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

/*
创建索引（相当于数据库）
mapping 如果为空("")则表示不创建模型
*/
func (es *MyElastic) CreateIndex(index_name, mapping string) (result bool) {
	es.Err = nil
	exists, err := es.Client.IndexExists(index_name).Do(es.Ctx)
	if err != nil {
		es.Err = err
		log.Println(es.Err)
		return false
	}

	if !exists {
		var re *elastic.IndicesCreateResult
		if len(mapping) == 0 {
			re, es.Err = es.Client.CreateIndex(index_name).Do(es.Ctx)
		} else {
			re, es.Err = es.Client.CreateIndex(index_name).BodyString(mapping).Do(es.Ctx)
		}

		if es.Err != nil {
			log.Println(es.Err)
			return false
		}

		return re.Acknowledged
	}

	return false
}

/*
	排序查询
	返回json数据集合
*/
func (es *MyElastic) SortQuery(index_name string, builder []elastic.Sorter, query []elastic.Query) (bool, []string) {

	searchResult := es.Client.Search().Index(index_name)

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
	_, err := searchResult.Do(es.Ctx) // execute
	if err != nil {
		log.Println(es.Err)
		return false, nil
	}
	//log.Println("Found a total of %d entity\n", es_result.TotalHits())

	// if es_result.Hits.TotalHits > 0 {
	// 	var result []string
	// 	//log.Println("Found a total of %d entity\n", searchResult.Hits.TotalHits)
	// 	for _, hit := range es_result.Hits.Hits {

	// 		result = append(result, string(*hit.Source))

	// 	}
	// 	return true, result
	// } else {
	// 	// No hits
	// 	return true, nil
	// }
	return true, nil
}

/*
   排序查询
   返回原始Hit
   builder：排序
   agg：聚合 类似group_by sum
   query：查询
*/
func (es *MyElastic) SortQueryReturnHits(index_name string, from, size int, builder []elastic.Sorter, query []elastic.Query) (bool, []*elastic.SearchHit) {

	searchResult := es.Client.Search().Index(index_name)

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
	_, err := searchResult.Do(es.Ctx) // execute
	if err != nil {
		log.Println(es.Err)
		return false, nil
	}

	//	log.Println("wwwwww", es_result.Aggregations)
	// if es_result.Hits.TotalHits > 0 {

	// 	return true, es_result.Hits.Hits
	// } else {

	// 	return true, nil
	// }

	return true, nil
}

/*
添加记录,覆盖添加
*/
func (es *MyElastic) Add(index_name, type_name, id string, data interface{}) (result bool) {
	result = false
	// Index a tweet (using JSON serialization)
	if len(id) > 0 {
		_, es.Err = es.Client.Index().
			Index(index_name).
			Type(type_name).
			Id(id).
			BodyJson(data).
			Do(es.Ctx)
	} else {
		_, es.Err = es.Client.Index().
			Index(index_name).
			Type(type_name).
			BodyJson(data).
			Do(es.Ctx)
	}

	if es.Err != nil {
		log.Println(es.Err)
		return false
	}
	_, es.Err = es.Client.Flush().Index(index_name).Do(es.Ctx)
	if es.Err != nil {
		log.Println(es.Err)
		return false
	}
	return true
}

/*
 批量新增
*/
func (es *MyElastic) BulkAdd(index_name, type_name, id string, data []interface{}) (result bool) {
	result = false
	// Index a tweet (using JSON serialization)
	bulkRequest := es.Client.Bulk()
	if len(id) > 0 {
		for _, doc := range data {
			esRequest := elastic.NewBulkIndexRequest().
				Index(index_name).
				Type(type_name).
				Id(id).Doc(doc)
			bulkRequest = bulkRequest.Add(esRequest)
		}
	} else {
		for _, doc := range data {
			esRequest := elastic.NewBulkIndexRequest().
				Index(index_name).
				Type(type_name).
				Doc(doc)
			bulkRequest = bulkRequest.Add(esRequest)
		}
	}

	_, es.Err = bulkRequest.Do(es.Ctx)
	if es.Err != nil {
		log.Println(es.Err)
		return false
	}
	_, es.Err = es.Client.Flush().Index(index_name).Do(es.Ctx)
	if es.Err != nil {
		log.Println(es.Err)
		return false
	}
	return true

}

/*
添加记录,覆盖添加
index_name
type_name
query interface{} //查询条件
out *[]Param //查询结果
*/
func (es *MyElastic) SearchMap(index_name, type_name string, query interface{}, out *[]map[string]interface{}) (result bool) {
	es_search := es.Client.Search()
	if len(type_name) > 0 {
		es_search = es_search.Type(type_name)
	}
	if len(index_name) > 0 {
		es_search = es_search.Index(index_name)
	}
	var es_result *elastic.SearchResult
	es_result, es.Err = es_search.Source(query).Do(es.Ctx)
	if es.Err != nil {
		log.Println(es.Err)
		return false
	}
	if es_result.Hits == nil {
		log.Println(errors.New("expected SearchResult.Hits != nil; got nil"))
		return false
	}

	for _, hit := range es_result.Hits.Hits {
		tmp := make(map[string]interface{})
		err := json.Unmarshal(*hit.Source, &tmp)
		if err != nil {
			log.Println(es.Err)
		} else {
			*out = append(*out, tmp)
		}
	}

	return true
}

/*
添加记录,覆盖添加
index_name
type_name
query interface{} //查询条件
out *[]Param //查询结果
*/
func (es *MyElastic) Search(index_name, type_name string, query interface{}, f func(e []byte) error) (result bool) {

	es_search := es.Client.Search()
	if len(type_name) > 0 {
		es_search = es_search.Type(type_name)
	}
	if len(index_name) > 0 {
		es_search = es_search.Index(index_name)
	}
	var es_result *elastic.SearchResult
	es_result, es.Err = es_search.Source(query).Do(es.Ctx)
	if es.Err != nil {
		log.Println(es.Err)
		return false
	}
	if es_result.Hits == nil {
		log.Println(errors.New("expected SearchResult.Hits != nil; got nil"))
		return false
	}

	for _, hit := range es_result.Hits.Hits {
		f(*hit.Source) //如果 inmprt github.com/olivere/elastic 需要去掉 *
	}

	return true
}
