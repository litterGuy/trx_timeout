package main

import (
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
	"trx_timeout/utils"
)

func main() {
	t := time.NewTimer(time.Second)
	uri := "grpc.trongrid.io:50051"

	count := Count{}

	defer func() {
		if err := recover(); err != nil {
			log.Println("panic err is", err)
		}
	}()

	// init signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	for {
		select {
		case s := <-c:
			switch s {
			case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
				log.Printf("statistics data : %+v", count)
				log.Println("quit")
				return
			case syscall.SIGHUP:
			default:
				return
			}
		case <-t.C:
			count.Total = count.Total + 1
			blockCount, err := utils.GetNowBlock(uri)
			if err != nil {
				t.Reset(time.Second * 1)
				count.Error = count.Error + 1
				log.Println("rpc err :", err)
				continue
			}
			count.Success = count.Success + 1
			log.Println("now block num :" + strconv.FormatInt(blockCount, 10))
			t.Reset(0)
		}
	}
}

func init() {
	file := filepath.Join("", "trx_timeout.log")
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	log.SetOutput(logFile) // 将文件设置为log输出的文件
	log.SetPrefix("[trx_timeout]")
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)
	return
}

type Count struct {
	Total   int64
	Success int64
	Error   int64
}
