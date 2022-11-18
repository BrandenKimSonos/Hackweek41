package benchmark

import (
	"context"
	"database/sql"
	"fmt"
	"hackweek41/utils"
	"log"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

var NUM_ITERATIONS = 100000

var db_user = utils.GoDotEnvVariable("DB_USER")
var db_pass = utils.GoDotEnvVariable("DB_PASSWORD")
var db_port = utils.GoDotEnvVariable("DB_PORT")
var db_server = utils.GoDotEnvVariable("DB_SERVER")

func queryDBConnection(db_conn *sql.DB) {
	record_rows, query_err := db_conn.Query("SELECT userid, usertype_id FROM www.dbo.users WHERE userid = 109272176")

	if query_err != nil {
		log.Println(query_err)
	}

	for record_rows.Next() {
		var (
			userid int
			usertype_id int
		)
		err := record_rows.Scan(
			&userid, 
			&usertype_id,
		)

		if err != nil {
			log.Fatal(err)
		}
		
		fmt.Println("USERID: ", userid)
		fmt.Println("USERTYPE_ID: ", usertype_id)
	}
}

func BlowingUpDBConnections() {
	fmt.Println("Blowing Up DB Connections...")
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
		log.Println(db_err)
	}
	
	for i := 0; i < NUM_ITERATIONS; i++ {
		
		if i % 1000 == 0 {
			fmt.Println("Finished " + strconv.Itoa(i) + " iterations out of " + strconv.Itoa(NUM_ITERATIONS) + "...")
			fmt.Println("Number of active connections: " + strconv.Itoa(db.Stats().OpenConnections))
		}
		go queryDBConnection(db)
	}
}

func queryRedisConnection(rdb *redis.Client) {
	for {
		_, rdb_err := rdb.Get("userid:109272176").Result()
		if rdb_err != nil {
			log.Fatal(rdb_err)
		}
		// fmt.Println("UserType Value: ", val)
	}
}

