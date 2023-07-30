package conduit

// package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type OCRResponse struct {
	Result string `json:"result"`
	Error  string `json:"error"`
}

func OcrResult(imageBytes []byte) (string, error) {
	// OCR 运行
	baseLocal, _ := os.Getwd()
	ocrLocal := filepath.Join(baseLocal, "PowerSpider", "conduit", "ocr.py")

	// 读取图片文件并进行base64编码
	encodedImage := base64.StdEncoding.EncodeToString(imageBytes)

	// 运行Python解释器并执行ocr.py
	cmd := exec.Command("python", ocrLocal, encodedImage)
	cmd.Stderr = os.Stderr
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error executing Python script: %v", err)
	}

	// 解析ocr.py的输出JSON数据
	var response OCRResponse

	err = json.Unmarshal(output, &response)
	if err != nil {
		return "", fmt.Errorf("error decoding JSON response: %v", err)
	}

	// 处理识别结果
	if response.Error != "" {
		return "", fmt.Errorf("error: %s", response.Error)
	} else {
		return response.Result, nil
	}
}

func main() {
	// 读取图片文件并进行base64编码
	imageBytes, err := os.ReadFile(`G:\Gocode\PowerSpider\conduit\AuthCode.jpg`)
	if err != nil {
		fmt.Println("error reading image:", err)
		return
	}

	result, err := OcrResult(imageBytes)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println("识别结果是 => ", result)
}
