package reader

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"
	"regexp"
)

//当新文件出现时，是否自动读取
func TestCheck(t *testing.T) {
	util(true)
}

//测试文件打开关闭
func TestStartAndStop(t *testing.T) {
	util(false)
}

func util(isnext bool) {
	stream := make(chan string, 100)
	rj, err := NewReader("/Users/anbaoyong/Project/test/aby.${%Y-%m-%d-%H}", stream, "")
	if err != nil {
		return
	}
	go rj.Start()
	go func() {
		time.Sleep(2 * time.Second) //2秒后创建文件
		now := time.Now()
		if isnext {
			now = now.Add(time.Hour)
		}
		suffix := now.Format("2006-01-02-15")
		filepath := fmt.Sprintf("/Users/anbaoyong/Project/test/aby.%s", suffix)

		{
			f, err := os.Create(filepath)
			if err != nil {
				log.Fatal(err)
			}
			time.Sleep(time.Millisecond * 250) //因为tail 的巡检周期是250ms
			defer f.Close()

			fmt.Fprint(f, "this is a test\n")

		}
		time.Sleep(250 * time.Millisecond) //延迟关闭
		rj.Stop()

	}()

	for line := range stream {
		fmt.Println(line)
	}
}

func TestRegex(t *testing.T) {
	s := "2018-10-15 09wewe:21:22,769 ERROR controller.AlarmController - send alarm error:"
	reg := `(2[0-9]{3})-(0[1-9]|1[012])-([012][0-9]|3[01])\s([01][0-9]|2[0-4])(:[012345][0-9]){2},\d+`
	rr, _ := regexp.Compile(reg)
	b := rr.Find([]byte(s))
	if b == nil {
		fmt.Println("sadfsadfasdf")
	}
	fmt.Println(string(b))
}
