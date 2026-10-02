// 경기도 시군구별 메모리 반도체(HS 854232) 수출을 조회하는 예제.
//
//	OPENDATA_API_KEY=... go run ./examples/basic
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/kenshin579/opendata-go"
	"github.com/kenshin579/opendata-go/customs"
)

func main() {
	c, err := opendata.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	s := customs.New(c)
	rows, err := s.SigunguItemTrade(context.Background(), customs.SigunguItemTradeParams{
		Start: "202601", End: "202603", SidoCd: "41", HsSgn: "854232",
	})
	if err != nil {
		log.Fatal(err)
	}
	for _, r := range customs.WithoutTotal(rows) {
		v, err := r.ExpUsdAmt.Int64() // 천 달러
		if err != nil {
			continue
		}
		fmt.Printf("%s %-16s %12d 천 달러\n", r.PriodTitle, r.SggNm, v)
	}
}
