package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"


	"database/sql"
	_"github.com/go-sql-driver/mysql"

	"strconv"
	//"go/ast"
	//"go/types"
	//"math"
	//"encoding/binary"
	//"reflect"
	//"context"
)




func main() {

	gin.SetMode(gin.DebugMode)
	router := gin.Default()

	db, err := sql.Open("mysql", "root:gh5536856@/test?charset=utf8")
	checkErrB(err)

	fmt.Println("数据库连接完成")

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS student(id INT NOT  NULL auto_increment PRIMARY KEY ,name CHAR (16) NOT  NULL,photo VARCHAR(160) NOT  NULL DEFAULT 'https://ss0.bdstatic.com/70cFuHSh_Q1YnxGkpoWK1HF6hhy/it/u=2938685437,2474894161&fm=27&gp=0.jpg'  ,abstract VARCHAR(200) )  ;")
	checkErrB(err)

	if false {
		for a :=0;a<10; a++{
			stmt, err := db.Prepare("INSERT student SET name=?,abstract=?")
			checkErrB(err)
			_, err = stmt.Exec("张伟"+strconv.Itoa(a), "追求真理的人")
			checkErrB(err)
		}
	}



	/*for u:=0;u<10;u++{
		u-=10
	}*/

	var (
		a bool
		b bool
		c *bool = &b

		f =  [3]float64{0.12,0.32,3.54}
	)


	fmt.Println("--->")

	slice:=make([]float64 ,20)
	slice = append(slice, 0.34)

	//range也可以用在map的键值对上。
	kvs := map[string]string{"a": "apple", "b": "banana"}
	for k, v := range kvs {
		fmt.Printf("%s -> %s\n", k, v)
	}

	var f0 float32 = 30
	delete(kvs,"a")
	fmt.Println(strconv.FormatFloat(f[0] ,'E',-1,32)+"<---")
	fmt.Println(strconv.FormatFloat(float64(f0) ,'E',-1,32)+"<-")
	//var m,n,v chan  int(f[0]))
	if a||*c {

	}

	//http://127.0.0.1:8080/students/?index=0&num=5
	router.POST("/students", func(context *gin.Context) {


		index,err := strconv.Atoi(context.Query("index"))
		checkErrB(err)
		num,err :=   strconv.Atoi(context.Query("num"))
		checkErrB(err)
		strSql := "SELECT * FROM student limit "+strconv.Itoa(index*num)+","+strconv.Itoa(num)+";"
		fmt.Println(strSql)
		rows, err := db.Query(strSql)
		checkErrB(err)
		/*str,err := rows.Columns()
		checkErrB(err)*/

		type Item struct {
			Id int `json:"Id"`
			Name string `json:"Name"`
			Photo string    `json:"Photo"`
			Abstract string    `json:"Abstract"`
		}

		type JsonHolder struct {
			Code int `json:code`
			Msg string `json:msg`
			Item []Item `json:data`
		}




		var holder JsonHolder
		var id int
		var name string
		var photo string
		var abstract string
		if !rows.Next(){
			holder.Code = 1004
			holder.Msg = "未查询到数据"
			holder.Item =append(holder.Item)
		}else {
			holder.Code = 1000
			holder.Msg = "获取成功"
			err = rows.Scan(&id, &name,&photo,&abstract)
			checkErrB(err)
			holder.Item =append(holder.Item,Item{Id:id,Name:name,Photo:photo,Abstract:abstract})
			//若返回json数据，可以直接使用gin封装好的JSON方法
			for rows.Next() {
				err = rows.Scan(&id, &name,&photo,&abstract)
				checkErrB(err)
				holder.Item =append(holder.Item,Item{Id:id,Name:name,Photo:photo,Abstract:abstract})
				//fmt.Println("-->id=%d,name=%s,photo=%s,abstract=%s",id,name,photo,abstract)
			}
		}







		context.JSON(http.StatusOK, holder)
		//return
	})

	router.Run("127.0.0.1:8080")

}



func checkErrB(err error) {
	if err != nil {
		panic(err)
	}
}
