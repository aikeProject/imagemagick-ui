package core

// 文件处理完毕，发送给前端的数据
type Complete struct {
	Id     string `json:"id"`
	Status int    `json:"status"`
}
