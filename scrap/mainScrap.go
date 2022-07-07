package scrap

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"golang.org/x/net/html"
	"golang.org/x/net/html/charset"
)

var (
	/*
								URI中的#符号代表的是fragment,
						如:

							https://domain/index#L18

							　　这个 URI 中 L18 就是 fragment 的值。这有哪些特殊的地方呢？

							1. # 有别于 ?，? 后面的查询字符串会被网络请求带上服务器，而 fragment 不会被发送的服务器；
							2. fragment 的改变不会触发浏览器刷新页面，但是会生成浏览历史；
							3. fragment 会被浏览器根据文件媒体类型（MIME type）进行对应的处理；
							4. Google 的搜索引擎会忽略 # 及其后面的字符串。

				利用1,2特性,使得前端js可以知道位置hash值的变化,实现页面的滚动或跳转,从而实现了一个简单的页面路由

				利用特新3,页面内定位。在页面中通过设置标签的 id 属性来定义锚点，从而实现锚点定位。实际上锚点定位的实现正是依赖了 fragment 的特性 3。
			如这个 URI https://domain/index.html#L18，假设返回的文件类型是 text/html，则浏览器会读取 URI’s fragment，然后在页面中寻找 L18 这个锚点，并将页面滚动到该锚点的位置

				利用特性4,因为 fragment 会被 Google 搜索引擎忽略掉，因此对于用 hash 模式前端路由的应用的 SEO 来说是很不友好的。
		不过 Google 给了一个方案，就是在 # 紧跟一个 ! ，这样Google 搜索引擎就会将这个 URI 进行转换，
		如 https://domain/index.html#!L18转换后就成为了 https://domain/index.html?_escaped_fragment_=L18。
		这样搜索引擎就会携带上 URI’s fragment 直接去访问这个 URI，开发者可以利用这个 trick 优化网站的 SEO。

	*/
	EscapedFragment string = "_escaped_fragment_=" // #!将被转换为 _escaped_fragment_= ,目的是某些网站为google搜索引擎做的优化
	fragmentRegexp         = regexp.MustCompile("#!(.*)")
)

type Scraper struct {
	Url                *url.URL //用户输入的URL
	EscapedFragmentUrl *url.URL //TODO:WTF?
	MaxRedirect        int      //最大跳转次数?
}

// Document 结果
type Document struct {
	Body    bytes.Buffer
	Preview DocumentPreview
}

type DocumentPreview struct {
	Icon        string
	Name        string
	Title       string
	Description string
	Images      []string
	Link        string
}

func Scrape(uri string, maxRedirect int) (*Document, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	return (&Scraper{Url: u, MaxRedirect: maxRedirect}).Scrape()
}

