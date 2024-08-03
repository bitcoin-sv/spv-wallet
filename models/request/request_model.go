package request

type Pageable struct {
	Page int
	Size int
	Sort Sort
}

type Order struct {
	Property  string
	Direction string // TODO: convert to enumerable
}

type Sort struct {
	Orders []Order
}
