package gateway

type Paginator interface {
	PerPage() int
	Page() int
	SetPage(int)
	SetLimit(int)
	Total() int
	SetTotal(int)
}

type paginator struct {
	limit int
	page  int
	total int
}

func NewPaginator() Paginator {
	return &paginator{}
}

func (p *paginator) Page() int {
	if p.page == 0 {
		p.page = 1
	}
	return p.page
}

func (p *paginator) PerPage() int {
	if p.limit == 0 {
		p.limit = 10
	}
	return p.limit
}

func (p *paginator) SetPage(i int) {
	p.page = i
}

func (p *paginator) SetLimit(i int) {
	p.limit = i
}

func (p *paginator) Total() int {
	return p.total
}

func (p *paginator) SetTotal(total int) {
	p.total = total
}
