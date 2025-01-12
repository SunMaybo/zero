package release

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SunMaybo/zero/zctl/cmd"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/c-bata/go-prompt"
	"go.uber.org/zap"
	"io/fs"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var envWhiteList = map[string]bool{
	"qa39":    false,
	"qa40":    false,
	"qa41":    false,
	"qa53":    false,
	"qa67":    false,
	"qa76":    false,
	"qa86":    false,
	"sandbox": false,
	"format":  true,
	"all":     false,
}
var dingTalkToken = "44c8b95ce7bb6fa11f674598a2c7f60e782d809e75e2dd6a475edba4d43ebf46"

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

type OssConfig struct {
	Bucket    string `json:"bucket"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
}

func Delay(env string, path, cdnUrl, bucket string, isScale bool, pk, cdnPk, dingTalkSecret string, isSaas bool, isCommon bool) {
	if pk == "" {
		log.Println("please config front_web_pk on .zctl.yaml")
	}
	if dingTalkSecret == "" {
		log.Fatal("please config ding_talk_secret on .zctl.yaml")
	}
	env = strings.TrimSpace(env)
	//1. 指定环境
	//2. 校验环境白名单
	var ok, isOnline bool
	if isOnline, ok = envWhiteList[strings.ToLower(env)]; !ok {
		log.Fatal("current env is no existed")
	}
	if isOnline {
		log.Println("check your env is online", "evn", env)
	} else {
		log.Println("check your env is offline", "evn", env)
	}

	branchName := getGitBranchName(path)
	if branchName == "" {
		log.Fatal("current branch is nil errs")
	}
	log.Println("checkout current branch name is " + branchName)
	//3. 校验当前分支是否合法
	if isOnline && !isScale && branchName != "master" {
		log.Fatal("you must publish online through the Master branch")
	}
	//4. 获取当前项目名称

	projectName := getProjectName(path)
	if projectName == "" {
		log.Fatal("git project name err")
	}
	log.Println("checkout current project name is " + projectName)
	delayDir := projectName
	version := time.Now().Format("01021504")
	log.Println("please input your password")

	passowrd := ""
	var err error
	if passowrd, err = input("password"); err != nil {
		log.Fatal(err)
	}
	var cfgs []OssConfig
	accessKey := ""
	secretKey := ""
	endpoint := "oss-cn-beijing.aliyuncs.com"
	if result, err := DecryptByAes(passowrd, pk); err != nil {
		log.Fatalf("you entered the password incorrectly")
	} else {
		err = json.Unmarshal([]byte(result), &cfgs)
		if err != nil {
			log.Fatal(err)
		}
		for _, cfig := range cfgs {
			if cfig.Bucket == bucket {
				log.Println("very good,please waiting......")
				accessKey = cfig.AccessKey
				secretKey = cfig.SecretKey
				break
			}
		}

	}
	if accessKey == "" {
		log.Fatalf("you oss Ak config err")
	}
	cdnAccessKey := ""
	cdnSecretKey := ""
	if result, err := DecryptByAes(passowrd, cdnPk); err != nil {
		log.Fatalf("you entered the password incorrectly")
	} else {
		log.Println("very good,please waiting......")
		cdnAccessKey = strings.Split(string(result), "-")[0]
		cdnSecretKey = strings.Split(string(result), "-")[1]
	}

	var delayDirs []string
	//4. 线上版本打tag并上传
	if isOnline && !isScale {
		if result, err := cmd.Run("git tag release-"+version, path); err != nil {
			log.Fatalf("git tag release-%s,err:%s", version, err.Error())
		} else {
			log.Println(result)
			log.Println("git tag release-%s success " + version)
		}
		if result, err := cmd.Run("git push origin release-"+version, path); err != nil {
			log.Fatalf("git push origin release-%s,err:%s", version, err.Error())
		} else {
			log.Println(result)
			log.Println("git push origin release-%s success " + version)
		}
		delayDirs = append(delayDirs, delayDir)
	} else {
		if "all" == env {
			for str := range envWhiteList {
				if strings.Contains(str, "qa") {
					delayDirs = append(delayDirs, delayDir+"-"+str)
				}
			}
		} else {
			delayDirs = append(delayDirs, delayDir+"-"+env)
		}
	}
	for _, delayDir := range delayDirs {
		//3. 前端项目build
		if isOnline {
			log.Printf("Please enter Yes to confirm online project %s with version %s.", projectName, "release-"+version)
			if confirm, err := input("ensure"); err != nil {
				log.Fatal(err)
			} else if strings.TrimSpace(confirm) != "yes" {
				log.Println("You have been terminated.")
			}
			// SaaS 项目 pathname
			var npmDelayDir string
			if isSaas {
				npmDelayDir = strings.Split(delayDir, "-")[1]
			} else if isCommon {
				npmDelayDir = strings.Split(delayDir, "-")[1]
			}
			if err := cmd.Execute("/bin/bash", path, func(lines string) {
				log.Println(lines)
			}, "-c", "npm i --registry https://registry.npm.taobao.org  && npm run build --projectdir="+npmDelayDir); err != nil {
				log.Fatalf("npm build err:%s", err.Error())
			} else {
				log.Println("npm run build success")
			}
		} else {
			if err := cmd.Execute("/bin/bash", path, func(lines string) {
				log.Println(lines)
			}, "-c", "npm i --registry https://registry.npm.taobao.org  && npm run build --projectdir="+delayDir); err != nil {
				log.Fatalf("npm build err:%s", err.Error())
			} else {
				log.Println("npm run build success")
			}
		}
		client, err := oss.New(endpoint, accessKey, secretKey)
		if err != nil {
			log.Println("Error:", err)
			os.Exit(-1)
		}

		// 填写存储空间名称，例如examplebucket。
		bucket, err := client.Bucket(bucket)
		if err != nil {
			log.Fatalf("Error:%v", err)
			os.Exit(-1)
		}
		currentUser, _ := user.Current()
		username := currentUser.Username
		if isOnline && (!isScale || !isCommon) {
			err := DingTalkNew(dingTalkSecret, dingTalkToken).
				Talk("【前端项目发布通知】", fmt.Sprintf("[*%s*同学～上线了前端---%s---项目---当前版本---%s]", username, delayDir, "release-"+version), nil, nil, true)
			if err != nil {
				log.Fatal("ding talk err abort publish")
			}
		} else if isOnline {
			err := DingTalkNew(dingTalkSecret, dingTalkToken).
				Talk("【前端项目发布通知】", fmt.Sprintf("[*%s*同学～回滚了前端---%s---项目---当前版本---%s]", username, delayDir, branchName), nil, nil, true)
			if err != nil {
				log.Fatal("ding talk err abort publish")
			}
		}
		//4. 代码上传
		if isOnline && (isSaas || isCommon) {
			prefix := "prod/"
			if isSaas {
				prefix = prefix + "fe-xbbcloud/"
			} else if isCommon {
				prefix = prefix + "fe-common/"
			}
			delayDir = strings.Split(delayDir, "-")[1]
			fmt.Println(prefix + delayDir)
			uploadDirectoryFileTree(bucket, path+"/dist", prefix+delayDir)
		} else if isOnline {
			prefix := "prod/"
			fmt.Println(prefix + delayDir + "false")
			uploadDirectoryFileTree(bucket, path+"/dist", prefix+delayDir)
		} else {
			uploadDirectoryFileTree(bucket, path+"/dist", "test/"+delayDir)
		}
	}
	//5.线上
	zap.S().Info("uploader success.....")
	if isOnline && cdnUrl != "" {
		//6.刷新cdn
		err = RefreshCdn(cdnUrl, cdnAccessKey, cdnSecretKey)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("refresh cdn[%s] success.....", cdnUrl)
	}

}

type uploadFile struct {
	objectKey string
	filepath  string
}

func uploadDirectoryFileTree(bucket *oss.Bucket, contextPath, output string) {
	uploadFiles := make(chan uploadFile, 20)
	concurrent := 3
	wait := sync.WaitGroup{}
	wait.Add(concurrent)
	for i := 0; i < concurrent; i++ {
		go func() {
			defer wait.Done()
			for file := range uploadFiles {
				if err := bucket.PutObjectFromFile(file.objectKey, file.filepath); err != nil {
					zap.S().Fatalf("uploader err,%s", err.Error())
				}
				if err := bucket.SetObjectACL(file.objectKey, oss.ACLPublicRead); err != nil {
					zap.S().Fatalf("uploader err,%s", err.Error())
				}
			}
		}()
	}
	if err := filepath.Walk(contextPath, func(path string, info fs.FileInfo, err error) error {
		if contextPath == path {
			return nil
		}
		if !info.IsDir() {
			uploadFiles <- uploadFile{
				objectKey: output + path[len(contextPath):],
				filepath:  path,
			}
			return nil
		}
		return nil
	}); err != nil {
		zap.S().Error(err)
	}
	close(uploadFiles)
	wait.Wait()
}
func getProjectName(path string) string {
	if result, err := cmd.Run("git remote show origin |grep \"Push  URL\"", path); err != nil {
		zap.S().Fatalf("git remote show origin err,%s", err.Error())
	} else {
		for _, s := range strings.Split(result, "/") {
			if strings.HasSuffix(strings.TrimSpace(s), ".git") {
				return strings.ReplaceAll(strings.TrimSpace(strings.ReplaceAll(s, ".git", "")), "\n", "")
			}
		}
	}
	return ""
}

func getGitBranchName(dir string) string {
	if result, err := cmd.Run("git branch", dir); err != nil {
		zap.S().Fatal("git branch err", err)
	} else {
		items := strings.Split(DBC2SBC(result), "\n")
		for _, item := range items {
			if len(item) <= 0 {
				return ""
			}
			if item[0] == 42 {
				item = item[1:]
				return strings.TrimSpace(item)
			}
		}
	}
	return ""
}
func DBC2SBC(s string) string {
	var strLst []string
	for _, i := range s {
		insideCode := i
		if insideCode == 12288 {
			insideCode = 32
		} else {
			insideCode -= 65248
		}
		if insideCode < 32 || insideCode > 126 {
			strLst = append(strLst, string(i))
		} else {
			strLst = append(strLst, string(insideCode))
		}
	}
	return strings.Join(strLst, "")
}
func input(prefix string) (string, error) {
	result := prompt.Input(prefix+"> ", completer)
	if result == "exit" {
		fmt.Println("---------------------------程序退出-----------------------------------")
		time.Sleep(3 * time.Second)
		return "", errors.New("exit")
	}
	return result, nil
}
