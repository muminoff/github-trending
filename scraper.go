package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	_ "strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	//loop
	for {

		dateString := dateString()
		filename := dateString + ".md"

		//create markdown file
		createMarkDown(dateString, filename)

		//TODO: use goroutinez
		scrape("objective-c", filename)
		scrape("go", filename)
		scrape("javascript", filename)

		gitPull()
		gitAddAll()
		gitCommit(dateString)
		gitPush()

		time.Sleep(time.Duration(24) * time.Hour)
	}
}

func dateString() string {
	y, m, d := time.Now().Date()
	mStr := fmt.Sprintf("%d", m)
	dStr := fmt.Sprintf("%d", d)
	if m < 10 {
		mStr = fmt.Sprintf("0%d", m)
	}
	if d < 10 {
		dStr = fmt.Sprintf("0%d", d)
	}
	return fmt.Sprintf("%d-%s-%s", y, mStr, dStr)

}

func createMarkDown(date string, filename string) {

	// open output file
	fo, err := os.Create(filename)
	if err != nil {
		panic(err)
	}

	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	// make a write buffer
	w := bufio.NewWriter(fo)
	w.WriteString("###" + date + "\n")
	w.Flush()
}

func scrape(language string, filename string) {
	var doc *goquery.Document
	var e error
	// var w *bufio.Writer

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(fmt.Sprintf("\n####%s\n", language)); err != nil {
		panic(err)
	}

	if doc, e = goquery.NewDocument(fmt.Sprintf("https://github.com/trending?l=%s", language)); e != nil {
		panic(e.Error())
	}

	doc.Find("li.repo-leaderboard-list-item").Each(func(i int, s *goquery.Selection) {
		title := s.Find("div h2 a").Text()
		owner := s.Find("span.owner-name").Text()
		repoName := s.Find("strong").Text()
		description := s.Find("p.repo-leaderboard-description").Text()
		url, _ := s.Find("h2 a").Attr("href")
		url = "https://github.com" + url
		fmt.Println("owner: ", owner)
		fmt.Println("repo: ", repoName)
		fmt.Println("URL: ", url)
		if _, err = f.WriteString("* [" + title + "](" + url + "): " + description + "\n"); err != nil {
			panic(err)
		}
	})
}
func gitPull() {
	app := "git"
	arg0 := "pull"
	arg1 := "origin"
	arg2 := "master"
	cmd := exec.Command(app, arg0, arg1, arg2)
	out, err := cmd.Output()

	if err != nil {
		println(err.Error())
		return
	}

	print(string(out))
}

func gitAddAll() {
	app := "git"
	arg0 := "add"
	arg1 := "."
	cmd := exec.Command(app, arg0, arg1)
	out, err := cmd.Output()

	if err != nil {
		println(err.Error())
		return
	}

	print(string(out))
}

func gitCommit(date string) {
	app := "git"
	arg0 := "commit"
	arg1 := "-am"
	arg2 := date
	cmd := exec.Command(app, arg0, arg1, arg2)
	out, err := cmd.Output()

	if err != nil {
		println(err.Error())
		return
	}

	print(string(out))
}
func gitPush() {
	app := "git"
	arg0 := "push"
	arg1 := "origin"
	arg2 := "master"
	cmd := exec.Command(app, arg0, arg1, arg2)
	out, err := cmd.Output()

	if err != nil {
		println(err.Error())
		return
	}

	print(string(out))
}
