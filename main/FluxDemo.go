package main

import (
	"context"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"log"
	"math/rand"
	"time"
)

//func main() {
//	setTimeZone()
//	//insertData()
//	influxSimpleQuery()
//
//}

func insertData() {
	client := createInfluxClient()
	org := "huskyui"
	bucket := "mydb"
	writeApi := client.WriteAPIBlocking(org, bucket)

	heroArr := []string{
		"百里守约",
		"花木兰",
		"后羿",
		"小乔",
	}

	for i := 0; i < 1000; i++ {
		time.Sleep(1 * time.Nanosecond)
		var points []*write.Point
		for idx := range heroArr {
			tags := map[string]string{
				"hero": heroArr[idx],
			}

			fields := map[string]interface{}{
				"score": rand.Intn(15),
			}
			point := write.NewPoint("heroscoremeasurement", tags, fields, time.Now())
			points = append(points, point)
		}
		if err := writeApi.WritePoint(context.Background(), points...); err != nil {
			log.Fatal(err)
		}

	}
}
