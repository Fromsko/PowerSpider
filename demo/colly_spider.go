package main

import (
	"PowerSpider/conduit"
	"PowerSpider/config"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
)

var conf config.Config

type Spider struct {
	name    string
	spider  *colly.Collector
	cookies []*http.Cookie
	params  map[string]string
	data    string
}

// 初始化任务
func InitSpider() (c *colly.Collector) {
	conf = config.Config{
		User: "202127530334",
		Pwd:  "102018",
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
	c.spider.OnHTML(`input`, func(e *colly.HTMLElement) {
		ID := e.Attr("id")
		switch ID {
		case "__VIEWSTATE":
			c.params["__VIEWSTATE"] = e.Attr("value")
		case "__VIEWSTATEGENERATOR":
			c.params["__VIEWSTATEGENERATOR"] = e.Attr("value")
		case "__EVENTVALIDATION":
			c.params["__EVENTVALIDATION"] = e.Attr("value")
		default:
			break
		}
	})

	c.spider.OnResponse(func(r *colly.Response) {
		if cookieJar := c.spider.Cookies(fetchUrl); len(cookieJar) > 0 {
			fmt.Printf(
				"[Spider-{%s}-Run] |Cookies| Get success => %v\n",
				c.name, cookieJar,
			)
			c.cookies = cookieJar
		}
	})

	c.spider.OnScraped(func(r *colly.Response) {
		var encoded string
		fmt.Printf(
			"[Spider-{%s}-Stop] Task is finished!\n",
			c.name,
		)
		params := url.Values{
			"__LASTFOCUS":          {""},
			"__EVENTTARGET":        {"UserLogin$ImageButton1"},
			"__EVENTARGUMENT":      {""},
			"__VIEWSTATE":          {c.params["__VIEWSTATE"]},
			"__VIEWSTATEGENERATOR": {c.params["__VIEWSTATEGENERATOR"]},
			"__EVENTVALIDATION":    {c.params["__EVENTVALIDATION"]},
			"UserLogin:txtUser":    {conf.User},
			"UserLogin:txtPwd":     {conf.Pwd},
		}

		for key, values := range params {
			for _, value := range values {
				encoded += key + "=" + value + "&"
			}
		}

		c.data = encoded[:len(encoded)-1] + "&UserLogin%3AddlPerson=%BF%A8%BB%A7&UserLogin%3AtxtSure="
	})

	c.spider.Visit(fetchUrl)
}

func (c *Spider) FinishTask(verfity string) {
	c.spider.OnScraped(func(r *colly.Response) {
		fmt.Printf(
			"[Spider-{%s}-Stop] Task is finished!\n",
			c.name,
		)
	})
}

func (c *Spider) VerifyCode(fetchUrl string) {
	c.spider.SetCookies(fetchUrl, c.cookies)

	c.spider.OnResponse(func(r *colly.Response) {
		if r.StatusCode == http.StatusOK {
			saveImage(r.Body, "info")
		}

		result, err := conduit.OcrResult([]byte(r.Body))
		if err != nil {
			fmt.Println("error:", err)
			return
		}
		fmt.Println("识别结果是 => ", result)
	})

	c.spider.Visit(fetchUrl)
}

func saveImage(content []byte, url string) {
	fileName := filepath.Base(url)
	filePath := filepath.Join("res/imgs", fileName+".jpg")

	// 创建目录
	err := os.MkdirAll("res/imgs", os.ModePerm)
	if err != nil {
		log.Printf("创建目录失败：%s", err)
		return
	}

	// 创建文件
	file, err := os.Create(filePath)
	if err != nil {
		log.Printf("创建文件失败：%s", err)
		return
	}
	defer file.Close()

	// 保存图片内容
	_, err = file.Write(content)
	if err != nil {
		log.Printf("保存图片失败：%s", err)
	} else {
		fmt.Printf("已下载并保存图片：%s\n", filePath)
	}
}

func main() {
	HomeSpider := &Spider{
		name:   "主页爬虫",
		spider: InitSpider(),
		params: make(map[string]string, 10),
	}
	HomeSpider.SetHeaders()
	HomeSpider.VisitHomePage(conf.BaseUrl)

	VerifySpider := &Spider{
		name:    "验证码识别",
		spider:  InitSpider(),
		data:    HomeSpider.data,
		cookies: HomeSpider.cookies,
	}
	VerifySpider.SetHeaders(map[string]string{
		"Accept":          "q=0.9,image/webp,image/apng,*/*;",
		"Accept-Encoding": "gzip, deflate",
		"Content-Type":    "application/x-www-form-urlencoded",
		"accept-language": "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6",
		"cache-control":   "max-age=0",
	})
	VerifySpider.VerifyCode(config.VerifyUrl)
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
