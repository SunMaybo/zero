package parser

import (
	"github.com/SunMaybo/zero/common/zlog"
	"github.com/xwb1989/sqlparser"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func ParserCreatedSQL(filePath string) []*sqlparser.DDL {
	var ddls []*sqlparser.DDL
	_ = filepath.Walk(filePath, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".sql") {
			if buff, err := ioutil.ReadFile(path); err != nil {
				zlog.S.Fatal(err)
			} else {
				tree, err := sqlparser.ParseStrictDDL(string(buff))
				if err != nil {
					zlog.S.Errorf("sql parser err,%s,path:%s", err, path)
					return nil
				}
				ddls = append(ddls, tree.(*sqlparser.DDL))
			}
		}
		return nil
	})
	return ddls
}
