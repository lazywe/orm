package orm

type Sql struct {
	table   string
	where   string
	orderBy string
	limit   string
}
type SqlWhere struct {
	Field     string
	Condition string
	Value     string
	Operation Operation
}

type SqlInterface interface {
	Init()
	GetAll()
	GetOne()
	Where()
	Order()
	Limit()
	Table()
}

/// 定义sql操作枚举
type Operation int

// iota 初始化后会自动递增
const (
	And    Operation = iota // value --> 0
	Or                      // value --> 1
	AndNot                  // value --> 2
	OrNot                   // value --> 3
	Nil                     // value --> 4
)

func (this Operation) String() string {
	switch this {
	case And:
		return " and "
	case Or:
		return " or "
	case AndNot:
		return " and not "
	case OrNot:
		return " or not "
	default:
		return ""
	}
}

var (
	exp = map[string]string{
		"eq":          "=",
		"neq":         "!=",
		"gt":          ">",
		"egt":         ">=",
		"lt":          "<",
		"elt":         "<=",
		"like":        "LIKE",
		"notlike":     "NOT LIKE",
		"in":          "IN",
		"notin":       "NOT IN",
		"find_in_set": "FIND_IN_SET",
	}
)
