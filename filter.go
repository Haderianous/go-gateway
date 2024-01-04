package gateway

type Operation string

const (
	Eq    Operation = "eq"
	Ct    Operation = "ct"
	Bt    Operation = "bt"
	NotEq Operation = "neq"
)

type Filter struct {
	Key   string        `json:"k"`
	Value []interface{} `json:"v"`
	Op    Operation     `json:"op"`
}

type Sort struct {
	Key   string `json:"k"`
	Value any    `json:"v"`
}

type FilterParams struct {
	Filters []Filter `json:"filters"`
	Sorts   []Sort   `json:"sorts"`
	Page    int      `json:"page"`
	Limit   int      `json:"limit"`
}

func (op Operation) ToSql() string {
	switch op {
	case Bt:
		return "between"
	case Eq:
		return "="
	case NotEq:
		return "!="
	case Ct:
		return "like"
	default:
		return "="
	}
}
