# esLog
elasticsearch log golang 的elasticsearch 

Log encapsulation, including search, query, add, etc.

[中文版](README_zh-CN.md)

# elasticsearch Log Encapsulation Types

- install

```
 go get gopkg.in/olivere/elastic.v5

```

- init 
  
```go 
	es.Url = "http://192.168.198.17:9200/"
	es.Index = "wms_log"

	e := es.GetClient()
```

- add

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
- search

```go
es.Url = "http://192.168.198.17:9200/"
es.Index = "wms_log"
    
//Precise search
term := make(map[string]interface{})
...
//Fuzzy matching
match := make(map[string]interface{})
...
//Time Search
timeCase := make(map[string]es.CaseSection)
...

eslist := tools.Search(term, match, timeCase, req.Page, req.Limit)

```
- case

if import github.com/olivere/elastic must remove *hit.Source on *

- more
[link](https://xie1xiao1jun.github.io/post/loglistdef/)