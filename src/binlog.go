package main

import (
	"fmt"
	"github.com/siddontang/go-mysql/canal"
	"runtime/debug"
)

type binlogHandler struct {
	canal.DummyEventHandler
	BinlogParser
}

func (h *binlogHandler) OnRow(e *canal.RowsEvent) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Print(r, " ", string(debug.Stack()))
		}
	}()

	// base value for canal.DeleteAction or canal.InsertAction
	var n = 0
	var k = 1

	if e.Action == canal.UpdateAction {
		n = 1
		k = 2
	}

	for i := n; i < len(e.Rows); i += k {

		key := e.Table.Schema + "." + e.Table.Name
		fmt.Printf("监听的表key = %s \n", key)

		switch key {
		case User{}.SchemaName() + "." + User{}.TableName():
			user := User{}
			h.GetBinLogData(&user, e, i)
			switch e.Action {
			case canal.UpdateAction:
				oldUser := User{}
				h.GetBinLogData(&oldUser, e, i-1)
				fmt.Printf("检测到表%s Id=%d is changed from %+v\n to %+v \n", e.Table.Name, user.Id, oldUser, user)

			case canal.InsertAction:
				fmt.Printf("检测到表%s Id=%d is created with %+v \n", e.Table.Name, user.Id, user)

			case canal.DeleteAction:
				fmt.Printf("检测到表%s Id=%d is deleted with %s\n", e.Table.Name, user.Id, user.Title)

			default:
				fmt.Printf("Unknown action")
			}
		}

	}
	return nil
}

func (h *binlogHandler) String() string {
	return "binlogHandler"
}

func binlogListener() {
	c, err := getDefaultCanal()
	if err == nil {
		coords, err := c.GetMasterPos()
		if err == nil {
			c.SetEventHandler(&binlogHandler{})
			c.RunFrom(coords)
		}
	}
}

func getDefaultCanal() (*canal.Canal, error) {
	cfg := canal.NewDefaultConfig()
	//cfg.Addr = fmt.Sprintf("%s:%d", "mariadb", 3307)
	cfg.Addr = fmt.Sprintf("%s:%d", "127.0.0.1", 3308)
	cfg.User = "root"
	cfg.Password = "root"
	cfg.Flavor = "mysql"

	cfg.Dump.ExecutionPath = ""

	return canal.NewCanal(cfg)
}
