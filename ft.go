package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

var (
	C = flag.String("C", "#", "")
)

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
	endl := 0
	fps := make([]*(os.File), 0)
	rdMap := make(map[int]*(bufio.Reader))
	for index, filePath := range files {
		fp, err := os.Open(filePath)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fps = append(fps, fp)
		rdMap[index] = bufio.NewReader(fp)
	}

	for {
		if endl >= lenth {
			break
		}
		for i := 0; i < lenth; i++ {
			rd := rdMap[i]
			if rd == nil {
				fmt.Printf("%-9s ", *C)
				continue
			}
			line, err1 := rd.ReadString('\n')
			if err1 != nil || io.EOF == err1 {
				endl++
				rdMap[i] = nil
				fmt.Printf("%-9s ", *C)
				continue
			}
			line = strings.TrimSpace(line)
			fmt.Printf("%-9s ", line)
		}
		fmt.Printf("\n")

	}

}
