package main

import (
	"log"
	"runtime"
	"sync"
	"time"
	"fmt"
	"github.com/henson/proxypool/api"
	"github.com/henson/proxypool/getter"
	"github.com/henson/proxypool/pkg/initial"
	"github.com/henson/proxypool/pkg/models"
	"github.com/henson/proxypool/pkg/storage"

)



func main() {

	//init the database
	initial.GlobalInit()

	runtime.GOMAXPROCS(runtime.NumCPU())
	ipChan := make(chan *models.IP, 2000)

	// Start HTTP
	go func() {
		api.Run()
	}()

	// Check the IPs in DB
	go func() {
		storage.CheckProxyDB()
	}()
		//tunnel translation

	go func(){
		api.Translation(":8881")
	}()
	for i := 0; i < 10; i++ {
		go func() {
			for {

				fmt.Println("-------------------begin--------------------------")
				time.Sleep(10 * time.Minute)
				storage.Rcheckp()
				
				fmt.Println("-------------------end--------------------------")
				
				
			}
		}()
	}
	
	
	// Check the IPs in channel
	for i := 0; i < 50; i++ {
		go func() {
			for {
				storage.CheckProxy(<-ipChan)
			}
		}()

	}

	// Start getters to scraper IP and put it in channel
	for {
		n := models.CountIPs()
		log.Printf("Chan: %v, IP: %v\n", len(ipChan), n)
		if len(ipChan) < 100 {
			go run(ipChan)
		}
		time.Sleep(10 * time.Minute)
	}
}

func run(ipChan chan<- *models.IP) {
	var wg sync.WaitGroup
	funs := []func() ([]*models.IP, error){
		//getter.FQDL,  //c
		//getter.PZZQZ, //新代理
		//getter.Data5u, //c
		//getter.Feiyi, //c
		//getter.IP66, //need to remove it
		//getter.IP3306,
		//getter.KDL,
		//getter.GBJ,	//因为网站限制，无法正常下载数据
		//getter.Xici, //c
		//getter.XDL, //close
		//getter.IP181,  // 已经无法使用
		getter.YDL,	
		//getter.PLP, //need to remove it
		//getter.PLPSSL,
		//getter.IP89,
	}
	for _, f := range funs {
		wg.Add(1)
		go func(f func() ([]*models.IP, error)) {
		
		
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Recovered from panic in getter function: %v", r)
				}
			}()
			temp, err := f()
			if err != nil {
				// Handle error here. For example, you can just log it:
				log.Printf("Error in getter function: %v", err)
			
			} else {
			//log.Println("[run] get into loop")
			for _, v := range temp {
				//log.Println("[run] len of ipChan %v",v)
				ipChan <- v
			}
			wg.Done()
			}
		}(f)
	}
	wg.Wait()
	log.Println("All getters finished.")
}
