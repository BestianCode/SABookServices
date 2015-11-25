package main

import (
	"github.com/BestianRU/SABModules/SBMConnect"
	"github.com/BestianRU/SABModules/SBMSystem"
)

func main() {
	var (
		jsonConfig SBMSystem.ReadJSONConfig
		rLog       SBMSystem.LogFile
		pg         SBMConnect.PgSQL
	)

	const (
		pName = string("SABook CardDAVMaker Poller")
		pVer  = string("5 2015.11.25.21.00")
	)

	jsonConfig.Init("./CardDAVMaker.log", "./CardDAVMaker.json")

	rLog.ON(jsonConfig)

	if pg.Init(jsonConfig, "insert into aaa_dav_ntu values (0,123);") != 0 {
		rLog.Log("POLLER: Poll insert error!")
	}
	defer pg.Close()

	rLog.OFF()

}
