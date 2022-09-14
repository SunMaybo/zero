package gen

import (
	"bytes"
	"database/sql"
	"encoding/binary"
	"fmt"
	"github.com/SunMaybo/zero/common/zlog"
	"github.com/SunMaybo/zero/zctl/parser"
	"io/ioutil"
	"regexp"
	"sort"
	"strings"

	"github.com/chuckpreslar/inflect"
	"github.com/serenize/snaker"
)

const (
	// proto3 is a describing the proto3 syntax type.
	proto3 = "proto3"

	// indent represents the indentation amount for fields. the style guide suggests
	// two spaces
	indent = "  "
)

// GenerateSchema generates a protobuf schema from a database connection and a package name.
// A list of tables to ignore may also be supplied.
// The returned schema implements the `fmt.Stringer` interface, in order to generate a string
// representation of a protobuf schema.
// Do not rely on the structure of the Generated schema to provide any context about
// the protobuf types. The schema reflects the layout of a protobuf file and should be used
// to pipe the output of the `Schema.String()` to a file.
func GenerateSchema(filepath string, serviceName string) error {
	s := &Schema{}
	s.Syntax = proto3
	s.ServiceName = serviceName
	s.Package = serviceName
	s.GoPackage = serviceName
	ddls := parser.ParserCreatedSQL(filepath)
	var cols []Column
	for _, ddl := range ddls {
		for _, column := range ddl.TableSpec.Columns {
			var data int64
			if column.Type.Length != nil {
				bytebuff := bytes.NewBuffer(column.Type.Length.Val)
				binary.Read(bytebuff, binary.BigEndian, &data)
			}
			isNo := "NO"
			if !column.Type.NotNull {
				isNo = "Yes"
			}
			cols = append(cols, Column{
				TableComment:  ddl.TableSpec.Options,
				TableName:     ddl.NewName.Name.String(),
				ColumnName:    column.Name.String(),
				IsNullable:    isNo,
				ColumnType:    column.Type.Type,
				DataType:      column.Type.Type,
				ColumnComment: string(column.Type.Comment.Val),
				CharacterMaximumLength: sql.NullInt64{
					Int64: data,
					Valid: false,
				},
				NumericPrecision: sql.NullInt64{
					Int64: data,
					Valid: false,
				},
				NumericScale: sql.NullInt64{
					Int64: 0,
					Valid: false,
				},
			})
		}

	}
	err := typesFromColumns(s, cols, nil)
	if nil != err {
		return err
	}
	sort.Sort(s.Imports)
	sort.Sort(s.Messages)
	sort.Sort(s.Enums)
	if err := ioutil.WriteFile(filepath+"/"+serviceName+".proto", []byte(s.String()), 0777); err != nil {
		zlog.S.Fatal(err)
	}
	fmt.Println(s)
	return nil
}

// typesFromColumns creates the appropriate schema properties from a collection of column types.
func typesFromColumns(s *Schema, cols []Column, ignoreTables []string) error {
	messageMap := map[string]*Message{}
	ignoreMap := map[string]bool{}
	for _, ig := range ignoreTables {
		ignoreMap[ig] = true
	}

	for _, c := range cols {
		if _, ok := ignoreMap[c.TableName]; ok {
			continue
		}

		messageName := snaker.SnakeToCamel(c.TableName)
		messageName = inflect.Singularize(messageName)

		msg, ok := messageMap[messageName]
		if !ok {
			messageMap[messageName] = &Message{Name: messageName, Comment: c.TableComment}
			msg = messageMap[messageName]
		}

		err := parseColumn(s, msg, c)
		if nil != err {
			return err
		}
	}

	for _, v := range messageMap {
		s.Messages = append(s.Messages, v)
	}

	return nil
}

// Schema is a representation of a protobuf schema.
type Schema struct {
	Syntax      string
	ServiceName string
	GoPackage   string
	Package     string
	Imports     sort.StringSlice
	Messages    MessageCollection
	Enums       EnumCollection
}

// MessageCollection represents a sortable collection of messages.
type MessageCollection []*Message

func (mc MessageCollection) Len() int {
	return len(mc)
}

