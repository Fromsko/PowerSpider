package main

import (
	"PowerSpider/config"
	"fmt"
	"net/http"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
)

var conf config.Config

type Spider struct {
	name    string
	spider  *colly.Collector
	cookies []*http.Cookie
}

// 初始化任务
func InitSpider() (c *colly.Collector) {
	conf = config.Config{
		Timer: config.Timer{
			TimeUnit: "hourse",
			TimeInfo: 2,
		},
		ResDir: "res",
		Porxy:  "http://localhost:7890",
	}
	config.InitConfig(&conf)

	c = colly.NewCollector(
		// 允许重复访问
		colly.AllowURLRevisit(),
		// 异步请求
		colly.Async(false),
		// 设置请求头
		colly.UserAgent(config.UserAgent),
	)

	// 设置代理地址
	c.SetProxy(conf.Porxy)

	// 随机UserAgent
	extensions.RandomUserAgent(c)

	if err := c.Limit(&colly.LimitRule{
		// Filter domains affected by this rule
		// 筛选受此规则影响的域
		DomainGlob: conf.BaseUrl + "/*",
		// Set a delay between requests to these domains
		// 设置对这些域的请求之间的延迟
		Delay: 1 * time.Second,
		// Add an additional random delay
		// 添加额外的随机延迟
		RandomDelay: 1 * time.Second,
		// 设置并发
		Parallelism: 3,
	}); err != nil {
		fmt.Println(err)
	}

	// 限制采集规则
	/*
		在Colly里面非常方便控制并发度，只抓取符合某个(些)规则的URLS
		colly.LimitRule{DomainGlob: "*.douban.*", Parallelism: 5}，
		表示限制只抓取域名是douban(域名后缀和二级域名不限制)的地址，
		当然还支持正则匹配某些符合的 URLS

		Limit方法中也限制了并发是5。为什么要控制并发度呢？
		因为抓取的瓶颈往往来自对方网站的抓取频率的限制，
		如果在一段时间内达到某个抓取频率很容易被封，所以我们要控制抓取的频率。
		另外为了不给对方网站带来额外的压力和资源消耗，也应该控制你的抓取机制。
	*/
	return c
}

// 设置 headers
func (c *Spider) SetHeaders(header ...map[string]string) {
	var h map[string]string

	if len(header) > 0 {
		h = header[0]
	}

	if h == nil {
		h = map[string]string{
			"Accept":       "*/*",
			"Content-Type": "application/x-www-form-urlencoded",
			"Origin":       "http://10.14.0.124",
			"User-Agent":   config.UserAgent,
		}
	}

	c.spider.OnRequest(func(r *colly.Request) {
		for key, value := range h {
			r.Headers.Add(key, value)
		}
		fmt.Printf(
			"[Spider-{%s}-Run] |Visit| Page => %s \n",
			c.name, r.URL.String(),
		)
	})
}

func (c *Spider) VisitHomePage(fetchUrl string) {
	c.spider.OnResponse(func(r *colly.Response) {
		if cookieJar := c.spider.Cookies(fetchUrl); len(cookieJar) > 0 {
			fmt.Printf(
				"[Spider-{%s}-Run] |Cookies| Get success => %v\n",
				c.name, cookieJar,
			)
			c.cookies = cookieJar
		}
	})
	c.spider.Visit(fetchUrl)
}

func (c *Spider) GetParams() {
	c.spider.OnHTML(`#imgAuthCode`, func(e *colly.HTMLElement) {
		fmt.Println(e.Attr("src"))
	})

	c.spider.OnHTML(`input`, func(e *colly.HTMLElement) {
		// e.Request.Visit(e.Attr("href"))
		// fmt.Println(e.Attr("id"))

		if e.Attr("id") == "__VIEWSTATE" {
			fmt.Println(e.Attr("value"))
		}
		fmt.Println(e.Attr("src"))
	})
}

func (c *Spider) FinishTask() {
	c.spider.OnScraped(func(r *colly.Response) {
		fmt.Printf(
			"[Spider-{%s}-Stop] Task is finished!\n",
			c.name,
		)
	})
}

func main() {
	HomeSpider := &Spider{
		name:   "主页爬虫",
		spider: InitSpider(),
	}
	HomeSpider.SetHeaders()
	HomeSpider.GetParams()
	HomeSpider.FinishTask()
	HomeSpider.VisitHomePage(conf.BaseUrl)
}

/*
请求执行之前调用
	- OnRequest
响应返回之后调用
	- OnResponse
监听执行 selector
	- OnHTML
监听执行 selector
	- OnXML
错误回调
	- OnError
完成抓取后执行，完成所有工作后执行
	- OnScraped
取消监听，参数为 selector 字符串
	- OnHTMLDetach
取消监听，参数为 selector 字符串
	- OnXMLDetach
*/
