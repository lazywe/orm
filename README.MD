## go orm for mysql

#### 安装

````
go get github.com/lazywe/orm

````

#### 依赖
````
github.com/go-sql-driver/mysql
````

#### 使用
````
import "github.com/go-sql-driver/mysql"
// 目前只支持mysql
// mysql1
orm.GetMysql("default1").RegistDriver("mysql", "default", "user:pass@tcp(host)/")
// mysql2
orm.GetMysql("default1").RegistDriver("mysql", "default1", "user:pass@tcp(host)/")
````

#### 实例

````
// query,自定义查询
orm1 := orm.NewMaster()
res, _ := orm1.GetMysql("default").Query("select a.title from an_novel.pa_article as a limit 0 , 10")
fmt.Println(res)
````

#### 语法糖模式

````golang
 where := []orm.SqlWhere{
	orm.SqlWhere{"a","eq","666",And},
 	orm.SqlWhere{"b","neq","44",Or},
 	orm.SqlWhere{"b","in","44,33,5,1",Or},
	orm.SqlWhere{"c","find_in_set","6",Nil},
	orm.SqlWhere{"c","find_in_set_or","6,7,8",Nil},
 }
````
- where目前支持
````golang
    exp = map[string]string{
    	"eq":              "=",
    	"neq":             "!=",
    	"gt":              ">",
    	"egt":             ">=",
    	"lt":              "<",
    	"elt":             "<=",
    	"like":            "LIKE",
    	"notlike":         "NOT LIKE",
    	"in":              "IN",
    	"notin":           "NOT IN",
    	"find_in_set":     "FIND_IN_SET",
    	"find_in_set_and": "FIND_IN_SET_AND",
    	"find_in_set_or":  "FIND_IN_SET_OR",
    }
````

#### order
````golang
order := map[string]string{"id":"desc","sort":"desc"}
// id倒序，sort倒序
````


#### 获取总数
````golang
err,result := orm1.GetMysql("default1").table("default").Where([]SqlWhere{{"id","eq","1",Nil}}).GetCount("id")
// err bool
// result int64 总数量
````

#### 获取列表

````golang
err,result := orm1.GetMysql("default1").table("default").Where(where).Order(order).Limit("0,20").GetAll("*")

// err bool
// result []map[string]interface{}
````

#### 获取单条数据

````golang
err,result := orm1.GetMysql("default1").table("default").Where([]SqlWhere{{"id","eq","1",Nil}}).GetOne("*")

// err bool
// result map[string]interface{} 
````

#### 插入

````golang
datas := map[string]string{
    "title" : "孙中山",
    "age":"18",
}
err,result := orm1.GetMysql("default1").table("default").Insert(datas)
// err bool
// result int64 主键
````

#### 修改
````golang
datas := map[string]string{
    "title" : "孙中山",
    "age":"18",
}
 where := []orm.SqlWhere{
	orm.SqlWhere{"id","eq","1",Nil},
 }
/// 不要忘记添加where条件哦，后果很严重
err := orm1.GetMysql("default1").table("default").Where(where).Updates(datas)
// err bool
````

#### 删除
````golang

 where := []orm.SqlWhere{
	orm.SqlWhere{"id","eq","1",Nil},
 }
/// 不要忘记添加where条件哦，后果很严重
err := orm1.GetMysql("default1").table("default").Where(where).Delete()
// err bool
````