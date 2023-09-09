package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
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

// 检查名字
func CheckName(ipath string) string {
	fpath := "./output/" + ipath[8:]
	fpath = strings.ReplaceAll(fpath, "-", "_")
	lastIndex := strings.LastIndex(fpath, ".")
	fpath = fpath[:lastIndex] + ".MPF"
	fmt.Println("正在格式化 --> " + fpath[9:])
	// fmt.Println(fpath)
	return fpath
}

func FormatFile(path string) {
	// 检查可疑行
	CheckLineList := []string{}
	oldLine := ""
	lastLine := ""
	//	fmt.Println(path)
	fpath := CheckName(path)
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
		// 分离有效行
		if line[:1] != "N" {
			line = ""
		}
		// 程序第一行
		if strings.Index(line, "G64") != -1 {
			line = "N1 G54 G64 G0 G90"
			//			fmt.Println(line)
		}
		// 删除换刀
		if strings.Index(line, "T") != -1 {
			line = ""
		}
		// 删除G53
		if strings.Index(line, "G53") != -1 {
			line = strings.ReplaceAll(line, "G53", "")
			line = strings.TrimSpace(line)
			if strings.Index(line, " ") == -1 {
				line = ""
			}

		}
		// 删除G43
		if strings.Index(line, "G43") != -1 {
			line = strings.ReplaceAll(line, "G43 ", "")
			if strings.Index(line, "H") != -1 {
				index := strings.Index(line, "H")
				line = line[:index]
			}

		}
		// 删除G94
		if strings.Index(line, "G94") != -1 {
			line = strings.ReplaceAll(line, "G94 ", "")
		}
		// 修改进给速度为F3000
		if strings.Index(line, "F") != -1 {
			index := strings.Index(line, "F")
			line = line[:index] + "F3000."
		}
		if strings.Index(line, "M30") != -1 {
			line = "M9\r\n" + line
		}

		// 写入line
		if line != "" {
			line = line + "\r\n"
			// fmt.Print(line)
			if strings.Index(line, "Z") != -1 && (strings.Index(line, "X") != -1 || strings.Index(line, "Y") != -1) {
				// fmt.Print("old:" + oldLine)
				// fmt.Print("last:" + lastLine)
				if oldLine == lastLine {
					oldLine = ""
				} else {
					oldLine = "\r\n" + oldLine
				}
				CheckLineList = append(CheckLineList, oldLine)
				CheckLineList = append(CheckLineList, line)
				lastLine = line
			}
			fo.WriteString(line)
			oldLine = line
		}
	}
	if len(CheckLineList) > 20 {
		fmt.Println(fpath[9:] + "-->判断为精铣!!")
	} else if len(CheckLineList) > 0 && len(CheckLineList) < 21 {
		fmt.Println(fpath[9:] + "-->可能会撞刀!!请检查...")
		for i := 0; i < len(CheckLineList); i++ {
			fmt.Print(CheckLineList[i])
		}
	}
}

// 求文件的md5
func FileHash(data []byte) string {
	m := md5.New()
	m.Write(data)
	return hex.EncodeToString(m.Sum(nil))
}

// 比较文件大小
func BJM(A []string, B []string) {
	// c := 0
	for i := 0; i < len(B); i++ {
		for j := i + 1; j < len(B); j++ {
			if B[i] == B[j] {
				fmt.Println(A[i] + "与" + A[j] + "文件内容相同，请检查...")
			}
			// fmt.Println("比较了", c, "次")
			// c = c + 1
		}
	}
	fmt.Println()
}

func main() {
	BJA := []string{}
	BJB := []string{}
	fs := NcInit()
	for i := 0; i < len(fs); i++ {
		f, _ := os.Stat(fs[i])
		fb, _ := ioutil.ReadFile(fs[i])
		BJA = append(BJA, f.Name())
		BJB = append(BJB, FileHash(fb))

		FormatFile(fs[i])
		fmt.Println("")
	}
	BJM(BJA, BJB)

	fmt.Print("程序结束!!!")
	// 暂停
	zt := ""
	fmt.Scan(&zt)
	zt = ""
}