func BlowingUpRedisConnections() {
	fmt.Println("Blowing Up Redis Connections...")

	redis_url := utils.GoDotEnvVariable("REDIS_READ_URL")

	rdb := redis.NewClient(&redis.Options{
		Addr:     redis_url,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	for i := 0; i < NUM_ITERATIONS; i++ {
		if i % 1000 == 0 {
			fmt.Println("Finished " + strconv.Itoa(i) + " iterations out of " + strconv.Itoa(NUM_ITERATIONS) + "...")
		}

		go queryRedisConnection(rdb)
	}
}

func SingleThreadDB(wg *sync.WaitGroup, iterations int) {
	write_file, write_file_err := os.Create("./SingleThreadedDB.csv")

	if write_file_err != nil {
		log.Fatal(write_file_err)
	}
	defer write_file.Close()

	_, n1_err := write_file.WriteString("Latency (ms)\n")

	if n1_err != nil {
		log.Fatal(n1_err)
	}

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
		log.Println(db_err)
	}

	total_time := time.Now()
	var average_latency float64 = 0
	for i := 0; i < iterations; i++ {

		if i % 100 == 0 {
			fmt.Println("[DB]: Finished " + strconv.Itoa(i) + " iterations out of " + strconv.Itoa(iterations) + "...")
		}

		latency := time.Now()
		record_rows, query_err := db.Query("SELECT userid, usertype_id FROM www.dbo.users WHERE userid = 109272176")

		if query_err != nil {
			log.Println(query_err)
		}

		for record_rows.Next() {
			var (
				userid int
				usertype_id int
			)
			err := record_rows.Scan(
				&userid, 
				&usertype_id,
			)

			if err != nil {
				log.Fatal(err)
			}
		}
		latency_duration := time.Since(latency)

		_, n1_err := write_file.WriteString(strconv.FormatInt(latency_duration.Milliseconds(), 10)+"\n")

		if n1_err != nil {
			log.Fatal(n1_err)
		}

		average_latency += float64(latency_duration.Milliseconds())
	}
	average_latency /= float64(iterations)

	total_time_duration := time.Since(total_time)

	fmt.Println("[DB]: Total Time for ", iterations, " iterations: ", strconv.FormatInt(total_time_duration.Milliseconds(), 10), "ms")
	fmt.Println("[DB]: Average Latency for ", iterations, " iterations: ", average_latency, "ms")

	// file, err := os.Open("./data.txt")
    // if err != nil {
    //     log.Fatal(err)
    // }
    // defer file.Close()

    // scanner := bufio.NewScanner(file)
    // for scanner.Scan() {
    //     fmt.Println(scanner.Text())
    // }

    // if err := scanner.Err(); err != nil {
    //     log.Fatal(err)
    // }

	wg.Done()
}

func SingleThreadRedis(wg *sync.WaitGroup, iterations int) {
	// read_file, read_file_err := os.Open("./data.txt")
    // if read_file_err != nil {
    //     log.Fatal(read_file_err)
    // }
    // defer read_file.Close()

	write_file, write_file_err := os.Create("./SingleThreadedRedis.csv")

	if write_file_err != nil {
		log.Fatal(write_file_err)
	}
	defer write_file.Close()

	_, n1_err := write_file.WriteString("Latency (ms)\n")

	if n1_err != nil {
		log.Fatal(n1_err)
	}

	redis_url := utils.GoDotEnvVariable("REDIS_READ_URL")
	rdb := redis.NewClient(&redis.Options{
		Addr:     redis_url,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	total_time := time.Now()
	var average_latency float64 = 0
	for i := 0; i < iterations; i++ {

		if i % 100 == 0 {
			fmt.Println("[REDIS]: Finished " + strconv.Itoa(i) + " iterations out of " + strconv.Itoa(iterations) + "...")
		}

		latency := time.Now()
		_, rdb_err := rdb.Get("userid:109272176").Result()
		latency_duration := time.Since(latency)
		if rdb_err != nil {
			log.Fatal(rdb_err)
		}

		_, n1_err := write_file.WriteString(strconv.FormatInt(latency_duration.Milliseconds(), 10)+"\n")

		if n1_err != nil {
			log.Fatal(n1_err)
		}

		average_latency += float64(latency_duration.Milliseconds())
	}
	average_latency /= float64(iterations)

	total_time_duration := time.Since(total_time)

	fmt.Println("[REDIS]: Total Time for ", iterations, " iterations: ", strconv.FormatInt(total_time_duration.Milliseconds(), 10), "ms")
	fmt.Println("[REDIS]: Average Latency for ", iterations, " iterations: ", average_latency, "ms")

    // scanner := bufio.NewScanner(read_file)
    // for scanner.Scan() {
	// 	var latency float64 = 0
    //     fmt.Println(scanner.Text())
    // }

    // if err := scanner.Err(); err != nil {
    //     log.Fatal(err)
    // }

	wg.Done()
}

func SingleRoutineQuery() {
	fmt.Println("Testing Latency Using Only One Connection...")
	fmt.Println("Version", runtime.Version())
    fmt.Println("NumCPU", runtime.NumCPU())
    fmt.Println("GOMAXPROCS", runtime.GOMAXPROCS(0))

	var wg sync.WaitGroup
	wg.Add(2)

	go SingleThreadDB(&wg, 1000)
	go SingleThreadRedis(&wg, 1000)

	wg.Wait()
}

func runMultiThreadedDBQuery(db *sql.DB, write_file *os.File, average_latency *float64, db_wg *sync.WaitGroup) {
	defer db_wg.Done()
	latency := time.Now()
	record_rows, query_err := db.Query("SELECT userid, usertype_id FROM www.dbo.users WHERE userid = 109272176")

	if query_err != nil {
		log.Println(query_err)
	}

	for record_rows.Next() {
		var (
			userid int
			usertype_id int
		)
		err := record_rows.Scan(
			&userid, 
			&usertype_id,
		)

		if err != nil {
			log.Fatal(err)
		}
	}
	latency_duration := time.Since(latency)

	_, n1_err := write_file.WriteString(strconv.FormatInt(latency_duration.Milliseconds(), 10)+"\n")

	if n1_err != nil {
		log.Fatal(n1_err)
	}

	*average_latency += float64(latency_duration.Milliseconds())
}

func MultiThreadDB(wg *sync.WaitGroup, iterations int) {
	write_file, write_file_err := os.Create("./MultiThreadedDB.csv")

	if write_file_err != nil {
		log.Fatal(write_file_err)
	}
	defer write_file.Close()

	_, n1_err := write_file.WriteString("Latency (ms)\n")

	if n1_err != nil {
		log.Fatal(n1_err)
	}

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
		log.Println(db_err)
	}

	total_time := time.Now()
	var average_latency float64 = 0
	var db_wg sync.WaitGroup
	db_wg.Add(iterations)
	for i := 0; i < iterations; i++ {

		if i % 1000 == 0 {
			fmt.Println("[DB]: Finished " + strconv.Itoa(i) + " iterations out of " + strconv.Itoa(iterations) + "...")
		}

		go runMultiThreadedDBQuery(db, write_file, &average_latency, &db_wg)
	}
	db_wg.Wait()
	average_latency /= float64(iterations)

	total_time_duration := time.Since(total_time)

	fmt.Println("[DB]: Total Time for ", iterations, " iterations: ", strconv.FormatInt(total_time_duration.Milliseconds(), 10), "ms")
	fmt.Println("[DB]: Average Latency for ", iterations, " iterations: ", average_latency, "ms")
}

func runMultiThreadedRedisQuery(rdb *redis.Client, write_file *os.File, average_latency *float64, redis_wg *sync.WaitGroup, i int) {
	defer redis_wg.Done()
	latency := time.Now()
	_, rdb_err := rdb.Get("userid:109272176").Result()
	latency_duration := time.Since(latency)
	if rdb_err != nil {
		log.Fatal(rdb_err)
	}
	// fmt.Println("Iteration: ", i, " Value: ", val)
	_, n1_err := write_file.WriteString(strconv.FormatInt(latency_duration.Milliseconds(), 10)+"\n")

	if n1_err != nil {
		log.Fatal(n1_err)
	}

	*average_latency += float64(latency_duration.Milliseconds())
}

func MultiThreadRedis(wg *sync.WaitGroup, iterations int) {
	write_file, write_file_err := os.Create("./MultiThreadedRedis.csv")

	if write_file_err != nil {
		log.Fatal(write_file_err)
	}
	defer write_file.Close()

	_, n1_err := write_file.WriteString("Latency (ms)\n")

	if n1_err != nil {
		log.Fatal(n1_err)
	}

	redis_url := utils.GoDotEnvVariable("REDIS_READ_URL")
	rdb := redis.NewClient(&redis.Options{
		Addr:     redis_url,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	total_time := time.Now()
	var average_latency float64 = 0
	var redis_wg sync.WaitGroup
	redis_wg.Add(iterations)

	for i := 0; i < iterations; i++ {

		if i % 1000 == 0 {
			fmt.Println("[REDIS]: Finished " + strconv.Itoa(i) + " iterations out of " + strconv.Itoa(iterations) + "...")
		}

		go runMultiThreadedRedisQuery(rdb, write_file, &average_latency, &redis_wg, i)
	}

	redis_wg.Wait()
	average_latency /= float64(iterations)

	total_time_duration := time.Since(total_time)

	fmt.Println("[REDIS]: Total Time for ", iterations, " iterations: ", strconv.FormatInt(total_time_duration.Milliseconds(), 10), "ms")
	fmt.Println("[REDIS]: Average Latency for ", iterations, " iterations: ", average_latency, "ms")

	wg.Done()
}

func MultiThreadedQuery() {
	fmt.Println("Testing Latency Using Only Multi Connection...")
	fmt.Println("Version", runtime.Version())
    fmt.Println("NumCPU", runtime.NumCPU())
    fmt.Println("GOMAXPROCS", runtime.GOMAXPROCS(0))

	var wg sync.WaitGroup
	wg.Add(1)

	go MultiThreadDB(&wg, 1000)
	go MultiThreadRedis(&wg, 1000)

	wg.Wait()
}

func dbConcurrencyQuery(db *sql.DB, write_file *os.File, db_wg *sync.WaitGroup, iterations int, index int) {
	defer db_wg.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	number_of_successes := 0
	number_of_failures := 0
	for i := 0; i < iterations; i++ {

		if i % 10 == 0 {
			fmt.Println("Running iteration ", i, " out of ", iterations)
		}

		record_rows, query_err := db.QueryContext(ctx, "SELECT userid, usertype_id FROM www.dbo.users WHERE userid = 109272176")

		if query_err != nil {
			log.Println(query_err)
			number_of_failures++
			continue
		}

		for record_rows.Next() {
			var (
				userid int
				usertype_id int
			)
			err := record_rows.Scan(
				&userid, 
				&usertype_id,
			)

			if err != nil {
				number_of_failures++
				break
			} else {
				number_of_successes++
			}
		}
	}

	write_file.WriteString("[DB]: Iteration " + strconv.Itoa(index) + " Successes: " + strconv.Itoa(number_of_successes) + " Failures: " +strconv.Itoa(number_of_failures) + "\n")
}

func runDBConcurrencyTest(n int, iterations int) {
	write_file, write_file_err := os.Create("./DBConcurrencyTest.txt")

	if write_file_err != nil {
		log.Fatal(write_file_err)
	}
	defer write_file.Close()

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
		log.Println(db_err)
	}
	var db_wg sync.WaitGroup
	db_wg.Add(n)

	for i := 0; i < n; i++ {
		fmt.Println("Running goroutine number: ", i)
		go dbConcurrencyQuery(db, write_file, &db_wg, iterations, i)
	}

	db_wg.Wait()
}

func redisConcurrencyQuery(rdb *redis.Client, write_file *os.File, redis_wg *sync.WaitGroup, iterations int, index int) {
	defer redis_wg.Done()
	
	number_of_successes := 0
	number_of_failures := 0
	for i := 0; i < iterations; i++ {

		if i % 10 == 0 {
			fmt.Println("Running iteration ", i, " out of ", iterations)
		}

		val, rdb_err := rdb.Get("userid:109272176").Result()
		if rdb_err != nil {
			number_of_failures += 1
			continue
		}

		if val != string(redis.Nil) {
			number_of_successes += 1
			continue
		}
	}
	
	write_file.WriteString("[REDIS]: Iteration " + strconv.Itoa(index) + " Successes: " + strconv.Itoa(number_of_successes) + " Failures: " +strconv.Itoa(number_of_failures) + "\n")
}

func runRedisConcurrencyTest(n int, iterations int) {
	write_file, write_file_err := os.Create("./RedisConcurrencyTest.txt")

	if write_file_err != nil {
		log.Fatal(write_file_err)
	}
	defer write_file.Close()

	redis_url := utils.GoDotEnvVariable("REDIS_READ_URL")
	rdb := redis.NewClient(&redis.Options{
		Addr:     redis_url,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rdb = rdb.WithContext(ctx)

	var redis_wg sync.WaitGroup
	redis_wg.Add(n)

	for i := 0; i < n; i++ {
		fmt.Println("Running goroutine number: ", i)
		go redisConcurrencyQuery(rdb, write_file, &redis_wg, iterations, i)
	}

	redis_wg.Wait()
}

func TestConcurrencyLimits() {
	fmt.Println("Testing the Concurrency Limits...")
	fmt.Println("Version", runtime.Version())
    fmt.Println("NumCPU", runtime.NumCPU())
    fmt.Println("GOMAXPROCS", runtime.GOMAXPROCS(0))

	number_of_goroutines := 10
	number_of_iterations := 1000
	runDBConcurrencyTest(number_of_goroutines, number_of_iterations)
	runRedisConcurrencyTest(number_of_goroutines, number_of_iterations)
}