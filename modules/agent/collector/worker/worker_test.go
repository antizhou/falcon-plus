package worker

import (
	"fmt"
	"testing"
	"time"
)

func TestWorkerStart(t *testing.T) {
	c := make(chan string, 10)
	go func() {
		for i := 0; i < 1000; i++ {
			for j := 0; j < 10; j++ {
				c <- fmt.Sprintf("memeda--%d--%d", i, j)
			}
			fmt.Println()
			time.Sleep(time.Second * 1)
		}
	}()
	wg := NewWorkerGroup("memeda", c, nil)
	wg.Start()
	time.Sleep(10 * time.Second)
	wg.Stop()
	time.Sleep(1 * time.Second)
}

func TestWorkerStartsdd(t *testing.T) {
	//fmt.Println(time.Now().Format("2006-01-02 15:04:05.999"))
	//fmt.Println(os.O_APPEND|os.O_CREATE|os.O_WRONLY)
	var pat = ""
	pat = `([012][0-9]|3[01])-[JFMASOND][a-z]{2}-(2[0-9]{3})\s([01][0-9]|2[0-4])(:[012345][0-9]){2}`
	fmt.Println(pat)

	pat = "([012][0-9]|3[01])-[JFMASOND][a-z]{2}-(2[0-9]{3})\\s([01][0-9]|2[0-4])(:[012345][0-9]){2}"
	fmt.Println(pat)

}
