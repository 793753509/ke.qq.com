package scrape

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"ke.qq.com/model"
	"strings"
	"time"
)

type Scraper interface {
	ScrapeCourseInfo(tableName string)
}

type ScrapeManagerInst struct {
	MsgChan chan model.CourseInfo
}

func NewScrapeManager() (Scraper, chan model.CourseInfo) {
	curseChan := make(chan model.CourseInfo, 10)
	return &ScrapeManagerInst{
		MsgChan: curseChan,
	}, curseChan
}

var typeFilter = map[string]struct{}{
	"IT·互联网": struct{}{},
	"兴趣·生活":  struct{}{},
	"升学·考研":  struct{}{},
	"电商·营销":  struct{}{},
	"职业·考证":  struct{}{},
	"设计·创作":  struct{}{},
	"语言·留学":  struct{}{},
}

func (s *ScrapeManagerInst) ScrapeCourseInfo(tableName string) {
	t := time.Now()
	number := 1
	c := colly.NewCollector(func(c *colly.Collector) {
		colly.MaxDepth(0)
		extensions.RandomUserAgent(c)
		c.Async = true
		c.AllowURLRevisit = false
		err := c.Limit(&colly.LimitRule{
			DomainRegexp: "ke.qq.com",
			RandomDelay:  2 * time.Second,
			Parallelism:  4,
		})
		if err != nil {
			fmt.Println(err.Error())
		}

	})
	c.OnHTML("div.sort-menu-border1 dl.sort-menu.sort-menu1.clearfix", func(e *colly.HTMLElement) {
		e.ForEach("dd a", func(_ int, h *colly.HTMLElement) {
			link := h.Attr("href")
			//fmt.Printf("find type: %s\n",e.Request.AbsoluteURL(link))
			c.Visit(h.Request.AbsoluteURL(link))

		})
	})

	c.OnHTML("div.sort-page a.page-next-btn.icon-font.i-v-right", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		c.Visit(e.Request.AbsoluteURL(link))

	})

	c.OnHTML("div.main-left ul.course-card-list a", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		c.Visit(e.Request.AbsoluteURL(link))

	})
	courseInfo := model.CourseInfo{
		TableName: tableName,
	}

	c.OnHTML("body.page-course.l-media.l-aside-right", func(e *colly.HTMLElement) {
		e.DOM.Each(func(i int, selection *goquery.Selection) {

			courseInfo.Name = strings.TrimSpace(selection.Find("a.breadcrumb-item.item--tt").Last().Text())
			url, _ := selection.Find("a.breadcrumb-item.item--tt").Attr("href")
			tmp := strings.Split(url, "/")
			if len(tmp) > 1 {
				courseInfo.ID = tmp[len(tmp)-1]
			}
			t := selection.Find("a.tt-link.js-teacher-name")
			teachers := strings.TrimSpace(t.First().Text())
			for i := 1; i < len(t.Nodes); i++ {
				teachers = teachers + "," + strings.TrimSpace(t.Eq(i).Text())
			}
			if teachers == "" {
				teachers = "无"
			}
			courseInfo.Teachers = teachers

			price := strings.TrimSpace(selection.Find("span.price.custom-string").First().Text())
			if price == "" {
				price = "免费"
			}
			courseInfo.Price = price
			courseType := selection.Find("nav.breadcrumb.inner-center").ChildrenFiltered("a.breadcrumb-item").Eq(1).Text()
			if _, ok := typeFilter[courseType]; !ok {
				courseType = "其他"
			}
			courseInfo.CourseType = courseType
			fmt.Printf("%d --> %s, %s, %s, %s, %s \n", number, courseInfo.Name, courseInfo.ID, courseInfo.Teachers, courseInfo.Price, courseInfo.CourseType)
			number++
			s.MsgChan <- courseInfo

		})
	})
	c.OnError(func(response *colly.Response, err error) {
		fmt.Println(err)
	})
	c.Visit("https://ke.qq.com/course/list")

	c.Wait()
	close(s.MsgChan)
	fmt.Printf("花费时间：%s", time.Since(t))
}
