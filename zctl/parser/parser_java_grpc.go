package parser

import (
	"github.com/SunMaybo/zero/zctl/file"
	"strings"
)

const (
	IsNotStream MethodType = iota
	DuplexStream
	ReturnSimplexStream
)

type MethodType int

type JavaRpcImpl struct {
	PackageName     string
	ServiceName     string
	GrpcFileName    string
	ServiceBaseName string
	MethodSigns     []Method
}

type Method struct {
	Method        string
	Param1        string
	Param2        string
	Param2T       string
	ReturnParam   string
	MethodComment string
	IsStream      MethodType
}
type MethodBlock struct {
	CodeBlock string
	Comment   string
}

func ParserJavaGrpc(filePath string) JavaRpcImpl {

	buff, err := file.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	if !strings.Contains(string(buff), "implements io.grpc.BindableService") {
		return JavaRpcImpl{}
	}
	codeLines := strings.Split(string(buff), "\n")
	javaRpcImpl := JavaRpcImpl{}
	isSkip := true
	var methodBlocks []MethodBlock
	starIdx := 0
	endIdx := 0
	currentComment := ""
	for idx, line := range codeLines {
		if isSkip && strings.Contains(line, "implements io.grpc.BindableService") {
			//pre class  implements
			preIdx := strings.Index(line, "class")
			sufidx := strings.Index(line, "implements")
			if preIdx > 0 && sufidx > 0 && sufidx > preIdx {
				javaRpcImpl.ServiceBaseName = strings.TrimSpace(line[preIdx+6 : sufidx])
			} else {
				panic("parser java grpc error")
			}
			isSkip = false
			continue
		} else if isSkip {
			continue
		}
		if strings.Contains(line, "io.grpc.ServerServiceDefinition bindService") {
			break
		}
		if strings.Contains(line, "/**") {
			codeLines[idx] = strings.TrimSpace(line)
			for i := idx + 1; i < len(codeLines); i++ {
				if strings.Contains(codeLines[i], "*/") {
					currentComment = strings.Join(codeLines[idx:i+1], "\n")
					break
				}

			}
		}
		if strings.Contains(line, "public") {
			starIdx = idx
		}
		if strings.Contains(line, "}") {
			endIdx = idx
			block := MethodBlock{
				Comment:   currentComment,
				CodeBlock: strings.Join(codeLines[starIdx:endIdx+1], ""),
			}
			methodBlocks = append(methodBlocks, block)
			currentComment = ""
			starIdx = 0
			endIdx = 0
		}

	}
	for _, block := range methodBlocks {
		mtd := Method{
			MethodComment: block.Comment,
		}
		if strings.Contains(block.CodeBlock, "public void") {
			//获取方法名称
			mtdStartIdx := strings.Index(block.CodeBlock, "public void")
			mtdEndIdx := strings.Index(block.CodeBlock, "(")
			mtd.Method = strings.TrimSpace(block.CodeBlock[mtdStartIdx+len("public void") : mtdEndIdx])
			//获取参数
			blockNext := block.CodeBlock[mtdEndIdx+1:]
			paramEndIdx := strings.Index(blockNext, ",")
			mtd.Param1 = strings.TrimSpace(blockNext[:paramEndIdx])
			blockNext = blockNext[paramEndIdx+1:]
			paramEndIdx = strings.Index(blockNext, ")")
			mtd.IsStream = IsNotStream
			mtd.Param2 = strings.TrimSpace(blockNext[:paramEndIdx])
			mtd.Param2T = strings.TrimSpace(mtd.Param2[strings.Index(mtd.Param2, "<")+1 : strings.Index(mtd.Param2, ">")])
		} else {
			mtdStartIdx := strings.Index(block.CodeBlock, "public")
			mtdEndIdx := strings.Index(block.CodeBlock, ">")
			mtd.ReturnParam = strings.TrimSpace(block.CodeBlock[mtdStartIdx+len("public") : mtdEndIdx+1])
			mtdMthodEndIdx := strings.Index(block.CodeBlock, "(")
			mtd.Method = strings.TrimSpace(block.CodeBlock[mtdEndIdx+1 : mtdMthodEndIdx])
			mtdEndIdx = strings.Index(block.CodeBlock, ")")
			mtd.Param1 = strings.TrimSpace(block.CodeBlock[mtdMthodEndIdx+1 : mtdEndIdx])
			mtd.IsStream = DuplexStream
		}

		javaRpcImpl.MethodSigns = append(javaRpcImpl.MethodSigns, mtd)
		javaRpcImpl.ServiceName = "Abstract" + javaRpcImpl.ServiceBaseName
		grpcJavaName := strings.Split(filePath, "/")[len(strings.Split(filePath, "/"))-1]
		grpcJavaName = strings.ReplaceAll(grpcJavaName, ".java", "")
		javaRpcImpl.GrpcFileName = grpcJavaName
	}
	for _, line := range codeLines {
		for i, sign := range javaRpcImpl.MethodSigns {
			if strings.Contains(line, sign.Method+"(") && strings.Contains(line, "java.util.Iterator") {
				sign.IsStream = ReturnSimplexStream
				javaRpcImpl.MethodSigns[i] = sign
			}
		}
	}
	return javaRpcImpl
}
