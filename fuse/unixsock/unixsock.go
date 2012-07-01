package main
import (
	"fmt"
	"net"
	"math/rand"
	"os"
	"time"
)
func main() {
	if len(os.Args) != 2 {
		println("need one argument")
		return
	}
	ua, err := net.ResolveUnixAddr("unix", os.Args[1])
	if err != nil {
		println("resolve unix", err.Error())
		return
	}
	ul, err := net.ListenUnix("unix", ua)
	if err != nil {
		println("listen unix", err.Error())
		return
	}
	//for {
		uc, err := ul.AcceptUnix()
		if err != nil {
			println("accept", err.Error())
			panic("accept")
		}
		go func() {
			r := rand.Int()
			for {
				uc.Write([]byte(fmt.Sprintf("%d\n", r)))
				time.Sleep(time.Second)
			}
		}()
	//}
	ul.Close()
	<-(chan int)(nil)
}
