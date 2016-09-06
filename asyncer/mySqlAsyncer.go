package asyncer

import (
	//	"database/sql"
	"fmt"
	//	_ "github.com/go-sql-driver/mysql"

	"github.com/zhutingle/gotrix/global"
)

var ch chan *global.CheckedParams

type MySqlAsyncer struct {
}

func (this MySqlAsyncer) Init() {

	ch = make(chan *global.CheckedParams, 1024)

	go runAsync()

	//	db, err := sql.Open("mysql", "root:root@tcp(192.168.0.248:3306)/thribu?charset=utf8")
	//	checkErr(err)
	//
	//	stmt, err := db.Prepare("INSERT INTO goods VALUES(?,?,?,?,?,?)")
	//	checkErr(err)
	//
	//	res, err := stmt.Exec(3, "name", "des", "imgs", 0, nil)
	//	checkErr(err)
	//
	//	id, err := res.LastInsertId()
	//	checkErr(err)
	//	fmt.Println(id)
	//
	//	affected, err := res.RowsAffected()
	//	checkErr(err)
	//	fmt.Println(affected)
	//
	//	db.Close()
}

func (this MySqlAsyncer) Async(checkedParams *global.CheckedParams) {

	ch <- checkedParams

}

func runAsync() {
	for checkedParams := range ch {
		fmt.Println(checkedParams.V)
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
