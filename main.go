package main

import (
	"encoding/json"
	"fmt"
	"github.com/liuzl/gocc"
	"log"
	"os"
	"time"
)

const rootPath = "/www/server/panel/BTPanel/languages"

func main() {
	// 判断文件夹是否存在
	if _, err := os.Stat(rootPath); os.IsNotExist(err) {
		log.Fatal(rootPath + "不存在")
	}
	// 重命名 all/zh.json 和 zh/server.json
	if _, err := os.Stat(fmt.Sprintf("%s/all/zh.json", rootPath)); err == nil {
		if err := os.Rename(fmt.Sprintf("%s/all/zh.json", rootPath), fmt.Sprintf("%s/all/zh-%d.json", rootPath, time.Now().Unix())); err != nil {
			log.Fatal("重命名all/zh.json失败:", err)
		}
	}
	if _, err := os.Stat(fmt.Sprintf("%s/zh/server.json", rootPath)); err == nil {
		if err := os.Rename(fmt.Sprintf("%s/zh/server.json", rootPath), fmt.Sprintf("%s/zh/server-%d.json", rootPath, time.Now().Unix())); err != nil {
			log.Fatal("重命名zh/server.json失败:", err)
		}
	}
	// 读取文件
	allzh, err := readJSONFile(fmt.Sprintf("%s/all/cht.json", rootPath))
	if err != nil {
		log.Fatal(err)
	}
	zhServer, err := readJSONFile(fmt.Sprintf("%s/cht/server.json", rootPath))
	if err != nil {
		log.Fatal(err)
	}
	// 初始化CC
	cc, err := gocc.New("t2s")
	if err != nil {
		log.Fatal(err)
	}
	// 转换
	allzhstr := convertToSimplified(cc, allzh)
	zhServerstr := convertToSimplified(cc, zhServer)
	// 写入文件
	err = writeJSONFile(fmt.Sprintf("%s/all/zh.json", rootPath), allzhstr)
	if err != nil {
		log.Fatal(err)
	}
	err = writeJSONFile(fmt.Sprintf("%s/zh/server.json", rootPath), zhServerstr)
	if err != nil {
		log.Fatal(err)
	}
	// 修改 settings.json
	settings, err := os.ReadFile(fmt.Sprintf("%s/settings.json", rootPath))
	if err != nil {
		log.Fatal("读取settings.json失败:", err)
	}
	type languageStruct struct {
		Name   string `json:"name"`
		Google string `json:"google"`
		Title  string `json:"title"`
		Cn     string `json:"cn"`
	}
	type settingsStruct struct {
		Default   string           `json:"default"`
		Languages []languageStruct `json:"languages"`
	}
	var settingsMap settingsStruct
	err = json.Unmarshal(settings, &settingsMap)
	if err != nil {
		log.Fatal("解析settings.json失败:", err)
	}
	exists := false
	for _, lang := range settingsMap.Languages {
		if lang.Name == "zh" {
			exists = true
			break
		}
	}
	// 插入zh-cn
	if !exists {
		settingsMap.Languages = append(settingsMap.Languages, languageStruct{
			Name:   "zh",
			Google: "zh-cn",
			Title:  "简体中文",
			Cn:     "简体中文",
		})
	}
	err = writeJSONFile(fmt.Sprintf("%s/settings.json", rootPath), settingsMap)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("转换完成")
}

func readJSONFile(filepath string) (map[string]any, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("读取%s失败: %w", filepath, err)
	}
	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("解析%s失败: %w", filepath, err)
	}
	return result, nil
}

func writeJSONFile(filepath string, data any) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("编码%s失败: %w", filepath, err)
	}
	err = os.WriteFile(filepath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("写入%s失败: %w", filepath, err)
	}
	return nil
}

func convertToSimplified(cc *gocc.OpenCC, data any) any {
	switch v := data.(type) {
	case map[string]any:
		for key, value := range v {
			v[key] = convertToSimplified(cc, value)
		}
		return v
	case []any:
		for i, item := range v {
			v[i] = convertToSimplified(cc, item)
		}
		return v
	case string:
		result, err := cc.Convert(v)
		if err != nil {
			log.Fatal("转换失败:", err)
		}
		return result
	default:
		return v
	}
}
