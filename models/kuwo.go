package models

import (
	"fmt"
	"math/rand"
	"net/url"
	"time"
)

const KW_HOST string = "http://www.kuwo.cn"
const KW_MUSIC_DETAIL string = "http://www.kuwo.cn/play_detail/"

// 格式化获取音乐数据Url http://www.kuwo.cn/url?format=mp3&rid=324244&response=url&type=convert_url3&br=320kmp3&from=web
func FormatMusicListUrl(keyword string, pn uint) string {
	return fmt.Sprintf("%s/api/www/search/searchMusicBykeyWord?key=%s&pn=%d&rn=30", KW_HOST, url.QueryEscape(keyword), pn)
}

func FormatMP3Url(rid int) string {
	return fmt.Sprintf("%s/url?format=mp3&rid=%d&response=url&type=convert_url3&br=320kmp3&from=web", KW_HOST, rid)
}

// 生成一串随机字符串
func RandStr(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)

}
