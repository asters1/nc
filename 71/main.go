package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
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
	IsExists("./bak")

	fs := GetFiles("./output")
	for i := 0; i < len(fs); i++ {
		os.RemoveAll(fs[i])
	}
	return GetFiles("./input")
}

func ClearInput() {
	fs := GetFiles("./input")
	// fmt.Println("ClearInput")
	for i := 0; i < len(fs); i++ {
		os.RemoveAll(fs[i])
	}
}

// 检查名字
func CheckName(ipath string) string {
	if strings.HasSuffix(ipath, ".NC") || strings.HasSuffix(ipath, ".nc") || strings.HasSuffix(ipath, ".Nc") || strings.HasSuffix(ipath, ".nC") {
		fmt.Println("正在格式化 --> " + ipath)
		return strings.ReplaceAll(ipath, "input", "output")
	}
	return ""
}

func FormatFile(path string) {
	// 检查可疑行
	CheckLineList := []string{}
	oldLine := ""
	lastLine := ""
	//	fmt.Println(path)
	fpath := CheckName(path)
	if fpath == "" {
		return
	}
	//	fmt.Println(f)
	fo, _ := os.OpenFile(fpath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	defer fo.Close()
	fo.WriteString("%\r\n1000\r\n")

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
			if !(strings.Contains(line, "TOOL")) {
				line = ""
			}
			// fmt.Println(line)
		}
		// fmt.Println(line)
		// 删除换刀
		if strings.Index(line, "T") != -1 {
			if !strings.Contains(line, "TOOL") {
				line = "G5.1Q1"
			}
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

func GetTimeName() string {
	t := time.Now()
	// fmt.Println(T)
	// fmt.Println(s)
	s := strconv.Itoa(t.Year()) + "年" + strconv.Itoa(int(t.Month())) + "月" + strconv.Itoa(t.Day()) + "日" + strconv.Itoa(t.Hour()) + "时" + strconv.Itoa(t.Minute()) + "分" + strconv.Itoa(t.Second()) + "秒"
	return s
}

// 备份
func Bak(fs []string) {
	if len(fs) == 0 {
		os.Exit(1)
	}
	str_dir := "./bak/" + fs[0][8:] + "_" + GetTimeName() + "/"

	os.MkdirAll(str_dir, 0666)
	for i := 0; i < len(fs); i++ {
		CopyFile(fs[i], str_dir+fs[i][8:])
	}
}

// CopyFile 拷贝文件函数
func CopyFile(inputName, outputName string) (written int64, err error) {
	src, err := os.Open(inputName)
	if err != nil {
		fmt.Printf("打开 %s 失败!, 错误:%v.\n", inputName, err)
		return
	}
	defer src.Close()
	dst, err := os.OpenFile(outputName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("打开 %s 失败, 错误:%v.\n", outputName, err)
		return
	}
	defer dst.Close()
	return io.Copy(dst, src)
}

func formatXY(l string) string {
	l_list := strings.Split(l, " ")
	// fmt.Println(l_list)
	line := ""
	for i := 0; i < len(l_list); i++ {
		if strings.Contains(l_list[i], "X") {
			line = line + l_list[i] + " "
		}
		if strings.Contains(l_list[i], "Y") {
			line = line + l_list[i]
		}
	}
	return line
}

func ZKXH(path string) {
	frist_x_y_switch := true
	fpath := CheckName(path)
	// fmt.Println(fpath)
	fo, _ := os.OpenFile(fpath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	defer fo.Close()
	fo.WriteString("%\r\nO1000\r\nG49 G64 G17 G80 G0 G90 G40 G99\r\nG5.1Q1\r\nZ100.\r\n")

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
		line = formatXY(line)
		if line != "" {
			if frist_x_y_switch {
				if strings.Contains(line, "X") && !strings.Contains(line, "Y") {
					line = line + " Y0."
				}
				if strings.Contains(line, "Y") && !strings.Contains(line, "X") {
					line = line + " X0."
				}
				line = line + " S800 M3\r\nG98 G83 Z-3. R2 Q3 F100."
				frist_x_y_switch = false
			}
		}
		if line != "" {
			line = line + "\r\n"
			fo.WriteString(line)
		}

	}
	fo.WriteString("G80\r\nZ100.\r\nM9\r\nM30\r\n")
}

func main() {
	fs := NcInit()
	Bak(fs)

	BJA := []string{}
	BJB := []string{}
	for i := 0; i < len(fs); i++ {
		f, _ := os.Stat(fs[i])
		fb, _ := ioutil.ReadFile(fs[i])
		BJA = append(BJA, f.Name())
		BJB = append(BJB, FileHash(fb))

		if strings.Contains(strings.ToLower(fs[i]), strings.ToLower("ZKXH")) {
			ZKXH(fs[i])
		} else {
			FormatFile(fs[i])
		}
		fmt.Println("")
	}
	BJM(BJA, BJB)
	ClearInput()

	fmt.Print("程序结束!!!")
	// 暂停
	zt := ""
	fmt.Scan(&zt)
	zt = ""
}
