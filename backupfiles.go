/*
* @Author: shibengen
* @Date:   2018-04-28 14:59:28
* @Last Modified by:   shibengen
* @Last Modified time: 2018-05-04 17:37:06
 */
package main

import (
	"flag"
	"fmt"
	cp "github.com/hacdias/fileutils"
	lconfig "github.com/larspensjo/config"
	termbox "github.com/nsf/termbox-go"
	"log"
	"os"
	"runtime"
	"strconv"
	"time"
)

var (
	configFile = flag.String("configfile", "config.conf", "General configuration file")
	conf       = make(map[string]string)
)

//读取文件需要经常进行错误检查，这个帮助方法可以精简下面的错误检查过程。
func check(e error) {
	if e != nil {
		panic(e)
	}
}

//
func pause() {
	fmt.Println("请按任意键继续...")
Loop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			break Loop
		}
	}
}

//
func init() {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	termbox.SetCursor(0, 0)
	termbox.HideCursor()
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	cfg, err := lconfig.ReadDefault(*configFile) //读取配置文件，并返回其Config
	if err != nil {
		log.Fatalf("Fail to find %v,%v", *configFile, err)
		pause()
	}
	if cfg.HasSection("path") { //判断配置文件中是否有一级标签
		options, err := cfg.SectionOptions("path") //获取一级标签的所有子标签（只有标签没有值）
		if err == nil {
			for _, v := range options {
				optionValue, err := cfg.String("path", v) //根据一级标签和二级标签获取对应的值
				if err == nil {
					conf[v] = optionValue
				}
			}
		}
	}
	days, _ := strconv.ParseInt(conf["delete_day"], 10, 64)
	timestamp := time.Now().Unix() - (86400 * days)
	mode := conf["mode"]
	basepath := conf["to_dir"]
	subdir := time.Now().Format("2006_01_02")
	if mode == "hour" {
		subdir = subdir + "/" + time.Unix(timestamp, 0).Format("2006_01_02_15")
	} else if mode == "minute" {
		subdir = subdir + "/" + time.Unix(timestamp, 0).Format("2006_01_02_15_04")
	}
	//copy
	output := cp.CopyDir(conf["from_dir"], conf["to_dir"]+"/"+subdir)
	//删除N天前的备份 2006-01-02 15:04:05 默认按天删除
	del_path := basepath + "/" + time.Unix(timestamp, 0).Format("2006_01_02")
	os.RemoveAll(del_path)
	if output != nil {
		log.Println(del_path)
		log.Println(subdir)
		log.Println(output)
		pause()
	}

}
