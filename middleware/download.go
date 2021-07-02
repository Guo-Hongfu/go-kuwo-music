package middleware

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"

	"kuwo/models"
	"kuwo/request"
)

/**
下载器
*/

func Download(outputDir string, resData models.RspData) {
	for _, v := range resData.List {
		getMp3Url := models.FormatMP3Url(v.Rid)
		newReq := request.Requests{
			Method: "GET",
			Url:    getMp3Url,
			Headers: request.ReqH{
				"Host":models.KW_HOST,
				"Referer": fmt.Sprintf("%s%d",models.KW_MUSIC_DETAIL,v.Rid),
				"Cookie": fmt.Sprintf("kw_token=%s",models.RandStr(11)),
			},
		}
		_, data := newReq.News()
		var a = models.KWMP3{}
		_ = json.Unmarshal(data, &a)
		if a.Url != ""{
			downloader := NewFileDownloader(a.Url, v.Name+".mp3", outputDir, 10)
			err := downloader.Run()
			if err != nil {
				zap.S().Errorf("！！！！！！！【%s】下载失败,%s",v.Name,err.Error())
			}
		}else {
			zap.S().Errorf("！！！！！！！歌曲,%s,rid=%d,,获取mp3地址失败", v.Name,v.Rid)
		}

	}
}
