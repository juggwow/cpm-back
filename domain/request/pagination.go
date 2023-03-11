package request

const LIMIT = 10

type Pagination struct {
	Page  int
	Limit int
}

func (p Pagination) GetPage() int {
	page := p.Page
	if page <= 0 {
		return 1
	}
	return page
}

func (p Pagination) GetLimit() int {
	if p.Limit <= 0 {
		return LIMIT
	}
	return p.Limit
}

func (p Pagination) Offset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func GetPagination(page int, pageSize int) Pagination {
	return Pagination{page, pageSize}
}
