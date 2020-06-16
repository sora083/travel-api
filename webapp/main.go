package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

const (
	ENDPOINT = "https://app.rakuten.co.jp/services/api/Travel/SimpleHotelSearch/20170426"
	APP_ID   = ""
)

// レスポンスJSONデータ用構造体
type Result struct {
	PageInfo PageInfo `json:"pagingInfo"`
	Hotels   []*Hotel `json:"hotels"`
}

type PageInfo struct {
	RecordCount int64 `json:"recordCount"`
	PageCount   int64 `json:"pageCount"`
	Page        int64 `json:"page"`
	First       int64 `json:"first"`
	Last        int64 `json:"last"`
}

type Hotel struct {
	HotelInfo []*HotelInfo `json:"hotel"`
}

type HotelInfo struct {
	HotelBasicInfo HotelBasicInfo `json:"hotelBasicInfo"`
}

type HotelBasicInfo struct {
	HotelNo   int64  `json:"hotelNo"`
	HotelName string `json:"hotelName"`
}

func main() {
	e := echo.New()

	// 全てのリクエストで差し込みたいミドルウェア（ログとか）はここ
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// ルーティング
	e.GET("/", func(c echo.Context) error {

		// URL生成
		q := map[string]string{
			"applicationId":   APP_ID, // AppID
			"format":          "json",
			"largeClassCode":  "japan",
			"middleClassCode": "kanagawa",
			"smallClassCode":  "yokohama",
		}
		url := fmt.Sprintf("%s?%s", ENDPOINT, buildQuery(q))
		log.Printf("URL: %s\n", url)

		// URLを叩いてデータを取得
		resp, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		//log.Printf("BODY: %s", body)
		// 取得したデータをJSONデコード
		var result Result
		err = json.Unmarshal(body, &result)
		if err != nil {
			log.Fatal(err)
		}

		return c.JSON(http.StatusOK, echo.Map{
			"pageInfo": result,
		})
	})

	// サーバー起動
	e.Start(":8080")
}

func buildQuery(q map[string]string) string {
	queries := make([]string, 0)
	for k, v := range q {
		qq := fmt.Sprintf("%s=%s", k, v)
		queries = append(queries, qq)
	}
	return strings.Join(queries, "&")
}
