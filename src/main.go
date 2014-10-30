package main

import (
	"database/sql"
	_ "fmt"
	"github.com/PuerkitoBio/goquery"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strings"
)

func ExampleScrape() []string {
	var result []string
	doc, err := goquery.NewDocument("http://video.56.com/opera/22818.html")
	//13 doc, err := goquery.NewDocument("http://video.56.com/opera/6913.html")
	//12 doc, err := goquery.NewDocument("http://video.56.com/opera/6361.html")
	//11 doc, err := goquery.NewDocument("http://video.56.com/opera/6324.html")

	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".episode_cnt a").Each(func(i int, s *goquery.Selection) {
		title := s.Find("span").Text()
		// fmt.Println(s.Get(0).Attr[2].Val, "-", title)
		result = append(result, s.Get(0).Attr[2].Val+" "+title)
	})
	return result
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	db, err := sql.Open("sqlite3", "./data.db")
	checkErr(err)

	//插入数据
	stmt, err := db.Prepare("INSERT INTO MusicTea_2014(title, url, time, text) values(?,?,?,?)")
	checkErr(err)
	content := ExampleScrape()

	var (
		title string
		url   string
		time  string
		text  string
	)

	for i := 0; i < len(content); i++ {
		if len(strings.Fields(content[i])) == 4 {
			title = strings.Fields(content[i])[1] + " " + strings.Fields(content[i])[2] + " " + strings.Fields(content[i])[3]
			url = strings.Fields(content[i])[0]
			time = strings.Fields(content[i])[3]
			text = strings.Fields(content[i])[1]
		} else {
			if len(strings.Fields(content[i])) == 3 {
				title = strings.Fields(content[i])[1] + " " + strings.Fields(content[i])[2]
				url = strings.Fields(content[i])[0]
				time = strings.Fields(content[i])[2]
				text = strings.Fields(content[i])[1]
			} else {
				title = strings.Fields(content[i])[1]
				url = strings.Fields(content[i])[0]
				time = strings.Fields(content[i])[1]
				text = strings.Fields(content[i])[1]
			}
		}
		_, err := stmt.Exec(title, url, time, text)
		checkErr(err)
	}

	// id, err := res.LastInsertId()
	// checkErr(err)

	// fmt.Println(id)

	defer db.Close()
}
