package helper

type Condition struct {
	Query string
	Args  []interface{}
}

type QueryOptions struct {
	Where   []Condition
	OrderBy string
	Take    int
	Skip    int
}

type Meta struct {
	TotalItem   int64 `json:"total_item"`
	TotalPage   int   `json:"total_page"`
	CurrentPage int   `json:"current_page"`
}
