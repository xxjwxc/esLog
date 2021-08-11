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
	e,err := New(WithIndexName("wms_log"), WithAddrs("http://192.168.198.17:9200/"))
```

- add

```go
	e,err := New(WithIndexName("wms_log"), WithAddrs("http://192.168.198.17:9200/"))

	var eslog es.ESLog
    ...

	b := e.Add(es.Index, es.Index, "", eslog)
	if !b {
		fmt.Println(e.Err)
	}
```
- search

```go
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

[link](https://xxjwxc.github.io/post/loglistdef/)