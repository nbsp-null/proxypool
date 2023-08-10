package getter

import (
	"io/ioutil"
	"net/http"

	//"fmt"
	//"os"
	clog "unknwon.dev/clog/v2"
	"time"
	//"regexp"
	"strings"
	"encoding/json"
	"github.com/henson/proxypool/pkg/models"
)

type Proxy struct {
	Addr  string `json:"addr"`
	Type1 int    `json:"type"`
}

//IP89 get ip from www.89ip.cn
func FQDL()  (result []*models.IP, erro error) {
	clog.Info("checkerproxy] start test")
	currentTime := time.Now()
	//var ExprIP = regexp.MustCompile(`((25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\.){3}(25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\:([0-9]+)`)
	pollURL := "https://checkerproxy.net/api/archive/"+currentTime.Format("2006-01-02")
	
	resp, err := http.Get(pollURL)
	if err != nil {
		clog.Warn(err.Error())
		return
	}

	if resp.StatusCode != 200 {
		clog.Warn(err.Error())
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	bodyIPs := strings.TrimSpace(string(body))
	var proxy []Proxy
	err1 := json.Unmarshal([]byte(bodyIPs), &proxy)
	if err1 != nil {
		//fmt.Println("Failed to parse JSON:", err1)
		return
	}
	//fmt.Println("-----------------------------------------")
	//fmt.Println(proxy[0])
	//os.Exit(0)
	//ips := ExprIP.FindAllString(bodyIPs, 100)

	for index := 0; index < len(proxy); index++ {
		ip := models.NewIP()
		ip.Data = strings.TrimSpace(proxy[index].Addr)
		ip.Type1 = "https"
		if proxy[index].Type1==1  {
		
		}else if proxy[index].Type1==4  {
		ip.Type1 = "Socks5"
		} else {
		ip.Type1 = "https"
		}
		

		ip.Source = "checkerproxy"
		clog.Info("[checkerproxy] ip = %s, type = %s", ip.Data, ip.Type1)
		result = append(result, ip)
	}

	clog.Info("checkerproxy done.")
	return
}
