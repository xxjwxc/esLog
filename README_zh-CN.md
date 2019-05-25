# esLog
elasticsearch log golang 的elasticsearch 日志封装，包括搜索，查询，添加等

# elasticsearch 的日志封装类型

[English](README_zh-CN.md)
- 安装

```
 go get gopkg.in/olivere/elastic.v5

```

- 初始化 
  
```go 
	es.Url = "http://192.168.198.17:9200/"
	es.Index = "wms_log"

	e := es.GetClient()
```
- 添加

```go
	es.Url = "http://192.168.198.17:9200/"
	es.Index = "wms_log"

	e := es.GetClient()

	var eslog es.ESLog
    ...

	b := e.Add(es.Index, es.Index, "", eslog)
	if !b {
		fmt.Println(e.Err)
	}
```
- 搜索

```go
es.Url = "http://192.168.198.17:9200/"
es.Index = "wms_log"
    
//精确搜索
term := make(map[string]interface{})
...
//模糊匹配
match := make(map[string]interface{})
...
//时间段搜索
timeCase := make(map[string]es.CaseSection)
...

eslist := tools.Search(term, match, timeCase, req.Page, req.Limit)

```

- 注意

如果 import github.com/olivere/elastic 需要去掉 *hit.Source 的 *

- 详细文档

[详细说明](https://xie1xiao1jun.github.io/post/loglistdef/)