func (mc MessageCollection) Less(i, j int) bool {
	return mc[i].Name < mc[j].Name
}

func (mc MessageCollection) Swap(i, j int) {
	mc[i], mc[j] = mc[j], mc[i]
}

// EnumCollection represents a sortable collection of enums.
type EnumCollection []*Enum

func (ec EnumCollection) Len() int {
	return len(ec)
}

func (ec EnumCollection) Less(i, j int) bool {
	return ec[i].Name < ec[j].Name
}

func (ec EnumCollection) Swap(i, j int) {
	ec[i], ec[j] = ec[j], ec[i]
}

// AppendImport adds an import to a schema if the specific import does not already exist in the schema.
func (s *Schema) AppendImport(imports string) {
	shouldAdd := true
	for _, si := range s.Imports {
		if si == imports {
			shouldAdd = false
			break
		}
	}

	if shouldAdd {
		s.Imports = append(s.Imports, imports)
	}

}

// String returns a string representation of a Schema.
func (s *Schema) String() string {
	buf := new(bytes.Buffer)
	buf.WriteString(fmt.Sprintf("syntax = \"%s\";\n", s.Syntax))
	buf.WriteString("\n")
	buf.WriteString(fmt.Sprintf("option go_package =\"%s\";\n", s.GoPackage))
	buf.WriteString("\n")
	buf.WriteString(fmt.Sprintf("package %s;\n", s.Package))
	buf.WriteString(fmt.Sprintf("option java_package = \"cn.xunhou.grpc.proto.%s\";", UpperHead(s.ServiceName)))
	buf.WriteString("\n")
	buf.WriteString(fmt.Sprintf("option java_outer_classname = \"%sServiceProto\";", UpperHead(s.ServiceName)))
	buf.WriteString("\n")
	buf.WriteString("import \"google/protobuf/timestamp.proto\";")
	buf.WriteString("\n")
	buf.WriteString("import \"google/protobuf/empty.proto\";")
	buf.WriteString("\n")
	buf.WriteString("import \"validate/validate.proto\";")
	buf.WriteString("\n")
	buf.WriteString("\n")

	for _, m := range s.Messages {
		m.GenDefaultMessage(buf)
		//m.GenRpcAddReqRespMessage(buf)
		//m.GenRpcUpdateReqMessage(buf)
		//	m.GenRpcDelReqMessage(buf)
		//	m.GenRpcGetByIdReqMessage(buf)
		//	m.GenRpcSearchReqMessage(buf)
	}

	buf.WriteString("\n")

	if len(s.Enums) > 0 {
		buf.WriteString("// ------------------------------------ \n\n")
		for _, e := range s.Enums {
			buf.WriteString(fmt.Sprintf("%s\n", e))
		}
	}

	buf.WriteString("\n")
	//funcTpl := "service " + UpperHead(s.ServiceName) + "Server { \n\n"
	//for _, m := range s.Messages {
	//	funcTpl += "\t //-----------------------" + m.Comment + "----------------------- \n"
	//	funcTpl += "\t rpc Add" + m.Name + "(Add" + m.Name + "BeRequest) returns (google.protobuf.Empty); \n"
	//	funcTpl += "\t rpc Update" + m.Name + "(Update" + m.Name + "BeRequest) returns (google.protobuf.Empty); \n"
	//	funcTpl += "\t rpc Del" + m.Name + "(Del" + m.Name + "BeRequest) returns (Del" + m.Name + "BeResponse); \n"
	//	funcTpl += "\t rpc Get" + m.Name + "ById(Get" + m.Name + "ByIdBeRequest) returns (Get" + m.Name + "ByIdBeResponse); \n"
	//	funcTpl += "\t rpc Search" + m.Name + "(Query" + m.Name + "BeRequest) returns (Query" + m.Name + "BeResponse); \n"
	//}
	//funcTpl = funcTpl + "\n}"
	//buf.WriteString(funcTpl)

	return buf.String()
}
func UpperHead(s string) string {
	return strings.ToUpper(s[0:1]) + s[1:]
}

// Enum represents a protocol buffer enumerated type.
type Enum struct {
	Name    string
	Comment string
	Fields  []EnumField
}

