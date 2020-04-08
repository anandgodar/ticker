package controller

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"net/smtp"
	"strings"
	"sync"
)

func GetTickersymbols() []string  {
	return []string{"CTRM"}
}

func StartApplication() error{
	// TODO check length here
	ts := GetTickersymbols()
	if len(ts)<=0{
		fmt.Println("Sorry , empty ticker symbol")
		return errors.New("Empty Ticker symbol, not able to proceed")
	}
	done:= make(chan bool)
	defer close(done)

	symbolsInChannel := prepareTickerSymbols(done,ts)

	workers := make([]<-chan string,10)
	for i:=0;i<10;i++{
		workers[i] = startFetchingData(done,symbolsInChannel,i)
	}

	for range merge(done,workers...){
		fmt.Println("dfd")
	}
	return  nil
}

func merge(done <-chan bool, channels ...<-chan string) <-chan string{
	mergedChannel := make(chan string)
	var wg sync.WaitGroup
	wg.Add(len(channels))
	multiplex := func(c <-chan string) {
			defer wg.Done()
			for val:=range c{
				select {
					case <-done:
						return
					case mergedChannel<-val:
				}//end of select
			}// end of for
	}//end of variable

	for _,c :=range channels{
		go multiplex(c)
	}

	go func() {
		wg.Wait()
		close(mergedChannel)
	}()

	return mergedChannel

}

func startFetchingData(done <-chan bool, symbolsInChannel <-chan string,workerId int ) <-chan string{
  symbols := make(chan string)
  go func(){
  		for syms:= range symbolsInChannel{
			select {
  			 	case <-done:
					return
				case symbols<-syms:
					ScrapeHTML(syms)
					fmt.Println(syms)
			 }

	  }
	  close(symbols)
  }()
  return symbols
}

func prepareTickerSymbols( done <-chan bool, tsy []string) <-chan string {
	tcsc := make(chan string)
	go func() {
		for _,s := range tsy{
			select {
				case <-done:
					return
				case tcsc <-s:

			} // end of select

		}	// end of for loop
		close(tcsc) // close the channel
	}()
	return tcsc
}

func ScrapeHTML(sym string){
	resp, err := http.Get("https://quotes.wsj.com/"+sym)
	if err != nil{
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200{
		log.Fatalf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("ul.cr_newsSummary").Find(".cr_dateStamp").Each(func(index int, item *goquery.Selection) {
		//title := item.Text()
		//linkTag := item.Find("span")
		fmt.Println(item.Text())
		SendMail("127.0.0.1:25","anandgodar@gmail.com","Test","hello",[]string{"anandgodar@gmail.com"})
		//link, _ := linkTag.Attr("href")
		//fmt.Printf("Post #%d: %s - %s\n", index, title, link)
	})

}

//ex: SendMail("127.0.0.1:25", (&mail.Address{"from name", "from@example.com"}).String(), "Email Subject", "message body", []string{(&mail.Address{"to name", "to@example.com"}).String()})
func SendMail(addr, from, subject, body string, to []string) error {
	r := strings.NewReplacer("\r\n", "", "\r", "", "\n", "", "%0a", "", "%0d", "")

	c, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer c.Close()
	if err = c.Mail(r.Replace(from)); err != nil {
		return err
	}
	for i := range to {
		to[i] = r.Replace(to[i])
		if err = c.Rcpt(to[i]); err != nil {
			return err
		}
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	msg := "To: " + strings.Join(to, ",") + "\r\n" +
		"From: " + from + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"Content-Transfer-Encoding: base64\r\n" +
		"\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

	_, err = w.Write([]byte(msg))
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}