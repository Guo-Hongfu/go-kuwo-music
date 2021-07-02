package main

import (
	"encoding/json"
	"fmt"
	"kuwo/middleware"
	"os"
	"strconv"
	"time"

	"go.uber.org/zap"

	"kuwo/models"
	"kuwo/request"
)



// 下载全部的音乐
func downAllMusic(url string) *models.MusicResponse {
	token := models.RandStr(11)
	NewReq := request.Requests{
		Method: "GET",
		Url:    url,
		Headers: request.ReqH{
			"Cookie":     fmt.Sprintf("kw_token=%s",token),
			"csrf":       token,
			"Host":       models.KW_HOST,
			"Referer":    models.KW_HOST,
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.131 Safari/537.36",
		},
	}
	_, data := NewReq.News()
	var a = models.MusicResponse{}
	_ = json.Unmarshal(data, &a)
	return &a
}

func search(keyword string) {
	firstUrl := models.FormatMusicListUrl(keyword, 1)
	firstMusic := downAllMusic(firstUrl)
	_ = os.Mkdir(keyword,0666)
	// 开始下载 第一页数据
	zap.S().Info("===========开始下载第一页数据===========")
	middleware.Download(keyword, firstMusic.Data)
	total,_ := strconv.Atoi(firstMusic.Data.Total)
	page := total / 30 + 1
	for pn:=2;pn<=page;pn++{
		zap.S().Info("等待20秒下载下一页歌曲.....................")
		time.Sleep(20*time.Second) // 暂停30秒下载第二页
		zap.S().Infof("===========开始下载第%d页数据===========",pn)

		url := models.FormatMusicListUrl(keyword, uint(pn))
		musicList := downAllMusic(url)
		middleware.Download(keyword, musicList.Data)
	}


}
func main() {
	// 初始化logger
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
	search("周杰伦")
}
