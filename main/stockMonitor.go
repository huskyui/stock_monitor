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

func logMsg(stockInfo stock) {
	message := "[股票名称]" + stockInfo.name + "[当前价格]" + fmt.Sprintf("%v", stockInfo.currentPrice)
	fmt.Println(message)
}

func scheduleFetchStockInfoAndNotify() {
	c := cron.New()
	err := c.AddFunc("* * * ? * *", func() {
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

type stock struct {
	name         string
	id           string
	currentPrice float64
}

func setTimeZone() {
	err := os.Setenv("TZ", "Asia/Shanghai")
	if err != nil {
		log.Fatal(err)
	}
}

func fetchStockInfo(stockNum string) stock {
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
		return stock{name: stockInfoArr[1], id: stockNum, currentPrice: currentPrice}
	}
	return stock{}
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

func writeData(stockInfo *stock) {
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
