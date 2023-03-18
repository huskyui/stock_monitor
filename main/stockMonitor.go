package main

import (
	"bytes"
	"context"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/robfig/cron"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"gopkg.in/gomail.v2"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	setTimeZone()
	scheduleFetchStockInfoAndNotify()
}

func sendEmail(title, body string) {
	m := gomail.NewMessage()
	m.SetHeader("From", "2207019991@qq.com")
	m.SetHeader("To", "wangpeng91710@gmail.com")
	m.SetHeader("Subject", title)
	m.SetBody("text/html", body)

	d := gomail.NewDialer("smtp.qq.com", 465, "2207019991@qq.com", os.Getenv("QQ_MAIL_PASSWORD"))

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}

func logMsg(stockInfo Stock) {
	sendEmail("股票价格", stockInfo.String())
}

func scheduleFetchStockInfoAndNotify() {
	c := cron.New()
	err := c.AddFunc("0 * * ? * *", func() {
		stockNum := "sh600009"
		stockInfo := fetchStockInfo(stockNum)
		logMsg(stockInfo)
		writeData(&stockInfo)
		influxSimpleQuery()

	})
	if err != nil {
		log.Fatal(err)
	}
	c.Start()
	select {}
}

type Stock struct {
	name         string
	id           string
	currentPrice float64
}

func (stock Stock) String() string {
	return fmt.Sprintf("[股票名称]%s[当前价格]%v", stock.name, stock.currentPrice)
}

func setTimeZone() {
	err := os.Setenv("TZ", "Asia/Shanghai")
	if err != nil {
		log.Fatal(err)
	}
}

func fetchStockInfo(stockNum string) Stock {
	url := fmt.Sprintf("http://qt.gtimg.cn/q=s_%s", stockNum)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(resp.Body)
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		utf8Bytes, err := GbkToUtf8(bodyBytes)
		bodyString := string(utf8Bytes)
		str := strings.Split(bodyString, "=")[1]
		str = strings.TrimPrefix(str, "\"")
		str = strings.TrimSuffix(str, "\";\n")
		stockInfoArr := strings.Split(str, "~")
		currentPrice, err := strconv.ParseFloat(stockInfoArr[3], 64)
		if err != nil {
			log.Fatal(err)
		}
		return Stock{name: stockInfoArr[1], id: stockNum, currentPrice: currentPrice}
	}
	return Stock{}
}

func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := io.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func influxSimpleQuery() {
	client := createInfluxClient()
	org := "huskyui"
	queryApi := client.QueryAPI(org)

	query := `from(bucket: "mydb")
			|> range(start: -20m)
			|> filter(fn: (r) => r._measurement == "stockmeasurement")
            |> mean()`
	results, err := queryApi.Query(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}
	for results.Next() {
		fmt.Println(results.Record())
	}
	if err := results.Err(); err != nil {
		log.Fatal(err)
	}

}

func writeData(stockInfo *Stock) {
	client := createInfluxClient()
	org := "huskyui"
	bucket := "mydb"
	writeApi := client.WriteAPIBlocking(org, bucket)

	tags := map[string]string{
		"id":   stockInfo.id,
		"name": stockInfo.name,
	}

	fields := map[string]interface{}{
		"currentPrice": stockInfo.currentPrice,
	}
	point := write.NewPoint("stockmeasurement", tags, fields, time.Now())
	if err := writeApi.WritePoint(context.Background(), point); err != nil {
		log.Fatal(err)
	}
}

func createInfluxClient() influxdb2.Client {
	return influxdb2.NewClient("http://42.192.90.211:8086",
		"N0ZRZYK5xF8yFB43O7YIYukcty6eFqngli7sUL5ZZhpLN8Ev1FkD22abpmaCBWWkDhShzSQ4oBT7K7T-qNHb7Q==")
}
