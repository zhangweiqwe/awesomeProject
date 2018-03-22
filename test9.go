package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	//"strconv"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"log"
	"net/http"
	"user"
	"time"
	//"go/token"
)

func main() {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()

	db, err := sql.Open("mysql", "root:gh5536856@/test?charset=utf8")
	checkErrZ(err)

	if err := db.Ping(); err != nil {
		fmt.Println("数据库ping 失败！")
		log.Fatalln(err)
	}

	fmt.Println("数据库连接完成")

	//ALTER TABLE user CHANGE phone phone int unsigned not null;
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS user(id INT NOT  NULL auto_increment ,phone VARCHAR(16) NOT NUll,password VARCHAR (16) NOT NUll,name CHAR (16) NOT  NULL,photo VARCHAR(160) NOT  NULL DEFAULT 'https://ss0.bdstatic.com/70cFuHSh_Q1YnxGkpoWK1HF6hhy/it/u=2938685437,2474894161&fm=27&gp=0.jpg'  ,description VARCHAR(200),PRIMARY KEY (id,phone) )  ;")
	checkErrZ(err)

	defer db.Close()

	//127.0.0.1:8090/register?phone=15923549211&name=张伟&password=zxcvbnm
	r.POST("/register", func(context *gin.Context) {

		phone := context.Query("phone")
		checkErrZ(err)
		password := context.Query("password")
		checkErrZ(err)
		name := context.Query("name")
		checkErrZ(err)

		strSql := "SELECT * FROM user Where phone= " + phone + ";"
		fmt.Println(strSql)
		rows, err := db.Query(strSql)
		checkErrZ(err)

		type RegisterJson struct {
			Code  int    `json:code`
			Msg   string `json:msg`
			Token string `json:token`
		}

		var json RegisterJson

		if rows.Next() {
			json.Code = user.CODE_ERROR
			json.Msg = "用户已存在！"
			context.JSON(http.StatusOK, json)
			return
		}

		stmt, err := db.Prepare("INSERT user SET phone=?,name =?,password = ?")
		checkErrZ(err)
		res, err := stmt.Exec(phone, name, password)
		checkErrZ(err)
		i, err := res.LastInsertId()
		checkErrZ(err)
		println("数据库被改变的地方---->i=", i)

		token := jwt.New(jwt.SigningMethodHS256)
		claims := make(jwt.MapClaims)
		claims["phone"] = phone
		claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
		claims["iat"] = time.Now().Unix()
		token.Claims = claims

		tokenString, err := token.SignedString([]byte(user.SecretKey))
		checkErrZ(err)

		json.Code = user.CODE_SUCCESS
		json.Msg = "注册成功！"
		json.Token = tokenString

		context.JSON(http.StatusOK, json)
	})

	//127.0.0.1:8090/resource
//eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MjE2OTQyMDcsImlhdCI6MTUyMTY5MDYwNywicGhvbmUiOiIxNTkyMzU0OTIxMSJ9.LkQNohKS8g3fcof8Bj-F7nI1Rju-vAEOpd-c9IhF1F4

	r.POST("/resource", func(context *gin.Context) {

		token, err := request.ParseFromRequest(context.Request, request.AuthorizationHeaderExtractor,
			func(token *jwt.Token) (interface{}, error) {
				return []byte(user.SecretKey), nil
			})
		//checkErrZ(err)

		type AuthenticationSuccessJson struct {
			Code  int    `json:code`
			Msg   string `json:msg`
			Token string `json:token`
		}

		var sucessJson AuthenticationSuccessJson


		if err == nil{
			if token.Valid {
				sucessJson.Code = user.CODE_SUCCESS
				sucessJson.Msg = "可以尝试访问资源"


				claims,ok :=token.Claims.(jwt.MapClaims)
				if!ok{
					sucessJson.Code = user.CODE_ERROR
					sucessJson.Msg = "尝试访问失败，用户信息异常"
				}else {
					var phone string = claims["phone"].(string)
					fmt.Println("还原信息",phone)
				}




			} else {
				sucessJson.Code = user.CODE_ERROR
				sucessJson.Msg = "尝试访问失败"
			}
		}else {
			sucessJson.Code = user.CODE_ERROR
			sucessJson.Msg = "尝试访问失败"
		}


		context.JSON(http.StatusOK, sucessJson)

	})

	//127.0.0.1:8090/resource?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MjE2OTAyMTIsImlhdCI6MTUyMTY4NjYxMn0.SZUd29hNMxH7f9mfo7m9i_Butod6Jf0DS9l5yOd388Q
	//http.Handle("/resource", nil)

	r.Run(":8090")
}

func checkErrZ(err error) {
	if err != nil {
		panic(err)
	}
}
