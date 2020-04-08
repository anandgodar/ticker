package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"projects/internal/controller"
	"strings"
	"sync"
	"time"
	"github.com/gomodule/redigo/redis"
)
func runapp(w http.ResponseWriter, r *http.Request) {
	os.Exit(1000)
	message := r.URL.Path
	message = strings.TrimPrefix(message, "/")
	startApp()
	message = "Hello " + message
	w.Write([]byte(message))
}
func get(c redis.Conn) error {

	// Simple GET example with String helper
	key := "Favorite Movie"
	s, err := redis.String(c.Do("GET", key))
	if err != nil {
		return (err)
	}
	fmt.Printf("%s = %s\n", key, s)

	// Simple GET example with Int helper
	key = "Release Year"
	i, err := redis.Int(c.Do("GET", key))
	if err != nil {
		return (err)
	}
	fmt.Printf("%s = %d\n", key, i)

	// Example where GET returns no results
	key = "Nonexistent Key"
	s, err = redis.String(c.Do("GET", key))
	if err == redis.ErrNil {
		fmt.Printf("%s does not exist\n", key)
	} else if err != nil {
		return err
	} else {
		fmt.Printf("%s = %s\n", key, s)
	}

	return nil
}
func set(c redis.Conn) error {
	_, err := c.Do("SET", "Favorite Movie", "Repo Man")
	if err != nil {
		return err
	}
	_, err = c.Do("SET", "Release Year", 1984)
	if err != nil {
		return err
	}
	return nil
}

func setRedis(w http.ResponseWriter, r *http.Request){
	conn, err := redis.Dial("tcp", "redis-server:6379")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	set(conn)
	get(conn)
	w.Write([]byte("Finished Setting ang Getting"))

}

func main() {
	//http.HandleFunc("/stock", runapp)
	//http.HandleFunc("/setRedis", setRedis)
	//if err := http.ListenAndServe(":8080", nil); err != nil {
	//	panic(err)
	//}
	startApp()
}
func startApp(){
	/*
		Project Compontent

	// generate common channel
	// Fannout to go routines
	// Fanin to common channel

	inputs
	 */


	 start := time.Now()

	 tickerSymbols := controller.GetTickersymbols()
	 if len(tickerSymbols)>0{
	 	controller.StartApplication()
	 }else{
	 	fmt.Println("Not able to proceed , No ticker symbols")
	 }

	 fmt.Printf("Time to finish %f",time.Since(start).Seconds())
	//var wg sync.WaitGroup
	// in := make(chan string)
	// out := make(chan string)
	// for i:=0;i<10;i++{
	//	 wg.Add(1)
	// 	go worker(i,in,out,&wg)
	// }
	//
	// for _,v:=range tickerSymbols{
	// 	in <-v
	// }
	//
	// close(in)
	//
	// for  {
	//	 select {
	// 			case msg:=<-out:
	// 				fmt.Println(msg)
	//	 default:
	//		 fmt.Println("Stopped")
	//	}
	//
	// }
	//close(out)
	//for i:=0;i<len(tickerSymbols); i++{
	//	fmt.Println(<-out)
	//}
	// close(out)
	// chData := make(chan string,len(tickerSymbols))
	//
	// for _,ts := range tickerSymbols{
	//
	// 	go fetchData(&wg,ts,chData)
	// }
	// wg.Wait()
	// close(chData)
	//for val := range chData{
	//	fmt.Println(val)
	//}

	fmt.Println("main stopped")
}

func printMe(out chan string){

}
func worker(i int,in chan string , out chan string,wg *sync.WaitGroup){
	for  r:=range in {
		wg.Done()
		fmt.Printf("job %d value %s",i,r)
		out <- r

	}

}
func printData(chData chan string,chDone chan bool){
		for val := range chData{
			fmt.Println(val)
		}
		chDone <-true
		defer close(chDone)
}

func fetchData(wg *sync.WaitGroup,ts string, chData chan string){
 	defer wg.Done()
	chData <- ts
	//close(chData)
}
