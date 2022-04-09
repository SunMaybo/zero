package template

import (
	"bytes"
	"log"
	"strconv"
	"strings"
	temp "text/template"
	"unicode"
)

// Copied from golint
var commonInitialisms = []string{"ACL", "API", "ASCII", "CPU", "CSS", "DNS", "EOF", "GUID", "HTML", "HTTP", "HTTPS", "ID", "IP", "JSON", "LHS", "QPS", "RAM", "RHS", "RPC", "SLA", "SMTP", "SQL", "SSH", "TCP", "TLS", "TTL", "UDP", "UI", "UID", "UUID", "URI", "URL", "UTF8", "VM", "XML", "XMPP", "XSRF", "XSS"}
var commonInitialismsReplacer *strings.Replacer
var uncommonInitialismsReplacer *strings.Replacer

func init() {
	var commonInitialismsForReplacer []string
	var uncommonInitialismsForReplacer []string
	for i := len(commonInitialisms) - 1; i >= 0; i-- {
		initialism := commonInitialisms[i]
		commonInitialismsForReplacer = append(commonInitialismsForReplacer, initialism, strings.Title(strings.ToLower(initialism)))
		uncommonInitialismsForReplacer = append(uncommonInitialismsForReplacer, strings.Title(strings.ToLower(initialism)), initialism)
	}
	commonInitialismsReplacer = strings.NewReplacer(commonInitialismsForReplacer...)
	uncommonInitialismsReplacer = strings.NewReplacer(uncommonInitialismsForReplacer...)
}

func Parser[T any](template string, data T) (string, error) {
	buffer := new(bytes.Buffer)
	t := temp.Must(temp.New("anyname").Parse(template))
	if err := t.Execute(buffer, data); err != nil {
		return "", err
	}
	return buffer.String(), nil
}
func MarshalToCamelCase(name string) string {
	if name == "" {
		return ""
	}

	temp := strings.Split(name, "_")
	var s string
	for _, v := range temp {
		vv := []rune(v)
		if len(vv) > 0 {
			if bool(vv[0] >= 'a' && vv[0] <= 'z') { //首字母大写
				vv[0] -= 32
			}
			s += string(vv)
		}
	}

	s = uncommonInitialismsReplacer.Replace(s)
	return s
}
func MarshalToSnakeCase(name string) string {
	name = Lcfirst(name)
	buffer := NewBuffer()
	for i, r := range name {
		if unicode.IsUpper(r) {
			if i != 0 {
				buffer.Append('_')
			}
			buffer.Append(unicode.ToLower(r))
		} else {
			buffer.Append(r)
		}
	}
	return buffer.String()
}
func Lcfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

type Buffer struct {
	*bytes.Buffer
}

func NewBuffer() *Buffer {
	return &Buffer{Buffer: new(bytes.Buffer)}
}

func (b *Buffer) Append(i interface{}) *Buffer {
	switch val := i.(type) {
	case int:
		b.append(strconv.Itoa(val))
	case int64:
		b.append(strconv.FormatInt(val, 10))
	case uint:
		b.append(strconv.FormatUint(uint64(val), 10))
	case uint64:
		b.append(strconv.FormatUint(val, 10))
	case string:
		b.append(val)
	case []byte:
		b.Write(val)
	case rune:
		b.WriteRune(val)
	}
	return b
}

func (b *Buffer) append(s string) *Buffer {
	defer func() {
		if err := recover(); err != nil {
			log.Println("*****内存不够了！******")
		}
	}()
	b.WriteString(s)
	return b
}
