package gen

import (
	"fmt"
	"github.com/SunMaybo/zero/zctl/parser"
	"github.com/SunMaybo/zero/zctl/template"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"strings"
)

func GenerateDao(filepath string, serviceName string) error {
	ddls := parser.ParserCreatedSQL(filepath)
	for _, ddl := range ddls {
		table := String{
			source: ddl.NewName.Name.String(),
		}
		entities := template.JavaEntity{
			ServiceName:  serviceName,
			TableComment: ddl.TableSpec.Options,
			TableUpper:   strings.ToUpper(table.ToCamelWithStartLower()[0:1]) + table.ToCamelWithStartLower()[1:],
			Table:        ddl.NewName.Name.String(),
		}
		isSnowId := false
		isTimeStamp := false
		isTenant := false
		for _, column := range ddl.TableSpec.Columns {
			field := template.Field{
				Name: String{source: column.Name.String()}.ToCamelWithStartLower(),
			}
			if column.Type.Comment != nil {
				field.Comment = string(column.Type.Comment.Val)
			}
			if strings.Contains(field.Name, "tenantId") {
				isTenant = true
			}
			typ := strings.ToLower(column.Type.Type)
			switch typ {
			case "char", "varchar", "text", "longtext", "mediumtext", "tinytext", "json":
				field.Type = "String"
			case "enum", "set":
				field.Type = "String"
			case "blob", "mediumblob", "longblob", "varbinary", "binary":
				field.Type = "[]byte"
			case "timestamp":
				field.Type = "Timestamp"
				if field.Name == "createdAt" {
					isTimeStamp = true
				}
			case "date", "datetime":
				field.Type = "Date"
			case "time":
				field.Type = "Time"
			case "bool":
				field.Type = "Boolean"
			case "tinyint", "smallint", "int", "mediumint":
				field.Type = "Integer"
			case "double":
				field.Type = "Double"
			case "float":
				field.Type = "Float"
			case "bigint":
				field.Type = "Long"
				if field.Name == "id" {
					isSnowId = true
				}
			case "decimal":
				field.Type = "BigDecimal"
			}
			if field.Name == "id" || field.Name == "createdAt" || field.Name == "updatedAt" || field.Name == "deletedFlag" || field.Name == "tenantId" {
				continue
			}
			entities.Fields = append(entities.Fields, field)

			if "" == field.Type {
				return fmt.Errorf("no compatible protobuf type found for `%s`. column: `%s`.`%s`", field.Type, field.Name, field.Comment)
			}
		}
		if isTenant && isTimeStamp && isSnowId {
			entities.Entity = "XbbSnowTimeTenantEntity"
		} else if isTenant && isTimeStamp {
			entities.Entity = "XbbTenantTimeEntity"
		} else if isSnowId && isTimeStamp {
			entities.Entity = "XbbSnowTimeEntity"
		} else if isTenant && isSnowId {
			entities.Entity = "XbbSnowTenantEntity"
		} else if isTenant {
			entities.Entity = "XbbTenantEntity"
		} else if isSnowId {
			entities.Entity = "XbbSnowEntity"
		} else if isTimeStamp {
			entities.Entity = "XbbTimeEntity"
		} else {
			entities.Entity = "XbbEntity"
		}
		result, err := template.Parser(template.JavaEntityTemplate, entities)
		if err != nil {
			zap.S().Fatal(err)
		}
		fmt.Println(result)
		os.MkdirAll(filepath+"/entity", 0777)
		os.MkdirAll(filepath+"/dao", 0777)
		ioutil.WriteFile(filepath+"/entity/"+entities.TableUpper+"Entity.java", []byte(result), 0777)
		resultDao, err := template.Parser(template.JavaDaoTemplate, entities)
		if err != nil {
			zap.S().Fatal(err)
		}
		fmt.Println(resultDao)
		ioutil.WriteFile(filepath+"/dao/"+entities.TableUpper+"Repository.java", []byte(resultDao), 0777)

	}

	return nil
}
