package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-redis/redis/v8"
)

const flat_city string = "Krakow"
const countOfPage int = 5
const expectedPrice int = 2500

var ctx = context.Background()

func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func writeFile(date, filename string) {
	file, error := os.Create(filename)
	check(error)
	defer file.Close()
	file.WriteString(date)
}

func hamiltonParser(rdb *redis.Client) {
	for i := 1; i <= countOfPage; i++ {
		url := fmt.Sprintf("https://www.hamiltonmay.com/properties/search?page=%s&branchId=krakow&type=apartment&want=rent", strconv.Itoa(i))
		responce, error := http.Get(url)
		check(error)
		if responce.StatusCode > 400 {
			fmt.Println("Status code:", responce.StatusCode)
		}
		defer responce.Body.Close()

		doc, error := goquery.NewDocumentFromReader(responce.Body)
		check(error)

		file, error := os.Create("post.scv")
		writer := csv.NewWriter(file)

		doc.Find("div.white-container").Find("div.container").Find("div.row").Find("div.col-md-12").Find("div.tabcontent").Find("div.media").Find("div.col-md-5").Each(func(index int, item *goquery.Selection) {
			link, _ := item.Find("div.mb-3").Find("div.col-11").Find("a").Attr("href")
			type apprtment struct {
				link      string
				Prise_PLN int
				Prise_EUR string
				Flat_size int
			}
			var flat_price_pln int
			var flat_price_eur string
			var flat_size int
			item.Find("ul").Find("li").Each(func(index int, item *goquery.Selection) {
				flat_params := strings.TrimSpace(item.Find("div.float-right").Text())
				if strings.Contains(flat_params, "PLN") {
					flat_price := strings.Split(flat_params, " ")
					price, err := strconv.Atoi((strings.Split(flat_price[0], ",")[0]) + (strings.Split(flat_price[0], ",")[1]))
					if err != nil {
						fmt.Printf("You error error is :%s, please check output data %s\n", err, flat_price[1])
					}
					flat_price_pln = price
					//flat_price_eur = strings.Split(flat_price[1], "(")[1]
				}
				if strings.Contains(flat_params, "m2") {
					flat_size, _ = strconv.Atoi((strings.Split(flat_params, "m")[0]))
				}
			})
			if flat_price_pln < expectedPrice && flat_price_pln != 0 {
				//writer.Write(appartmanets)
				apprtments := apprtment{
					link:      link,
					Prise_PLN: flat_price_pln,
					Prise_EUR: flat_price_eur,
					Flat_size: flat_size,
				}
				fmt.Println(apprtments)
				//Redis write
				apartmentJSON, err := json.Marshal(apprtment{
					link:      link,
					Prise_PLN: flat_price_pln,
					Prise_EUR: flat_price_eur,
					Flat_size: flat_size})
				check(err)
				err = rdb.Set(ctx, link, apartmentJSON, 0).Err()
				check(err)
				val, err := rdb.Get(ctx, link).Result()
				check(err)
				fmt.Println("I am from redis")
				fmt.Println(val)
			}

		})
		check(error)
		writer.Flush()
	}

}

func redisClient_test(rdb *redis.Client) {

	err := rdb.Set(ctx, "test", "diego", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "test").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("test", val)

	val2, err := rdb.Get(ctx, "key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}

}

func webHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello %s. welcome to Dante server", r.URL.Path[1:])
}

func main() {
	fmt.Println("Starting HTTPS server")
	http.HandleFunc("/", webHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))

	// rdb := redis.NewClient(&redis.Options{
	//	fmt.Println("Starting Redis server")
	// 	Addr:     "localhost:6379",
	// 	Password: "", // no password set
	// 	DB:       0,  // use default DB
	// })
	// hamiltonParser(rdb)
	//redisClient_test(rdb)
}

// Test printer
// printer := false
// if printer == true {
// 	fmt.Println(link)
// 	fmt.Printf("Flat size is: %d\n", flat_size)
// 	fmt.Printf("Flat price is: %d PLN, %d EUR\n", flat_price_pln, flat_price_eur)
// }

// File writer
// i := 0
// doc.Find("div.white-container").Find("div.container").Find("div.row").Find("div.col-md-12").Find("div.tabcontent").Find("div.media").Find("div.col-md-5").Each(func(index int, item *goquery.Selection) {
// 	doc, err := item.Html()
// 	i += 1
// 	fileName := "index" + strconv.Itoa(i) + ".html"
// 	fmt.Println(strconv.Itoa(i))
// 	check(err)
// 	writeFile(doc, fileName)
// })

// Old parser:
// doc.Find("div.white-container").Find("div.row").Find("div.col-md-12").Find("div.owl-0").Find("div.item").Each(func(index int, item *goquery.Selection) {
// 	city := strings.TrimSpace(item.Find("div.offer-city").Text())
// 	if city == flat_city {
// 		link, _ := item.Find("a").Attr("href")
// 		item_base_details := item.Find("div.offer-description").Find("ul").Find("li")
// 		price := item_base_details.Find("span.dd").Find("strong").Text()
// 		appartmanets := []string{city, link, price}
// 		writer.Write(appartmanets)
// 		fmt.Println(city, link, price)
// 	}

// })
