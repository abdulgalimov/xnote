# xnote
## run

```
package main

import (
	"fmt"
	"github.com/abdulgalimov/xnote"
	"github.com/abdulgalimov/xnote/common"
	"github.com/abdulgalimov/xnote/db"
)

func main() {
	fmt.Println("run xnote")
	var dbConfig common.DbConnectConfig
	dbConfig.Host = "localhost"
	dbConfig.Port = 3306
	dbConfig.DriverName = "mysql"
	dbConfig.UserName = "root"
	dbConfig.Password = "123"
	dbConfig.DbName = "xnote_dev"
	xdb, err := db.Connect(dbConfig)
	if err != nil {
		panic(err)
	}

	xnote.Start(xdb)
}
```
