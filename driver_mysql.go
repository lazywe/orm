package orm

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type MysqlHandel struct {
	mysqlClient map[string]*MysqlClient
}

type MysqlClient struct {
	db  *sql.DB
	sql *Sql
}

var (
	Mh *MysqlHandel
)

func init() {
	Mh = &MysqlHandel{
		mysqlClient: map[string]*MysqlClient{},
	}
}

// reg mysql
func (mh *MysqlHandel) RegisterMysql(connect string, handel string) error {
	// Open database connection
	dbc, err := sql.Open("mysql", handel)
	if err != nil {
		return err
	}
	if mh.mysqlClient[connect] != nil {
		return fmt.Errorf("链接:" + connect + "已经存在了")
	}
	mh.mysqlClient[connect] = &MysqlClient{
		db: dbc,
	}
	log.Println("mysql链接成功")
	return nil
}

// mysqlhandel
func (mh *MysqlHandel) NewMysqlMaster(connect string) (*MysqlClient, error) {
	if conn := mh.mysqlClient[connect]; conn != nil {
		return conn, nil
	}
	return nil, fmt.Errorf("connect:" + connect + "不存在")
}

// query
func (m *MysqlClient) Query(str string) ([]Params, error) {
	rows, err := m.db.Query(str)
	if debug {
	}
	if err != nil {
		return nil, err
	}
	col, err1 := rows.Columns()
	if err1 != nil {
		return nil, err
	}
	var (
		refs    []interface{}
		result  Params
		results = []Params{}
	)
	// 解析列表
	for rows.Next() {
		refs = make([]interface{}, len(col))
		for i := range col {
			var ref sql.NullString
			refs[i] = &ref
		}
		if err := rows.Scan(refs...); err != nil {
			return nil, err
		}
		result = make(Params)
		for k, v := range refs {
			var vt = v.(*sql.NullString)
			if vt.Valid {
				result[col[k]] = vt.String
			} else {
				result[col[k]] = ""
			}
		}
		results = append(results, result)
	}
	return results, nil
}

// exec
func (m *MysqlClient) Exec(str string) (sql.Result, error) {
	res, err := m.db.Exec(str)
	return res, err
}