// String returns a string representation of an Enum.
func (e *Enum) String() string {
	buf := new(bytes.Buffer)

	buf.WriteString(fmt.Sprintf("// %s \n", e.Comment))
	buf.WriteString(fmt.Sprintf("enum %s {\n", e.Name))

	for _, f := range e.Fields {
		buf.WriteString(fmt.Sprintf("%s%s;\n", indent, f))
	}

	buf.WriteString("}\n")

	return buf.String()
}

// AppendField appends an EnumField to an Enum.
func (e *Enum) AppendField(ef EnumField) error {
	for _, f := range e.Fields {
		if f.Tag() == ef.Tag() {
			return fmt.Errorf("tag `%d` is already in use by field `%s`", ef.Tag(), f.Name())
		}
	}

	e.Fields = append(e.Fields, ef)

	return nil
}

// EnumField represents a field in an enumerated type.
type EnumField struct {
	name string
	tag  int
}

// NewEnumField constructs an EnumField type.
func NewEnumField(name string, tag int) EnumField {
	name = strings.ToUpper(name)

	re := regexp.MustCompile(`([^\w]+)`)
	name = re.ReplaceAllString(name, "_")

	return EnumField{name, tag}
}

// String returns a string representation of an Enum.
func (ef EnumField) String() string {
	return fmt.Sprintf("%s = %d", ef.name, ef.tag)
}

// Name returns the name of the enum field.
func (ef EnumField) Name() string {
	return ef.name
}

// Tag returns the identifier tag of the enum field.
func (ef EnumField) Tag() int {
	return ef.tag
}

// newEnumFromStrings creates an enum from a name and a slice of strings that represent the names of each field.
func newEnumFromStrings(name, comment string, ss []string) (*Enum, error) {
	enum := &Enum{}
	enum.Name = name
	enum.Comment = comment

	for i, s := range ss {
		err := enum.AppendField(NewEnumField(s, i))
		if nil != err {
			return nil, err
		}
	}

	return enum, nil
}

// Service represents a protocol buffer service.
// TODO: Implement this in a schema.
type Service struct{}

// Message represents a protocol buffer message.
type Message struct {
	Name    string
	Comment string
	Fields  []MessageField
}

//gen default message
func (m Message) GenDefaultMessage(buf *bytes.Buffer) {
	mOrginName := m.Name
	mOrginFields := m.Fields

	curFields := []MessageField{}
	var filedTag int
	for _, field := range m.Fields {
		if isInSlice([]string{"version", "del_state", "delete_time"}, field.Name) {
			continue
		}
		filedTag++
		field.tag = filedTag
		field.Name = From(field.Name).ToCamelWithStartLower()
		if field.Comment == "" {
			field.Comment = field.Name
		}
		curFields = append(curFields, field)
	}
	m.Fields = curFields
	buf.WriteString(fmt.Sprintf("%s\n", m))

	//reset
	m.Name = mOrginName
	m.Fields = mOrginFields
}

//gen add req message
func (m Message) GenRpcAddReqRespMessage(buf *bytes.Buffer) {
	mOrginName := m.Name
	mOrginFields := m.Fields

	//req
	m.Name = "Add" + mOrginName + "BeRequest"
	curFields := []MessageField{}
	var filedTag int
	for _, field := range m.Fields {
		if isInSlice([]string{"id", "create_time", "update_time", "version", "del_state", "delete_time"}, field.Name) {
			continue
		}
		filedTag++
		field.tag = filedTag
		field.Name = From(field.Name).ToCamelWithStartLower()
		if field.Comment == "" {
			field.Comment = field.Name
		}
		curFields = append(curFields, field)
	}
	m.Fields = curFields
	buf.WriteString(fmt.Sprintf("%s\n", m))

	//reset
	m.Name = mOrginName
	m.Fields = mOrginFields

	//resp
	m.Name = "Add" + mOrginName + "BeResponse"
	m.Fields = []MessageField{}
	buf.WriteString(fmt.Sprintf("%s\n", m))

	//reset
	m.Name = mOrginName
	m.Fields = mOrginFields

}

