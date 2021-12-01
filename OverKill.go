package OverKill

import (
	"bufio"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"time"
)

var ctx = context.Background()

func LoadClient(filename string) {
	if _, err := os.Stat("redis.lock"); err == nil {
		println("Redis lock file detected... Continue? \n y/n")
		var continueCheck string
		fmt.Scanf("%s", &continueCheck)
		if continueCheck != "y"{
			os.Exit(0)
		}

	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	n := 0
	for scanner.Scan() {
		rdb.Set(ctx, string(n), scanner.Text(),0)
		n = n + 1
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	ioutil.WriteFile("redis.lock", []byte("Time DB Loaded - " + time.Now().Format(time.UnixDate)),0644)
}

func scanRedis(hashCheck string){
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	var d []string
	var cursor uint64
	var n int

	CheckLoop:
		for {
			var keys []string
			var err error
			keys, cursor, err = rdb.Scan(ctx, cursor, "*", 100).Result()
			if err != nil {
				panic(err)
			}
			keyValues := rdb.MGet(ctx, keys...)
			keyValues.Scan(&d)
			for _, entry := range d{
				if matched, _ := regexp.MatchString(".*" + hashCheck + ".*", entry); matched == true{
					println("Found Match")
					println(entry)
					break CheckLoop
				}
			}
			n += len(keys)
			if cursor == 0 {
				break CheckLoop
			}
		}
}

func main(){
	LoadClient("list.nsrl")

	file, err := os.Open("hash.list")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		scanRedis(scanner.Text())
	}

}