// Scrape 起始入口
func (scraper *Scraper) Scrape() (*Document, error) {
	//开始爬取结果
	doc, err := scraper.getDocument()
	if err != nil {
		return nil, err
	}
	//将结果进行解析
	err = scraper.parseDocument(doc)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func (scraper *Scraper) getUrl() string {
	if scraper.EscapedFragmentUrl != nil { //如果有跳转地址,则返回跳转的地址
		return scraper.EscapedFragmentUrl.String()
	}
	return scraper.Url.String()
}

func (scraper *Scraper) toFragmentUrl() error {
	unescapedurl, err := url.QueryUnescape(scraper.Url.String())
	if err != nil {
		return err
	}
	matches := fragmentRegexp.FindStringSubmatch(unescapedurl)
	if len(matches) > 1 {
		escapedFragment := EscapedFragment
		for _, r := range matches[1] {
			b := byte(r)
			if avoidByte(b) {
				continue
			}
			if escapeByte(b) {
				escapedFragment += url.QueryEscape(string(r))
			} else {
				escapedFragment += string(r)
			}
		}

		p := "?"
		if len(scraper.Url.Query()) > 0 {
			p = "&"
		}
		fragmentUrl, err := url.Parse(strings.Replace(unescapedurl, matches[0], p+escapedFragment, 1))
		if err != nil {
			return err
		}
		scraper.EscapedFragmentUrl = fragmentUrl
	} else {
		p := "?"
		if len(scraper.Url.Query()) > 0 {
			p = "&"
		}
		fragmentUrl, err := url.Parse(unescapedurl + p + EscapedFragment)
		if err != nil {
			return err
		}
		scraper.EscapedFragmentUrl = fragmentUrl
	}
	return nil
}

func (scraper *Scraper) getDocument() (*Document, error) {
	//爬取一次,深度减一
	scraper.MaxRedirect -= 1
	//如果这uri包含fragment(一般是跳转时)
	if strings.Contains(scraper.Url.String(), "#!") {
		scraper.toFragmentUrl()
	}
	if strings.Contains(scraper.Url.String(), EscapedFragment) {
		scraper.EscapedFragmentUrl = scraper.Url
	}

	//构建请求
	tr := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: 3 * time.Second, //建立连接的超时
		}).DialContext,
		DisableKeepAlives: true, //告知服务端不启用长连接,避免占用端口资源
		MaxIdleConns:      100,
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true}, //跳过证书验证,确保自创证书的网站也能被收集
	}
	cli := http.Client{
		Transport: tr,
	}
	req, err := http.NewRequest("GET", scraper.getUrl(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", "GoScraper")

	resp, err := cli.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	if resp.Request.URL.String() != scraper.getUrl() {
		scraper.EscapedFragmentUrl = nil
		scraper.Url = resp.Request.URL
	}
	b, err := convertUTF8(resp.Body, resp.Header.Get("content-type"))
	if err != nil {
		return nil, err
	}
	doc := &Document{Body: b, Preview: DocumentPreview{Link: scraper.Url.String()}}

	return doc, nil
}

func convertUTF8(content io.Reader, contentType string) (bytes.Buffer, error) {
	buff := bytes.Buffer{}
	content, err := charset.NewReader(content, contentType)
	if err != nil {
		return buff, err
	}
	_, err = io.Copy(&buff, content)
	if err != nil {
		return buff, err
	}
	return buff, nil
}

func (scraper *Scraper) parseDocument(doc *Document) error {
	t := html.NewTokenizer(&doc.Body)
	var ogImage bool
	var headPassed bool
	var hasFragment bool
	var hasCanonical bool
	var canonicalUrl *url.URL
	doc.Preview.Images = []string{}
	// saves previews' link in case that <link rel="canonical"> is found after <meta property="og:url">
	link := doc.Preview.Link
	// set default value to site name if <meta property="og:site_name"> not found
	doc.Preview.Name = scraper.Url.Host
	// set default icon to web root if <link rel="icon" href="/favicon.ico"> not found
	doc.Preview.Icon = fmt.Sprintf("%s://%s%s", scraper.Url.Scheme, scraper.Url.Host, "/favicon.ico")
	for {
		tokenType := t.Next()
		if tokenType == html.ErrorToken {
			return nil
		}
		if tokenType != html.SelfClosingTagToken && tokenType != html.StartTagToken && tokenType != html.EndTagToken {
			continue
		}
		token := t.Token()

		switch token.Data {
		case "head":
			if tokenType == html.EndTagToken {
				headPassed = true
			}
		case "body":
			headPassed = true

		case "link":
			var canonical bool
			var hasIcon bool
			var href string
			for _, attr := range token.Attr {
				if cleanStr(attr.Key) == "rel" && cleanStr(attr.Val) == "canonical" {
					canonical = true
				}
				if cleanStr(attr.Key) == "rel" && strings.Contains(cleanStr(attr.Val), "icon") {
					hasIcon = true
				}
				if cleanStr(attr.Key) == "href" {
					href = attr.Val
				}
				if len(href) > 0 && canonical && link != href {
					hasCanonical = true
					var err error
					canonicalUrl, err = url.Parse(href)
					if err != nil {
						return err
					}
				}
				if len(href) > 0 && hasIcon {
					doc.Preview.Icon = href
				}
			}

		case "meta":
			if len(token.Attr) != 2 {
				break
			}
			if metaFragment(token) && scraper.EscapedFragmentUrl == nil {
				hasFragment = true
			}
			var property string
			var content string
			for _, attr := range token.Attr {
				if cleanStr(attr.Key) == "property" || cleanStr(attr.Key) == "name" {
					property = attr.Val
				}
				if cleanStr(attr.Key) == "content" {
					content = attr.Val
				}
			}
			switch cleanStr(property) {
			case "og:site_name":
				doc.Preview.Name = content
			case "og:title":
				doc.Preview.Title = content
			case "og:description":
				doc.Preview.Description = content
			case "description":
				if len(doc.Preview.Description) == 0 {
					doc.Preview.Description = content
				}
			case "og:url":
				doc.Preview.Link = content
			case "og:image":
				ogImage = true
				ogImgUrl, err := url.Parse(content)
				if err != nil {
					return err
				}
				if !ogImgUrl.IsAbs() {
					ogImgUrl, err = url.Parse(fmt.Sprintf("%s://%s%s", scraper.Url.Scheme, scraper.Url.Host, ogImgUrl.Path))
					if err != nil {
						return err
					}
				}

				doc.Preview.Images = []string{ogImgUrl.String()}

			}

		case "title":
			if tokenType == html.StartTagToken {
				t.Next()
				token = t.Token()
				if len(doc.Preview.Title) == 0 {
					doc.Preview.Title = token.Data
				}
			}

		case "img":
			for _, attr := range token.Attr {
				if cleanStr(attr.Key) == "src" {
					imgUrl, err := url.Parse(attr.Val)
					if err != nil {
						return err
					}
					if !imgUrl.IsAbs() {
						doc.Preview.Images = append(doc.Preview.Images, fmt.Sprintf("%s://%s%s", scraper.Url.Scheme, scraper.Url.Host, imgUrl.Path))
					} else {
						doc.Preview.Images = append(doc.Preview.Images, attr.Val)
					}

				}
			}
		}

		if hasCanonical && headPassed && scraper.MaxRedirect > 0 {
			if !canonicalUrl.IsAbs() {
				absCanonical, err := url.Parse(fmt.Sprintf("%s://%s%s", scraper.Url.Scheme, scraper.Url.Host, canonicalUrl.Path))
				if err != nil {
					return err
				}
				canonicalUrl = absCanonical
			}
			scraper.Url = canonicalUrl
			scraper.EscapedFragmentUrl = nil
			//再次调用,访问解析到的url,直到深度耗尽
			fdoc, err := scraper.getDocument()
			if err != nil {
				return err
			}
			*doc = *fdoc
			return scraper.parseDocument(doc)
		}

		if hasFragment && headPassed && scraper.MaxRedirect > 0 {
			scraper.toFragmentUrl()
			fdoc, err := scraper.getDocument()
			if err != nil {
				return err
			}
			*doc = *fdoc
			return scraper.parseDocument(doc)
		}

		if len(doc.Preview.Title) > 0 && len(doc.Preview.Description) > 0 && ogImage && headPassed {
			return nil
		}

	}

	return nil
}

func avoidByte(b byte) bool {
	i := int(b)
	if i == 127 || (i >= 0 && i <= 31) {
		return true
	}
	return false
}

func escapeByte(b byte) bool {
	i := int(b)
	if i == 32 || i == 35 || i == 37 || i == 38 || i == 43 || (i >= 127 && i <= 255) {
		return true
	}
	return false
}

func metaFragment(token html.Token) bool {
	var name string
	var content string

	for _, attr := range token.Attr {
		if cleanStr(attr.Key) == "name" {
			name = attr.Val
		}
		if cleanStr(attr.Key) == "content" {
			content = attr.Val
		}
	}
	if name == "fragment" && content == "!" {
		return true
	}
	return false
}

func cleanStr(str string) string {
	return strings.ToLower(strings.TrimSpace(str))
}
