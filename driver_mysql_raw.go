package orm

import (
	"strconv"
	"strings"

	"fmt"
)

//  设置数据表
//  @table string 数据表名称
//  @return *Sql
func (m *MysqlClient) Table(table string) *MysqlClient {
	m.sql = &Sql{
		table: table,
	}
	return m
}

// 获取所有数据
// @param string field 字段，跟原生写法一样 id as id
// @demo orm.getMysql("default").table("default").Where([]SqlWhere{{"id","eq","1",Nil}}).GetAll("*")
// @return []orm.Params, error
func (m *MysqlClient) GetAll(field string) ([]Params, error) {
	/// sql 构造
	sql := "SELECT " + field + " FROM " + m.sql.table + m.sql.where + m.sql.orderBy + m.sql.limit
	res, err := m.Query(sql)
	return res, err
}

/**
 * 获取数据总条数
 * @param string field 字段，跟原生写法一样 id as id
 * @demo sql.table("default").Where([]SqlWhere{{"id","eq","1",Nil}}).GetCount("id")
 * @return bool,int64
 */
func (m *MysqlClient) GetCount() (int64, error) {
	/// sql 构造
	sql := "SELECT count(*) as num FROM " + m.sql.table + m.sql.where
	res, err := m.Query(sql)
	if err != nil {
		return 0, err
	}
	num, err := strconv.Atoi(res[0]["num"].(string))
	if err != nil {
		return 0, err
	}
	return int64(num), nil
}

/**
 * 获取单条数据
 * @param string field 字段，跟原生写法一样 id as id
 * @demo sql.table("default").Where([]SqlWhere{{"id","eq","1",Nil}}).GetOne("*")
 * @return []Params,bool
 */
func (m *MysqlClient) GetOne(field string) (Params, error) {
	defer m.gc()
	if m.sql.table == "" {
		return nil, fmt.Errorf("table is nil.")
	}
	/// sql 构造
	sql := "SELECT " + field + " FROM " + m.sql.table + m.sql.where + m.sql.orderBy + " limit 0,1"
	res, err := m.Query(sql)
	if err != nil {
		return nil, err
	}
	if len(res) > 0 {
		return res[0], nil
	}
	return nil, nil
}

// 修改
// @param map[string]string datas 插入的字段
// @return bool
func (m *MysqlClient) Update(datas map[string]string) error {
	defer m.gc()
	if len(datas) < 0 {
		return fmt.Errorf("update datas is nil.")
	}
	if m.sql.table == "" || m.sql.where == "" {
		return fmt.Errorf("table or where is nil.")
	}
	datastring := ""
	for k, v := range datas {
		if datastring != "" {
			datastring += ", "
		}
		datastring += fmt.Sprintf("%s = '%s'", k, v)
	}
	sql := "UPDATE " + m.sql.table + " set " + datastring + m.sql.where
	_, err := m.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}

// 字段增加数值
// @param string  key 增加的字段
// @param int  var	  增加的值
// @return bool
func (m *MysqlClient) Inc(key string, val int) error {
	defer m.gc()
	if m.sql.where == "" || m.sql.table == "" {
		return fmt.Errorf("where or table is nil.")
	}
	datastring := ""
	datastring += fmt.Sprintf("%s = %s+%d", key, key, val)
	sql := "UPDATE " + m.sql.table + " set " + datastring + m.sql.where
	_, err := m.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}

// 删除
// @return bool
func (m *MysqlClient) Delete() error {
	defer m.gc()
	if m.sql.where == "" || m.sql.table == "" {
		return fmt.Errorf("table or where is nil.")
	}
	sql := fmt.Sprintf("DELETE from %s %s", m.sql.table, m.sql.where)
	_, err := m.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}

// 新增
// @param map[string]string datas key字段，value值
// @return int64, error
func (m *MysqlClient) Insert(datas map[string]string) (int64, error) {
	defer m.gc()
	if len(datas) < 0 {
		return 0, fmt.Errorf("insert datas datas is nil.")
	}
	datakey := ""
	datavalue := ""
	for k, v := range datas {
		if datakey != "" {
			datakey += ", "
			datavalue += ", "
		}
		datakey += k
		datavalue += fmt.Sprintf("'%s'", v)
	}
	sql := fmt.Sprintf("Insert into %s (%s)"+" VALUES (%s)", m.sql.table, datakey, datavalue)
	result, err := m.Exec(sql)
	if err != nil {
		return 0, err
	}
	num, err2 := result.LastInsertId()
	if err2 != nil {
		return 0, err2
	}
	return num, nil
}

// limit
// @param string str limit字符 0,10       0偏移量，10查询数量
// @return *MysqlClient
func (m *MysqlClient) Limit(offset string, pagesize string) *MysqlClient {
	if offset != "" && pagesize != "" {
		m.sql.limit = fmt.Sprintf(" limit %s,%s", offset, pagesize)
	}
	return m
}

// order
// @param map[string]string maps order
// @demo map[string]string{"id":"desc","sort":"desc"}
// @return *Sql
func (m *MysqlClient) Order(maps map[string]string) *MysqlClient {
	order := ""
	if len(maps) > 0 {
		for k, v := range maps {
			if order != "" {
				order += ","
			}
			order += k + " " + v
		}
		order = " order by " + order
		m.sql.orderBy = order
	}
	return m
}

// 根据数组构建sqlwhere*
// @param []string 搜索数据接口
// @demo
// where := []SqlWhere{
//		SqlWhere{"a","eq","666",And},
//		SqlWhere{"b","neq","44",Or},
//		SqlWhere{"c","find_in_set","6",Nil},
// @return *MysqlClient
func (m *MysqlClient) Where(where []SqlWhere) *MysqlClient {
	if where != nil && len(where) > 0 {
		wherestring := " where "
		for k, v := range where {
			result := m.sql.buildCondition(v.Field, v.Condition, v.Value)
			wherestring += result
			if k+1 != len(where) {
				wherestring += v.Operation.String()
			}
			if v.Operation == Nil {
				break
			}
		}
		m.sql.where = wherestring
	}
	return m
}

/**
 * 自定义条件类型
 * @param string field 条件的字段
 * @param string condition 条件类型
 * @param string value 条件的值，此处可扩展
 * @return string
 **/
func (s *Sql) buildCondition(field string, condition string, value string) string {
	result := ""
	switch condition {
	case "eq", "neq", "gt", "lt", "egt", "elt", "like", "notlike":
		result += fmt.Sprintf("%s %s '%s'", field, exp[condition], value)
	case "in", "notin":
		// result += field + " IN(" + value + ")"

		valuearr := strings.Split(value, ",")
		valuestring := strings.Join(valuearr, "','")
		result += fmt.Sprintf("%s %s('%s')", field, exp[condition], valuestring)

	case "find_in_set":
		result += "FIND_IN_SET('" + value + "'," + field + ")"
	case "find_in_set_or":
		arr := strings.Split(value, ",")
		result += "( "
		for k, v := range arr {
			if k != 0 {
				result += " or "
			}
			result += "FIND_IN_SET('" + v + "'," + field + ")"
		}
		result += ")"
	case "find_in_set_and":
		arr := strings.Split(value, ",")
		result += "( "
		for k, v := range arr {
			if k != 0 {
				result += " and "
			}
			result += "FIND_IN_SET('" + v + "'," + field + ")"
		}
		result += ")"
	default:
		return value
	}
	return result
}

// gc
func (m *MysqlClient) gc() {
	m.sql = nil
}
