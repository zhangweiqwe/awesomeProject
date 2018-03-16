package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func main() {
	fmt.Print("美女")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadBytes('\n')
	x := string(input[0 : len(input)-2])
	const url, page string = "http:/www.btcerise.me/search?keyword=", "&p="
	var Find string
	FileResult, _ := os.OpenFile("re.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 777)
	defer func() {
		time.Sleep(1e9 * 2)
		FileResult.Sync()
		FileResult.Close()
	}()

	for i := 1; i < 101; i++ {
		Find = url + x + page + strconv.Itoa(i)
		h := strings.Repeat("#", i/2) + strings.Repeat(" ", 50-i/2)
		fmt.Printf("\r%02d%%[%s]", i, h)
		time.Sleep(1e6 * 5)
		go Resolve(Find, FileResult)
	}
}

func Resolve(Find string, FileResult io.Writer) {
	Re0, _ := regexp.Compile("<h5.*h5>")
	Re1, _ := regexp.Compile(`^<h5 class="h" name="rsrc"`)
	Re2, _ := regexp.Compile("<span class='highlight'>")
	Re3, _ := regexp.Compile("</span")
	Re4, _ := regexp.Compile("</h5>")
	Re5, _ := regexp.Compile(">")
	Re6, _ := regexp.Compile(`data-hash="`)
	Resp, err := http.Get(Find)
	if err != nil {
		fmt.Println(err)
	}
	Buf, _ := ioutil.ReadAll(Resp.Body)
	buf := Re0.FindAll(Buf, 1000)
	for _, line := range buf {
		line = Re1.ReplaceAll(line, []byte(""))
		line = Re2.ReplaceAll(line, []byte(""))
		line = Re3.ReplaceAll(line, []byte(""))
		line = Re4.ReplaceAll(line, []byte(""))
		line = Re5.ReplaceAll(line, []byte(""))
		line = Re6.ReplaceAll(line, []byte("magnet:?xt=urn:btih:"))
		FileResult.Write(line)
		FileResult.Write([]byte("\n"))
	}
}
