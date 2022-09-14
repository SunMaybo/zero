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
		if info == nil {
			return nil
		}
		if (!info.IsDir()) && strings.HasSuffix(info.Name(), ".sql") {
			if buff, err := ioutil.ReadFile(path); err != nil {
				zlog.S.Fatal(err)
			} else {
				ddlStrs := cleanSQLFile(buff)
				for _, ddlStr := range ddlStrs {
					tree, err := sqlparser.ParseStrictDDL(ddlStr)
					if err != nil {
						zlog.S.Errorf("sql parser err,%s,path:%s", err, path)
						return nil
					}
					ddls = append(ddls, tree.(*sqlparser.DDL))
				}

			}
		}
		return nil
	})
	return ddls
}

func cleanSQLFile(buff []byte) []string {
	result := strings.Split(string(buff), "\n")
	data := ""
	for _, s := range result {
		if strings.Index(s, "-- ") >= 0 {
			data += s[:strings.Index(s, "-- ")]
		} else if strings.Index(s, "#") >= 0 {
			data += s[:strings.Index(s, "#")]
		} else {
			data += s
		}
	}
	data = strings.TrimSpace(data)
	sql := strings.Split(data, ";")
	var ddl []string
	for _, s := range sql {
		if strings.HasPrefix(strings.ToLower(strings.TrimSpace(s)), "create") {
			ddl = append(ddl, strings.ToLower(strings.TrimSpace(s)))
		}
	}
	return ddl
}
