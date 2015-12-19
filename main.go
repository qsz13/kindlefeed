package main

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/robfig/cron.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	config := Config{}
	err := decoder.Decode(&config)
	if err != nil {
		log.Fatal(err)
	}
	//	fmt.Println(config.FeedURL)
	c := cron.New()
	c.AddFunc(config.Cron, func() {
		fetchFeed(&config)
	})
	c.Start()
	//	fetchFeed(&config)

	select {}

}

func fetchFeed(config *Config) {

	dbconfig := config.Database
	datasource := fmt.Sprintf("%s:%s@%s(%s)/%s", dbconfig.Username, dbconfig.Password, dbconfig.Protocol, dbconfig.Address, dbconfig.DBname)
	db, err := sql.Open("mysql", datasource)
	if err != nil {
		log.Fatal("Database connection failure:\n", err)
		return
	}
	defer db.Close()
	fmt.Println(datasource)
	for _, url := range config.FeedURL {
		log.Println("fetching feed: " + url)
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("failed to download feed:\n", err)
			return
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		var feed Feed
		xml.Unmarshal(body, &feed)
		//	fmt.Println(feed.Channel.Title)
		//	for _, item := range feed.Channel.Items {
		//		fmt.Println(item.Title)
		//	}
		saveFeedToDB(&feed, db)
	}
}

func saveFeedToDB(feed *Feed, db *sql.DB) {
	fmt.Println("begin to save")
	for _, item := range feed.Channel.Items {
		stmt, err := db.Prepare(
			"INSERT INTO feed (title, link, guid, creator, pubDate, description, content, create_time, update_time) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			log.Fatal("Insert record failed:\n", err, item)
			return
		}
		pubDate, err := time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", item.PubDate)
		_, err = stmt.Exec(
			item.Title,
			item.Link,
			item.Guid,
			item.Creator,
			pubDate,
			item.Description,
			item.Content,
			time.Now(),
			time.Now(),
		)
		if err != nil && strings.Contains(err.Error(), "1062") {
			//log.Println("duplicated guid: ", item.Guid)
			continue
		} else if err != nil {
			log.Println("Insert data failed:", err, item)
		}

		log.Println("data inserted: ", item.Title)
	}

}
