package main

import (
	"flag"
	"fmt"
	"log"
	"github.com/gocolly/colly"
)

func main() {
	maxFlag := flag.Int("max", 10, "Max number of comments to show")
	flag.Parse()

	c := colly.NewCollector()

	counter := 0
	max := *maxFlag

	c.OnHTML(".push", func(e *colly.HTMLElement) {

		if counter >= max {
			return
		}

		userID := e.ChildText(".push-userid")
		content := e.ChildText(".push-content")
		ipdatetime := e.ChildText(".push-ipdatetime")

		fmt.Printf("%d. 名字：%s，留言%s，時間： %s\n", counter+1, userID, content, ipdatetime)

		counter++
	})


	err := c.Visit("https://www.ptt.cc/bbs/joke/M.1481217639.A.4DF.html")
	if err != nil {
		log.Fatalf("Failed to visit page: %v", err)
	}
}