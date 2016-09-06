package handler

import (
	"database/sql"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/zhutingle/gotrix/global"
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
	//TODO 待优化
	handleSql := sqlArgsReg.ReplaceAllStringFunc(job.Job, func(str string) string {
		var name = str[2 : len(str)-1]
		if _, ok := cp.V[name].([]interface{}); ok {
			var array []interface{} = cp.V[name].([]interface{})
			args = append(args, array...)
			var returnStr = "?"
			for k := 0; k < len(array)-1; k++ {
				returnStr += ",?"
			}
			return returnStr
		} else {
			args = append(args, cp.V[name])
			return "?"
		}
	})
	stmt, err := this.db.Prepare(handleSql)
	if err != nil {
		log.Println(err)
		gErr = global.SQLHANDLE_PREPARE_ERROR
		return
	}

	if strings.HasPrefix(handleSql, "select") {
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
		if job.One {
			if len(data) > 0 {
				result = data[0]
			}
		} else {
			result = data
		}
	} else {
		res, err := stmt.Exec(args...)
		if err != nil {
			log.Println(err)
			gErr = global.SQLHANDLE_EXEC_ERROR
			return
		}
		if strings.HasPrefix(handleSql, "insert") {
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
