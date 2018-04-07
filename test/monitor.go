package test

import (
	"io/ioutil"
	"os"
	"runtime/pprof"
	"strings"
	"time"
)

/**
test为这些情况提供了一个统一的监控功能，在单元测试运行目录下使用以下方法可以获取到单元测试过程中的信息：

echo 'lookup goroutine' > utest.cmd
以上shell脚本将使utest自动输出goroutine堆栈跟踪信息到utest.goroutine文件。

utest支持以下几种监控命令：

lookup goroutine  -  获取当前所有goroutine的堆栈跟踪信息，输出到utest.goroutine文件，用于排查死锁等情况
lookup heap       -  获取当前内存状态信息，输出到utest.heap文件，包含内存分配情况和GC暂停时间等
lookup threadcreate - 获取当前线程创建信息，输出到utest.thread文件，通常用来排查CGO的线程使用情况
此外你还可以通过注册utest.CommandHandler回调来添加自己的监控命令支持。
 */
// Command handler
var CommandHandler func(string) bool

func init() {
	go func() {
		for {
			if input, err := ioutil.ReadFile("utest.cmd"); err == nil && len(input) > 0 {
				ioutil.WriteFile("utest.cmd", []byte(""), 0744)

				cmd := strings.Trim(string(input), " \n\r\t")

				var (
					profile  *pprof.Profile
					filename string
				)

				switch cmd {
				case "lookup goroutine":
					profile = pprof.Lookup("goroutine")
					filename = "utest.goroutine"
				case "lookup heap":
					profile = pprof.Lookup("heap")
					filename = "utest.heap"
				case "lookup threadcreate":
					profile = pprof.Lookup("threadcreate")
					filename = "utest.thread"
				default:
					if CommandHandler == nil || !CommandHandler(cmd) {
						println("unknow command: '" + cmd + "'")
					}
				}

				if profile != nil {
					file, err := os.Create(filename)
					if err != nil {
						println("couldn't create " + filename)
					} else {
						profile.WriteTo(file, 2)
					}
				}
			}
			time.Sleep(2 * time.Second)
		}
	}()
}
