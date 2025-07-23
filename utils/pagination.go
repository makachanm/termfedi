package utils

type ItemAutoDemandPagination[T any] struct {
	max_item_limit    int
	max_item_per_page int
	items             []T

	currunt_page_pointer int
	total_page_count     int
}

func NewItemAutoDemandPagination[T any](max_item_limit int, max_per_page_limit int) *ItemAutoDemandPagination[T] {
	ptk := new(ItemAutoDemandPagination[T])
	ptk.max_item_limit = max_item_limit
	ptk.max_item_per_page = max_per_page_limit
	ptk.items = make([]T, 0)

	ptk.currunt_page_pointer = 0
	ptk.total_page_count = (ptk.max_item_limit / ptk.max_item_per_page) + (ptk.max_item_limit % ptk.max_item_per_page)

	return ptk
}

func (p *ItemAutoDemandPagination[T]) SetMaxItemPerPage(limit int) {
	//must reset all pagination status
	p.max_item_per_page = limit
	p.total_page_count = (p.max_item_limit / p.max_item_per_page) + (p.max_item_limit % p.max_item_per_page)

	p.currunt_page_pointer = 0
}

func (p *ItemAutoDemandPagination[T]) Clear() {
	p.currunt_page_pointer = 0
	p.items = p.items[:0]
}

func (p *ItemAutoDemandPagination[T]) PutItem(item T) {
	if len(p.items) >= p.max_item_limit {
		//shift
		tmp := p.items[1:]

		p.items = make([]T, p.max_item_limit)
		copy(p.items, tmp)
		p.items[p.max_item_limit-1] = item
	} else {
		p.items = append(p.items, item)
	}
}

func (p *ItemAutoDemandPagination[T]) GoPrev() {
	if p.currunt_page_pointer <= 0 {
		return
	} else {
		p.currunt_page_pointer--
	}
}

func (p *ItemAutoDemandPagination[T]) GoNext() {
	if p.currunt_page_pointer+1 >= p.total_page_count {
		return
	} else {
		p.currunt_page_pointer++
	}
}

func (p *ItemAutoDemandPagination[T]) GetTotalPage() int {
	return p.total_page_count - 1
}

func (p *ItemAutoDemandPagination[T]) GetCurruntPagePointer() int {
	return p.currunt_page_pointer
}

func (p *ItemAutoDemandPagination[T]) GetCurruntPage() []T {
	var start_pos, end_pos int

	if p.currunt_page_pointer <= 0 {
		start_pos = 0
		end_pos = p.max_item_per_page
	} else {
		start_pos = (p.currunt_page_pointer * p.max_item_per_page)

		if p.currunt_page_pointer >= p.total_page_count-1 {
			end_pos = p.max_item_limit
		} else {
			end_pos = start_pos + p.max_item_per_page
		}
	}

	return p.items[start_pos:end_pos]
}