//gen add resp message
func (m Message) GenRpcUpdateReqMessage(buf *bytes.Buffer) {
	mOrginName := m.Name
	mOrginFields := m.Fields

	m.Name = "Update" + mOrginName + "BeRequest"
	curFields := []MessageField{}
	var filedTag int
	for _, field := range m.Fields {
		if isInSlice([]string{"create_time", "update_time", "version", "del_state", "delete_time"}, field.Name) {
			continue
		}
		filedTag++
		field.tag = filedTag
		field.Name = From(field.Name).ToCamelWithStartLower()
		if field.Comment == "" {
			field.Comment = field.Name
		}
		curFields = append(curFields, field)
	}
	m.Fields = curFields
	buf.WriteString(fmt.Sprintf("%s\n", m))

	//reset
	m.Name = mOrginName
	m.Fields = mOrginFields

	//resp
	m.Name = "Update" + mOrginName + "BeResponse"
	m.Fields = []MessageField{}
	buf.WriteString(fmt.Sprintf("%s\n", m))

	//reset
	m.Name = mOrginName
	m.Fields = mOrginFields
}

//gen add resp message
func (m Message) GenRpcDelReqMessage(buf *bytes.Buffer) {
	mOrginName := m.Name
	mOrginFields := m.Fields

	m.Name = "Del" + mOrginName + "BeRequest"
	m.Fields = []MessageField{
		{Name: "id", Typ: "int64", tag: 1, Comment: "id"},
	}
	buf.WriteString(fmt.Sprintf("%s\n", m))

	//reset
	m.Name = mOrginName
	m.Fields = mOrginFields

	//resp
	m.Name = "Del" + mOrginName + "BeResponse"
	m.Fields = []MessageField{}
	buf.WriteString(fmt.Sprintf("%s\n", m))

	//reset
	m.Name = mOrginName
	m.Fields = mOrginFields
}

//gen add resp message
func (m Message) GenRpcGetByIdReqMessage(buf *bytes.Buffer) {
	mOrginName := m.Name
	mOrginFields := m.Fields

	m.Name = "Get" + mOrginName + "ByIdBeRequest"
	m.Fields = []MessageField{
		{Name: "id", Typ: "int64", tag: 1, Comment: "id"},
	}
	buf.WriteString(fmt.Sprintf("%s\n", m))

	//reset
	m.Name = mOrginName
	m.Fields = mOrginFields

	//resp
	firstWord := strings.ToLower(string(m.Name[0]))
	m.Name = "Get" + mOrginName + "ByIdBeResponse"
	m.Fields = []MessageField{
		{Typ: mOrginName, Name: From(firstWord + mOrginName[1:]).ToCamelWithStartLower(), tag: 1, Comment: From(firstWord + mOrginName[1:]).ToCamelWithStartLower()},
	}
	buf.WriteString(fmt.Sprintf("%s\n", m))

	//reset
	m.Name = mOrginName
	m.Fields = mOrginFields
}

//gen add resp message
func (m Message) GenRpcSearchReqMessage(buf *bytes.Buffer) {
	mOrginName := m.Name
	mOrginFields := m.Fields

	m.Name = "Query" + mOrginName + "BeRequest"
	curFields := []MessageField{
		{Typ: "int64", Name: "page", tag: 1, Comment: "page"},
		{Typ: "int64", Name: "pageSize", tag: 2, Comment: "pageSize"},
	}
	var filedTag = len(curFields)
	for _, field := range m.Fields {
		if isInSlice([]string{"version", "del_state", "delete_time"}, field.Name) {
			continue
		}
		filedTag++
		field.tag = filedTag
		field.Name = From(field.Name).ToCamelWithStartLower()
		if field.Comment == "" {
			field.Comment = field.Name
		}
		curFields = append(curFields, field)
	}
	m.Fields = curFields
	buf.WriteString(fmt.Sprintf("%s\n", m))

	//reset
	m.Name = mOrginName
	m.Fields = mOrginFields

	//resp
	firstWord := strings.ToLower(string(m.Name[0]))
	m.Name = "Query" + mOrginName + "BeResponse"
	m.Fields = []MessageField{
		{Typ: "repeated " + mOrginName, Name: From(firstWord + mOrginName[1:]).ToCamelWithStartLower(), tag: 1, Comment: From(firstWord + mOrginName[1:]).ToCamelWithStartLower()},
	}
	buf.WriteString(fmt.Sprintf("%s\n", m))

	//reset
	m.Name = mOrginName
	m.Fields = mOrginFields
}

