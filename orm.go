package orm

type Orm struct {
}

// debug
var (
	debug = true
)

type Params map[string]interface{}

// 注册驱动
func RegistDriver(d string, connect string, handel string) {
	// 链接数据库
	switch d {
	case "mysql":
		if err := Mh.RegisterMysql(connect, handel); err != nil {
			panic(err)
		}
		break
	}
}

//使用mysql
func (o *Orm) GetMysql(connect string) *MysqlClient {
	m, err := Mh.NewMysqlMaster(connect)
	if err != nil {
		panic(err)
	}
	return m
}
