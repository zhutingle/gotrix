package handler

import (
	"database/sql"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/zhutingle/gotrix/global"
	"regexp"
)

type handleSql struct {
	db *sql.DB
}

func (this *handleSql) init() *handleSql {
	var err error
	this.db, err = sql.Open("mysql", global.Config.Database.Url)
	if err != nil {
		log.Println(err)
		return this
	}
	this.db.SetMaxOpenConns(global.Config.Database.MaxOpenConns)
	this.db.SetMaxIdleConns(global.Config.Database.MaxIdleConns)
	this.db.Ping()
	return this
}

func (this *handleSql) handle(job *Job, cp *global.CheckedParams) (result interface{}, gErr *global.GotrixError) {
	args := make([]interface{}, 0)
	sqlStr := job.Job

	var sqlArgsFunc = func(str string) string {
		var name = str[strings.Index(str, "{") + 1:strings.LastIndex(str, "}")]
		if _, ok := cp.V[name].([]interface{}); ok {
			var array []interface{} = cp.V[name].([]interface{})
			args = append(args, array...)
			var returnStr = "?"
			for k := 0; k < len(array) - 1; k++ {
				returnStr += ",?"
			}
			return returnStr
		} else {
			args = append(args, cp.V[name])
			return "?"
		}
	}

	if job.auto {
		// 对 <auto></auto> 标签按实际参数进行处理
		// 如：<auto>tel = ${tel},email = ${email},icon = ${icon}</auto>  其中 tel 和 email 不为空，则处理之后变成 tel=?,email=?
		// 如：<auto>tel = ${tel} and email = ${email} and icon = ${icon}</auto> 其中 tel 和 email 不为空，则处理之后变成 tel = ? and email = ?
		sqlStr = autoTagReg.ReplaceAllStringFunc(sqlStr, func(str string) string {
			str = str[6:len(str) - 7] // 去除 <auto></auto> 标签
			str = autoTagItemReg.ReplaceAllStringFunc(str, func(item string) string {
				var name = item[strings.Index(item, "{") + 1:strings.LastIndex(item, "}")]
				if cp.V[name] == nil {
					return ""
				} else {
					return sqlArgsReg.ReplaceAllStringFunc(item, sqlArgsFunc)
				}
			})
			// 清除掉多余的空格、逗号、and
			str = strings.TrimSpace(str)
			str = strings.TrimSuffix(str, ",")
			str = strings.TrimSuffix(str, "and")
			str = strings.TrimSuffix(str, "where")
			str = strings.TrimSuffix(str, "set")
			return str
		})
	}

	sqlStr = sqlArgsReg.ReplaceAllStringFunc(sqlStr, sqlArgsFunc)

	total := 0
	// 具有分页参数的情况下，增加获取当前条件下的总条数
	if cp.V["pNum"] != nil && cp.V["pSize"] != nil {

		// 更改查询语句为 select count(1) from ...
		countSqlStr := regexp.MustCompile("select.*?from").ReplaceAllString(sqlStr, "select count(1) from")
		countStmt, err := this.db.Prepare(countSqlStr)
		if err != nil {
			log.Println(err)
			gErr = global.SQLHANDLE_PREPARE_ERROR
			return
		}
		// 查询数量并获取
		countRow := countStmt.QueryRow(args...)

		err = countRow.Scan(&total)
		if err != nil {
			log.Println(err)
			gErr = global.SQLHANDLE_QUERY_ERROR
			return
		}

		// 增加分页参数，构造分页查询的 SQL 语句。
		sqlStr += " limit ?,?"
		args = append(args, cp.V["pNum"], cp.V["pSize"])
	}

	stmt, err := this.db.Prepare(sqlStr)
	if err != nil {
		log.Println(err)
		gErr = global.SQLHANDLE_EXECUTE_ERROR
		return
	}

	if regexp.MustCompile("^(?i:select)").MatchString(sqlStr) {
		// 对 select 语句进行处理
		rows, err := stmt.Query(args...)
		if err != nil {
			log.Println(err)
			gErr = global.SQLHANDLE_QUERY_ERROR
			return
		}
		defer rows.Close()

		columnNames, err := rows.Columns()
		if err != nil {
			log.Println(err)
			gErr = global.SQLHANDLE_COLUMNS_ERROR
			return
		}
		columnCount := len(columnNames)
		data := make([]interface{}, 0)
		for rows.Next() {
			row := make([]interface{}, columnCount)
			for j := 0; j < len(row); j++ {
				var cell interface{}
				row[j] = &cell
			}
			err = rows.Scan(row...)
			if err != nil {
				log.Println(err)
				gErr = global.SQLHANDLE_SCAN_ERROR
				return
			}
			column := make(map[string]interface{})
			for j := 0; j < len(row); j++ {
				cell := *row[j].(*interface{})
				if _, ok := cell.([]uint8); ok {
					column[columnNames[j]] = string(cell.([]uint8))
				} else {
					column[columnNames[j]] = cell
				}
			}
			data = append(data, column)
		}
		// 对于不不同的 type 封装不同的数据格式
		// single ：单条数据，没有数据返回 nil
		if job.Type == "single" {
			if len(data) > 0 {
				result = data[0]
			}
			// pagination ： 多条数据，带总条数 {data:[...],total: 100} ，没有数据时返回 {data:[], total: 0}
		} else if job.Type == "pagination" {
			resultMap := make(map[string]interface{})
			resultMap["data"] = data
			resultMap["total"] = total
			result = resultMap
		} else {
			// 默认 ： 以数组的形式输入数据，没有数据时为 []
			result = data
		}

	} else {
		// 对 insert、delete、update 语句进行处理。
		res, err := stmt.Exec(args...)
		if err != nil {
			log.Println(err)
			gErr = global.SQLHANDLE_EXEC_ERROR
			return
		}
		if strings.HasPrefix(sqlStr, "insert") {
			id, err := res.LastInsertId()
			if err != nil {
				log.Println(err)
				return
			}
			result = id
		} else {
			affected, err := res.RowsAffected()
			if err != nil {
				log.Println(err)
				return
			}
			result = affected
		}
	}

	return
}
