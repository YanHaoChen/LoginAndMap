package main

import(
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gorp.v1"
	"log"
	"crypto/rand"
	"fmt"
	"time"
)

type User struct {
  Id int64 `db:"id" json:"id"`
  Account string `db:"account" json:"account"`
	Password string `db:"password" json:"password"`
	Token string `db:"token" json:"token"`
}

type LightMapMaker struct {
	Id int64 `db:"id" json:"id"`
	Time time.Time `db:time json:"time"`
	Value float64 `db:"value" json:"value"`
	Account string `db:"account" json:"account"`
	LocationName string `db:"locationName" json:"locationName"`
}

func main(){
	r := gin.Default()
  r.Use(Cors())
	user := r.Group("User")
	{
		user.POST("/Login", Login)
		user.POST("/AuthStatus", AuthStatus)
	}
	maker := r.Group("Maker")
	{
		maker.POST("/GetLightMapMakers",GetLightMapMakers)
	}

	r.Run(":8080")
}

func Cors() gin.HandlerFunc {
 return func(c *gin.Context) {
 c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
 c.Next()
 }
}

var dbmap = initDb()

func initDb() *gorp.DbMap {
	db, err := sql.Open("mysql", "root:root@unix(/Applications/MAMP/tmp/mysql/mysql.sock)/go_db?loc=Local&parseTime=true")
	checkErr(err, "sql.Open failed")
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	dbmap.AddTableWithName(User{}, "User").SetKeys(true, "Id")
	dbmap.AddTableWithName(LightMapMaker{}, "LightMapMaker").SetKeys(true,"Id")
	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create table failed")

	return dbmap
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg,err)
	}
}

func Login(c *gin.Context) {
  var user User
	account := c.PostForm("account")
  password := c.PostForm("password")
	if account != "" && password != "" {

		err := dbmap.SelectOne(&user,"SELECT * FROM User WHERE account=? AND password=?", account, password)
			if err == nil {
				b := make([]byte, 10)
				rand.Read(b)
				str := fmt.Sprintf("%x", b)
				user.Token = str
				dbmap.Update(&user)

				c.JSON(200, gin.H{"token": str,"account": user.Account})
			} else {

				c.JSON(404, gin.H{"error": err})
			}
		} else {
		c.JSON(422, gin.H{"error": "fields are empty"})
	}
}

func AuthStatus(c *gin.Context) {
	var user User
	account := c.PostForm("account")
	token := c.PostForm("token")

	if account != "" && token != "" {
		err := dbmap.SelectOne(&user, "SELECT * FROM User WHERE account=? AND token=?", account, token)
		if err == nil {
			c.JSON(200,gin.H{"status":"pass"})
		} else {
			c.JSON(404, gin.H{"error": err})
		}
	} else {
		c.JSON(422, gin.H{"error":"fields are empty"})
	}
}

func GetLightMapMakers(c *gin.Context) {
	var user User
	var makers []LightMapMaker
	account := c.PostForm("account")
	token := c.PostForm("token")

	if token != "" && account != "" {
		err := dbmap.SelectOne(&user, "SELECT * FROM User WHERE token=? AND account=?", token,account)
		if err == nil {
			_, err := dbmap.Select(&makers, "SELECT * FROM LightMapMaker WHERE account=?",account)
			if err == nil {
				c.JSON(200,makers)
			} else {
				c.JSON(404,gin.H{"error":"something wrong"})
			}

		} else {
			c.JSON(404,gin.H{"error":"user is wrong"})
		}
	} else {
		c.JSON(422, gin.H{"error":"fields are empty"})
	}

}
