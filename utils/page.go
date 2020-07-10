package utils

type Page struct {
	//当前页码
	PageNo int
	//每页大小
	PageSize int
	//一共的页数
	TotalPage int
	//总条数
	TotalCount int
	//是否是第一页
	FirstPage bool
	//是否是最后一页
	LastPage bool
	//数据
	List interface{}
}

//总条数  当前页码  每页大小   数据list
func PageUtil(count int, pageNo int, pageSize int, list interface{}) Page {
	tp := count / pageSize
	if count%pageSize > 0 {
		tp = count/pageSize + 1
	}
	return Page{PageNo: pageNo, PageSize: pageSize, TotalPage: tp, TotalCount: count, FirstPage: pageNo == 1, LastPage: pageNo == tp, List: list}
}
