package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
)

var (
	C = flag.String("C", "#", "")
	D = flag.String("D", "", "列转换字典")
	F = flag.String("F", "1", "列")
)

var dictMap map[string]string

func init() {
	dictMap = make(map[string]string)
}

func transDict(filename string) {
	if filename == "" {
		return
	}
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	json.Unmarshal(bs, &dictMap)
}

func GetFiles(paths []string) []string {
	fs := make([]string, 0)
	for _, p := range paths {
		dir := path.Dir(p)
		base := path.Base(p)
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			fmt.Println(dir, "文件夹不存在", err.Error())
			continue
		}

		for _, file := range files {
			fname := file.Name()
			isDir := file.IsDir()
			if isDir {
				continue
			}
			match, _ := path.Match(base, fname)
			if match {
				fs = append(fs, path.Join(dir, fname))
			}
		}
	}
	return fs
}

func main() {
	flag.Parse()
	files := GetFiles(flag.Args())
	lenth := len(files)
	if lenth == 0 {
		fmt.Printf("未找到文件,请检查参数是否在正确 \n")
		return
	}
	fcol, err := strconv.Atoi(*F)
	if err != nil {
		fmt.Println("列指定错误")
		return
	}
	fcol = fcol - 1
	if fcol < 0 {
		fcol = 0
	}

	//解析字典
	transDict(*D)
	endl := 0
	fps := make([]*(os.File), 0)
	rdMap := make(map[int]*(bufio.Reader))
	ssArray := make([]*string, lenth)
	for index, filePath := range files {
		fp, err := os.Open(filePath)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fps = append(fps, fp)
		rdMap[index] = bufio.NewReader(fp)
	}

	defer func() {
		for _, f := range fps {
			f.Close()
		}
	}()

	for {
		if endl >= lenth {
			break
		}
		for i := 0; i < lenth; i++ {
			rd := rdMap[i]
			if rd == nil {
				//fmt.Printf("%-9s ", *C)
				ssArray[i] = C
				continue
			}
			line, err1 := rd.ReadString('\n')
			if err1 != nil || io.EOF == err1 {
				endl++
				rdMap[i] = nil
				//fmt.Printf("%-9s ", *C)
				ssArray[i] = C
				continue
			}
			line = strings.TrimSpace(line)
			//fmt.Printf("%-9s ", line)
			ssArray[i] = &line
		}
		flagMark := 0
		for _, str := range ssArray {
			if *str != *C {
				flagMark++
				break
			}
			//fmt.Printf("%-9s ", *str)
		}
		if flagMark > 0 {
			for _, line := range ssArray {
				//切片
				strList := strings.Fields(*line)
				lstr := len(strList)
				if fcol <= lstr-1 {
					colStr := strings.TrimSpace(strList[fcol])
					strList[fcol] = getColStr(colStr)
					for _, cols := range strList {
						fmt.Printf("%-8s ", cols)
					}
				}
			}
			fmt.Printf("\n")
		}
	}

}

func getColStr(src string) string {
	targetStr := dictMap[src]
	if len(targetStr) == 0 {
		return src
	}
	return targetStr
}
