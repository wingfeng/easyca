package engine

import (
	"easyca/conf"
	"easyca/model"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
)

var Db *xorm.Engine

func InitEngine() {
	var err error
	if Db, err = xorm.NewEngine(conf.Default.DB.DBDriver, conf.Default.DB.DBConnection); err != nil {
		panic(fmt.Sprintf("newEngine error is %s", err.Error()))
	}

	//同步表
	if err = Db.Sync2(new(model.ServerCa), new(model.SignCa)); err != nil {
		panic(fmt.Sprintf("Sync2 error is %s", err.Error()))
	}
}

func PageNum(context *gin.Context) (page, rows int) {
	page, _ = strconv.Atoi(context.DefaultQuery("page", "1"))
	rows, _ = strconv.Atoi(context.DefaultQuery("rows", "1"))
	return
}
