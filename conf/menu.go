package conf

import (
	"os"
	"time"

	log "log/slog"

	"gopkg.in/yaml.v2"
)

type Menu struct {
	Items      []MenuItem `json:"items"`
	LastUpdate time.Time  `json:"lastupdate"`
}
type MenuItem struct {
	Name  string     `json:"name"`
	URL   string     `json:"url"`
	Icon  string     `json:"icon"`
	Class string     `json:"class"`
	Subs  []MenuItem `json:"subs"`
}

// 获取菜单内容
func GetMenu(f string) *Menu {
	data, err := os.ReadFile(f)
	if err != nil {
		log.Error("读取%s 文件失败,Error:", "file", f, "error", err)
	}
	result := &Menu{}
	result.LastUpdate = time.Now()

	err = yaml.Unmarshal(data, result)
	if err != nil {
		log.Error("解析Yaml失败,Error:", "error", err)
	}
	return result
}
