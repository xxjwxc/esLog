package es

type EsQuery struct {
	//qmp       map[string]interface{}
	much_term []interface{}
	offset    int32
	limit     int32
}

func (q *EsQuery) OnPages(offset, limit int32) *EsQuery {
	q.offset = offset
	q.limit = limit

	if q.offset < 0 {
		q.offset = 0
	}
	if q.limit <= 0 {
		q.limit = 1
	}
	return q
}

//精确搜索
func (q *EsQuery) OnTerm(parm map[string]interface{}) *EsQuery {
	if len(parm) == 0 {
		return q
	}

	for k, v := range parm {
		mp := make(map[string]interface{})
		mp[k] = v
		q.much_term = append(q.much_term, map[string]interface{}{
			"term": mp,
		})
	}

	return q
}

//模糊匹配
func (q *EsQuery) OnMatch(parm map[string]interface{}) *EsQuery {
	if len(parm) == 0 {
		return q
	}

	for k, v := range parm {
		mp := make(map[string]interface{})
		mp[k] = v
		q.much_term = append(q.much_term, map[string]interface{}{
			"match": mp,
		})
	}

	return q
}

//模糊匹配
func (q *EsQuery) OnWildcard(parm map[string]interface{}) *EsQuery {
	if len(parm) == 0 {
		return q
	}

	for k, v := range parm {
		mp := make(map[string]interface{})
		mp[k] = v
		q.much_term = append(q.much_term, map[string]interface{}{
			"wildcard": mp,
		})
	}

	return q
}

//区间<oplogger.TimeSection{BeginTime,EndTime}	>
func (q *EsQuery) OnRangeTime(timeCase map[string]CaseSection) *EsQuery {
	if timeCase == nil || len(timeCase) == 0 {
		return q
	}

	for k, v := range timeCase {
		mp := make(map[string]interface{})
		mp["gte"] = v.Min //大于等于
		mp["lt"] = v.Max  //小于 <左闭右开>
		pm := make(map[string]interface{})
		pm[k] = mp
		q.much_term = append(q.much_term, map[string]interface{}{
			"range": pm,
		})
	}
	return q
}

//区间<oplogger.TimeSection{BeginTime,EndTime}	>
func (q *EsQuery) OnSource() map[string]interface{} {

	if len(q.much_term) == 0 {
		return map[string]interface{}{
			"from":        q.offset,
			"size":        q.limit,
			"create_time": "desc", //默认时间降序
		}
	}

	return map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": q.much_term,
			},
		},
		"from": q.offset,
		"size": q.limit,
		"sort": map[string]interface{}{
			"create_time": "desc",
		}, //默认时间降序
	}
}
