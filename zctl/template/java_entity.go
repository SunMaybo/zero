package template

const JavaEntityTemplate = `
package cn.xunhou.grpc.{{.ServiceName}}.entity;

import cn.xunhou.cloud.core.select.IFieldsInfo;
import cn.xunhou.cloud.dao.annotation.XbbTable;
import cn.xunhou.cloud.dao.xhjdbc.*;
import java.sql.*;
import lombok.Getter;
import lombok.Setter;
import lombok.ToString;
import lombok.experimental.Accessors;

import java.util.Arrays;

/**
 * {{.TableComment}}
 *
 * @author system
 * @since 2022-08-03 16:03:26
 */
@Getter
@Setter
@ToString
@Accessors(chain = true)
@XbbTable(table = "{{.Table}}")
public class {{.TableUpper}}Entity extends {{.Entity}} {

 {{range $index, $field := .Fields}}
    /**
     * {{$field.Comment}}
     */
    private {{$field.Type}} {{$field.Name}};
    {{end}}
}
`

type JavaEntity struct {
	ServiceName  string
	TableComment string
	TableUpper   string
	Table        string
	Entity       string
	Fields       []Field
}
type Field struct {
	Name    string
	Type    string
	Comment string
}
