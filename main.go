package main

import (
	"encoding/json"
	"os"
	"fmt"
	"time"

	"github.com/mmcdole/gofeed"
	"log"
)

type Feed struct {
	Name    string
	Url     string
}

type Source struct {
	List    []*Feed
}

type Item struct {
	Title   string
	Time    time.Time
	Url     string
	Content string
}

type News struct {
	List []*Item
}

func InitMinSources() {
	f1 := &Feed{"lenta.ru", "https://lenta.ru/l/r/EX/import.rss"}
	f2 := &Feed{"yandex.ru", "https://news.yandex.ru/index.rss"}
	f3 := &Feed{"echo.msk.ru", "http://echo.msk.ru/news.rss"}

	source := new(Source)
	source.List = append(source.List, f1)
	source.List = append(source.List, f2)
	source.List = append(source.List, f3)

	file, _ := os.OpenFile("feed.json", os.O_CREATE|os.O_WRONLY, 0777)
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	enc.Encode(source)
}

func ReadListSources() *Source {
	file, _ := os.Open("feed.json")
	defer file.Close()

	sl := new(Source)

	dec := json.NewDecoder(file)
	dec.Decode(sl)

	return sl
}

func ReadNews(sourceFeed *Feed, news *News) *News {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(sourceFeed.Url)
	if err != nil {
		log.Fatal("Unable to get the feed. Error is ", err)
	}

	newsItems := feed.Items

	for _, newsItem := range newsItems {
		item := new(Item)
		item.Title = newsItem.Title
		item.Url = newsItem.Link
		item.Time = *newsItem.PublishedParsed
		item.Content = newsItem.Content

		news.List = append(news.List, item)
	}
	return news

}

func main() {
	if _, err := os.Stat("feed.json"); os.IsNotExist(err) {
		InitMinSources()
	}

	news := new(News)

	source := ReadListSources()
	for _, i := range source.List {
		_ = ReadNews(i, news)
	}

	file, _ := os.OpenFile("news.json", os.O_CREATE|os.O_WRONLY, 0777)
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	enc.Encode(news)

	fmt.Println(news)
}