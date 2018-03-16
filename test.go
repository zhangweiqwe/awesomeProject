package main

import (
	"database/sql"
	"fmt"
	_"github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:gh5536856@/test?charset=utf8")
	checkErr(err)

	db.Prepare("CREATE TABLE student(id INT NOT  NULL auto_increment PRIMARY KEY ,name CHAR (16) NOT  NULL,photo VARCHAR(160) NOT  NULL DEFAULT 'https://ss0.bdstatic.com/70cFuHSh_Q1YnxGkpoWK1HF6hhy/it/u=2938685437,2474894161&fm=27&gp=0.jpg'  ,abstract VARCHAR(200) );")


	db.Prepare("INSERT INTO student(name ,abstract) VALUES ('张伟','追求真理'),('糖霜','傻');")
	db.Prepare("SELECT * FROM student;")
	db.Prepare("DELETE FROM student WHERE id = 1 OR id=2;")
	// insert
	stmt, err := db.Prepare("INSERT user_info SET id=?,name=?")
	checkErr(err)

	res, err := stmt.Exec(1, "wangshubo")
	checkErr(err)

	// update
	stmt, err = db.Prepare("update user_info set name=? where id=?")
	checkErr(err)

	res, err = stmt.Exec("wangshubo_update", 1)
	checkErr(err)


	/*aff_nums, _ :=res.RowsAffected()
	fmt.Println(aff_nums)*/
		//int6,err = res.RowsAffected()

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println(affect)

	// query
	rows, err := db.Query("SELECT * FROM user_info")
	checkErr(err)

	for rows.Next() {
		var uid int
		var username string

		err = rows.Scan(&uid, &username)
		checkErr(err)
		fmt.Println(uid)
		fmt.Println(username)
	}

	// delete
	stmt, err = db.Prepare("delete from user_info where id=?")
	checkErr(err)

	res, err = stmt.Exec(1)
	checkErr(err)

	// query
	rows, err = db.Query("SELECT * FROM user_info")
	checkErr(err)

	for rows.Next() {
		var uid int
		var username string

		err = rows.Scan(&uid, &username)
		checkErr(err)
		fmt.Println(uid)
		fmt.Println(username)
	}

	db.Close()

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
