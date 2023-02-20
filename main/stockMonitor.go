package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/DanPlayer/randomname"
	"github.com/eatmoreapple/openwechat"
	"github.com/gin-gonic/gin"
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
	wechatBot := startWechat()
	scheduleFetchStockInfoAndNotify(wechatBot)
	blockWechatBot(wechatBot)
}

func blockWechatBot(weChatBot *openwechat.Bot){
	weChatBot.Block()
}


func startWechat() *openwechat.Bot{
	bot := openwechat.DefaultBot(openwechat.Desktop) // 桌面模式，上面登录不上的可以尝试切换这种模式

	// 注册消息处理函数
	bot.MessageHandler = func(msg *openwechat.Message) {
		if msg.IsText() && msg.Content == "ping" {
			msg.ReplyText("pong")
		}
	}
	// 注册登陆二维码回调
	bot.UUIDCallback = openwechat.PrintlnQrcodeUrl

	// 登陆
	if err := bot.Login(); err != nil {
		fmt.Println(err)
	}
	return bot
}

func sendWeChatMsg(weChatBot *openwechat.Bot,stockInfo stock){
	user, err := weChatBot.GetCurrentUser()
	if err!=nil {
		log.Fatal(err)
	}
	friends, err := user.Friends()
	if err!=nil {
		log.Fatal(err)
	}
	message := "[股票名称]"+stockInfo.name +"[当前价格]"+ fmt.Sprintf("%v",stockInfo.currentPrice)
	friends.SendText(message)

}

func scheduleFetchStockInfoAndNotify(weChatBot *openwechat.Bot){
	c := cron.New()
	err := c.AddFunc("* * * ? * *", func() {
		stockNum := "sh600009"
		stockInfo := fetchStockInfo(stockNum)
		sendWeChatMsg(weChatBot,stockInfo)
	})
	if err != nil {
		log.Fatal(err)
	}
	c.Start()
}

func cronFunc() {
	c := cron.New()
	err := c.AddFunc("* * * ? * *", func() {
		fmt.Println(time.Now(), "hello cron")
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
	os.Setenv("TZ", "Asia/Shanghai")
}

func fetchStockInfo(stockNum string) stock {
	url := fmt.Sprintf("http://qt.gtimg.cn/q=s_%s", stockNum)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
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
			|> range(start: -10m)
			|> filter(fn: (r) => r._measurement == "stockmeasurement")`
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

func simpleGinWebServer() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message":    "pong",
			"additional": "fuck u!",
		})
	})
	r.GET("/favicon.ico", func(c *gin.Context) {
		c.JSON(http.StatusGatewayTimeout, gin.H{
			"message": "this is favicon.ico",
		})
	})
	r.GET("/random", func(c *gin.Context) {
		name := randomname.GenerateName()
		c.JSON(http.StatusOK, gin.H{
			"nickname": name,
		})
	})
	r.Run()
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
