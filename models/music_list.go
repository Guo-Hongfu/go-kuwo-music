package models



type ListData []MusicList
type MusicResponse struct {
	Code    int   `json:"code"`
	CurTime int   `json:"curTime"`
	Data    RspData `json:"data"`
}

type RspData struct {
	Total string       `json:"total"`
	List  ListData `json:"list"`
}

type MusicList struct {
	Name string `json:"name"`
	Rid  int    `json:"rid"`
}