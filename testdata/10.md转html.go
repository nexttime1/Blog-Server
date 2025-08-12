package main

import (
	"fmt"
	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
	"github.com/russross/blackfriday"
	"strings"
)

func main() {
	unsafe := blackfriday.MarkdownCommon([]byte("### 你好\n ```python\nprint('你好')\n```\n - 123 \n \n<script>alert(123)</script>\n\n ![图片](http://xxx.com)"))
	fmt.Println(string(unsafe))

	//html 获取文本信息
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(string(unsafe)))
	//fmt.Println(doc.Text())
	doc.Find("script").Remove()
	fmt.Println("下面：：：：：：：")
	fmt.Println(doc.Text())

	fmt.Println("再下面：：：：：：：")
	converter := md.NewConverter("", true, nil)
	html, _ := doc.Html()
	markdown, err := converter.ConvertString(html)
	fmt.Println(markdown, err)
}
