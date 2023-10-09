package gutils

import (
	"context"
	"fmt"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"github.com/yueluoa/infrastructure/ghttp"
	"net/http"
	"strings"
)

// "/data/go/ip2region.xdb"
func GetRegionWithIP(dbPath, ip string) string {
	var region string

	cBuff, err := xdb.LoadContentFromFile(dbPath)
	if err != nil {
		return region
	}

	searcher, err := xdb.NewWithBuffer(cBuff)
	if err != nil {
		return region
	}
	res, err := searcher.SearchByStr(ip)
	if err != nil {
		return region
	}

	region = parseRegion(res)

	return region
}

func parseRegion(region string) string {
	if region == "" {
		return ""
	}
	content := strings.Replace(strings.Replace(region, "|0", "", -1), "0|", "", -1)
	cityData := strings.Replace(content, "|", "/", -1)
	if strings.Contains(cityData, "内网IP") {
		return "内网IP"
	}

	return cityData
}

type IPInfo struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Ipdata struct {
		Info1 string `json:"info1"`
		Info2 string `json:"info2"`
		Info3 string `json:"info3"`
		Isp   string `json:"isp"`
	} `json:"ipdata"`
}

func GetRealIPInfo(ip string) (*IPInfo, error) {
	var info = &IPInfo{}

	url := fmt.Sprintf("https://api.vore.top/api/IPdata?ip=%v", ip)

	client := ghttp.NewClient()
	req, err := client.NewRequest(context.Background(), http.MethodGet, url, ip)
	if err != nil {
		return info, err
	}

	err = client.SendRequest(req, info)

	return info, err
}
