// +build main
package main

import (
	"flag"
	"fmt"
	"github.com/guonaihong/gutil/file"
	"github.com/oraleval/ulog"
	"io"
	"os"
)

func myError(u *ulog.Ulog) {
	u.Error(1).ID("request-id").Msg("聚合错误函数，测试跳过当前调用栈功能")
}

func main() {
	//压缩文件前缀,
	prefix := flag.String("prefix", "ws-eval" /*根据不同的项目修改这里*/, "eval log file prefix")
	//压缩文件存放文件夹
	dir := flag.String("dir", "./http-eval-log/" /*根据不同的项目修改这里*/, "http eval log save directory")
	//单个日志文件最大可以达到多大
	maxSize := flag.String("max-size", "1G" /*可以调整单个文件大小*/, "Single log file size")
	//最多几个压缩文件
	maxArchive := flag.Int("max-archive", 5 /*可以调整保存归档压缩文件数*/, "How many compressed files to save at most")
	//是否写日志文件到文件系统中
	save := flag.Bool("save", false /*这里不需要修改，直接在命令行里面配置*/, "Whether to save the log to the hard disk")

	flag.Parse()

	w := []io.Writer{os.Stdout}
	if *save {
		if size, err := file.ParseSize(*maxSize); err != nil {
			fmt.Errorf("Invalid value -max-size %s, %s\n", *maxSize, err)
		} else {

			// 需要保证声明周期比较长
			file := ulog.NewFile(*prefix, *dir, ulog.Gzip, int(size), *maxArchive)
			w = append(w, file)
			defer file.Close()
		}
	}

	u := ulog.New(w...)

	u.Debug().ID("request-id").IP("127.0.0.1:8080").Msg("hello")
	u.Error().ID("request-id").IP("127.0.0.1:8080").Msg("test error")
	u.Warn().Caller().Msg("test warn")
	myError(u)
}
