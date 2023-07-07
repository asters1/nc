package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// 判断文件夹是否存在
func IsExists(path string) {
	_, err := os.Stat(path)
	if err != nil {
		fmt.Println("文件夹[ " + path + " ]不存在!\n请自行创建文件夹!\n正在退出...")
		os.Exit(1)
	}
}

// 获得目录下的所有文件
func GetFiles(path string) []string {
	s := []string{}
	files, _ := ioutil.ReadDir(path)
	for _, file := range files {
		s = append(s, path+"/"+file.Name())
	}
	return s

}

// 初始化
func NcInit() []string {
	IsExists("./input")
	IsExists("./output")

	fs := GetFiles("./output")
	for i := 0; i < len(fs); i++ {
		os.RemoveAll(fs[i])
	}
	return GetFiles("./input")
}
func FormatFile(path string) {
	//	fmt.Println(path)
	fpath := "./output/" + path[8:]
	//	fmt.Println(f)
	fo, _ := os.OpenFile(fpath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	defer fo.Close()
	fo.WriteString("T1000\r\n")

	fi, _ := os.Open(path)
	defer fi.Close()
	buf := bufio.NewReader(fi)

	for {
		line, err := buf.ReadString('\n')
		if err == io.EOF {
			break
		}
		line = strings.TrimSpace(line)
		//分离有效行
		if line[:1] != "N" {
			line = ""
		}
		//程序第一行
		if strings.Index(line, "G64") != -1 {
			line = "N1 G54 G64 G0 G90"
			//			fmt.Println(line)
		}
		//删除换刀
		if strings.Index(line, "T") != -1 {
			line = ""
		}
		//删除G53
		if strings.Index(line, "G53") != -1 {
			line = strings.ReplaceAll(line, "G53", "")
			line = strings.TrimSpace(line)
			if strings.Index(line, " ") == -1 {
				line = ""
			}

		}
		//删除G43
		if strings.Index(line, "G43") != -1 {
			line = strings.ReplaceAll(line, "G43 ", "")
			if strings.Index(line, "H") != -1 {
				index := strings.Index(line, "H")
				line = line[:index]
			}

		}
		//删除G94
		if strings.Index(line, "G94") != -1 {
			line = strings.ReplaceAll(line, "G94 ", "")

		}
		//修改进给速度为F3000
		if strings.Index(line, "F") != -1 {
			index := strings.Index(line, "F")
			line = line[:index] + "F3000."
		}
		if strings.Index(line, "M30") != -1 {
			line = "M9\r\n" + line
		}

		//写入line
		if line != "" {
			line = line + "\r\n"
			fmt.Print(line)
			fo.WriteString(line)
		}
	}

}

func main() {
	fs := NcInit()
	for i := 0; i < len(fs); i++ {
		FormatFile(fs[i])

	}

}
