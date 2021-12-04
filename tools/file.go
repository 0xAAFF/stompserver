package tools

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

//CreateDir  文件夹创建
func CreateDir(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	os.Chmod(path, os.ModePerm)
	return nil
}

// ExistsDirectory  文件夹是否存在
//  参数:
//  path	string	检测的目录地址
//  return	bool	目录是否存在
func ExistsDirectory(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// ExistsFile  文件是否存在
//  参数:
//  path	string	检测的文件地址
func ExistsFile(path string) bool {
	if path == "" {
		return false
	}
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !s.IsDir()
}

// GetCurrentPath  获取当前程序运行路径
//  return	string,error 当前运行路径,err
//  例如: E:/Go/ProjectTest/ <nil>
func GetCurrentPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	path, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}
	if runtime.GOOS == "windows" {
		path = strings.Replace(path, "\\", "/", -1)
	}
	i := strings.LastIndex(path, "/")
	if i < 0 {
		return "", errors.New(`Can't find "/" or "\".`)
	}
	return string(path[0 : i+1]), nil
}

// GetCurrentPathNoError  获取当前程序运行路径
//  return	string	当前运行路径,如果获取失败返回""
//  例如: E:/Go/ProjectTest/ <nil>
func GetCurrentPathNoError() string {
	path, err := GetCurrentPath()
	if err != nil {
		return ""
	}
	return path
}

func SaveBinary(filePath string, data []byte) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	file.Write(data)
	return nil
}
