package prewarm

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	url "net/url"
	"os"
	"strconv"

	_ "github.com/denisenkom/go-mssqldb"
	redis "github.com/go-redis/redis/v9"

	"hackweek41/utils"
	"hackweek41/wwwmodels"
)

var ctx = context.Background()

func MainDriver() {
	fmt.Println("In Prewarm")

	db_user := utils.GoDotEnvVariable("DB_USER")
	db_pass := utils.GoDotEnvVariable("DB_PASSWORD")
	db_port := utils.GoDotEnvVariable("DB_PORT")
	db_server := utils.GoDotEnvVariable("DB_SERVER")

	f, file_err := os.Create("data.txt")

	if file_err != nil {
		log.Fatal(file_err)
	}

	defer f.Close()

	query := url.Values{}
	query.Add("app name", "hackweek41")
	db_url := &url.URL{
		Scheme: "sqlserver",
		User: url.UserPassword(db_user, db_pass),
		Host: fmt.Sprintf("%s:%s", db_server, db_port),
		RawQuery: query.Encode(),
	}

	db, db_err := sql.Open("sqlserver", db_url.String())

	if db_err != nil {
		log.Fatal(db_err)
	}

	redis_url := utils.GoDotEnvVariable("REDIS_URL")

	rdb := redis.NewClient(&redis.Options{
		Addr:     redis_url,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	record_rows, _ := db.Query("SELECT userid, usertype_id FROM www.dbo.users")

	for record_rows.Next() {
		u := wwwmodels.User{}
		err := record_rows.Scan(
			&u.Userid, 
			&u.Usertype_id,
		)

		if err != nil {
			log.Fatal(err)
		}
		
		str_user_id := strconv.Itoa(u.Userid)
		str_usertype_id := strconv.Itoa(u.Usertype_id)
		val, rdb_err := rdb.Set(ctx, "userid:"+str_user_id, str_usertype_id, 0).Result()
		if rdb_err != nil {
			log.Fatal(rdb_err)
		}

		fmt.Println("REDIS RES: ", val)
		fmt.Println("USERID: ", u.Userid)
		fmt.Println("USERTYPE_ID: ", str_usertype_id)
		
		_, err2 := f.WriteString(str_user_id+"\n")
		
		if err2 != nil {
			log.Fatal(err2)
		}
	}
}
	