// String returns a string representation of a Message.
func (m Message) String() string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("message %s {\n", m.Name))
	for _, f := range m.Fields {
		buf.WriteString(fmt.Sprintf("//%s\n", f.Comment))
		buf.WriteString(fmt.Sprintf("%s%s;\n", indent, f))
	}
	buf.WriteString("}\n")

	return buf.String()
}

// AppendField appends a message field to a message. If the tag of the message field is in use, an error will be returned.
func (m *Message) AppendField(mf MessageField) error {
	for _, f := range m.Fields {
		if f.Tag() == mf.Tag() {
			return fmt.Errorf("tag `%d` is already in use by field `%s`", mf.Tag(), f.Name)
		}
	}

	m.Fields = append(m.Fields, mf)

	return nil
}

// MessageField represents the field of a message.
type MessageField struct {
	Typ     string
	Name    string
	tag     int
	Comment string
}

// NewMessageField creates a new message field.
func NewMessageField(typ, name string, tag int, comment string) MessageField {
	return MessageField{typ, name, tag, comment}
}

// Tag returns the unique numbered tag of the message field.
func (f MessageField) Tag() int {
	return f.tag
}

// String returns a string representation of a message field.
func (f MessageField) String() string {
	return fmt.Sprintf("%s %s = %d", f.Typ, f.Name, f.tag)
}

// Column represents a database column.
type Column struct {
	TableName              string
	TableComment           string
	ColumnName             string
	IsNullable             string
	DataType               string
	CharacterMaximumLength sql.NullInt64
	NumericPrecision       sql.NullInt64
	NumericScale           sql.NullInt64
	ColumnType             string
	ColumnComment          string
}

// Table represents a database table.
type Table struct {
	TableName  string
	ColumnName string
}

// parseColumn parses a column and inserts the relevant fields in the Message. If an enumerated type is encountered, an Enum will
// be added to the Schema. Returns an error if an incompatible protobuf data type cannot be found for the database column type.
func parseColumn(s *Schema, msg *Message, col Column) error {
	typ := strings.ToLower(col.DataType)
	var fieldType string

	switch typ {
	case "char", "varchar", "text", "longtext", "mediumtext", "tinytext", "json":
		fieldType = "string"
	case "enum", "set":
		// Parse c.ColumnType to get the enum list
		enumList := regexp.MustCompile(`[enum|set]\((.+?)\)`).FindStringSubmatch(col.ColumnType)
		enums := strings.FieldsFunc(enumList[1], func(c rune) bool {
			cs := string(c)
			return "," == cs || "'" == cs
		})

		enumName := inflect.Singularize(snaker.SnakeToCamel(col.TableName)) + snaker.SnakeToCamel(col.ColumnName)
		enum, err := newEnumFromStrings(enumName, col.ColumnComment, enums)
		if nil != err {
			return err
		}

		s.Enums = append(s.Enums, enum)

		fieldType = enumName
	case "blob", "mediumblob", "longblob", "varbinary", "binary":
		fieldType = "bytes"
	case "date", "time", "datetime", "timestamp":
		//s.AppendImport("google/protobuf/timestamp.proto")
		fieldType = "int64"
	case "bool":
		fieldType = "bool"
	case "tinyint", "smallint", "int", "mediumint", "bigint":
		fieldType = "int64"
	case "float", "decimal", "double":
		fieldType = "double"
	}

	if "" == fieldType {
		return fmt.Errorf("no compatible protobuf type found for `%s`. column: `%s`.`%s`", col.DataType, col.TableName, col.ColumnName)
	}

	field := NewMessageField(fieldType, col.ColumnName, len(msg.Fields)+1, col.ColumnComment)

	err := msg.AppendField(field)
	if nil != err {
		return err
	}

	return nil
}

func isInSlice(slice []string, s string) bool {
	for i, _ := range slice {
		if slice[i] == s {
			return true
		}
	}
	return false
}
