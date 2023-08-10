package getter

import (
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	clog "unknwon.dev/clog/v2"

	"github.com/henson/proxypool/pkg/models"
)
//https://acq.iemoapi.com/getProxyIp?lb=1&return_type=txt&protocol=http&num=50
// iemoapi get ip from 66ip.cn
func YDL() (result []*models.IP, erro error)  {
	var ExprIP = regexp.MustCompile(`((25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\.){3}(25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\:([0-9]+)`)

	pollURL := "https://acq.iemoapi.com/getProxyIp?lb=1&return_type=txt&protocol=http&num=50"
	resp, err := http.Get(pollURL)
	if err != nil {
		log.Println(err)
		return
	}

	if resp.StatusCode != 200 {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	bodyIPs := string(body)
	ips := ExprIP.FindAllString(bodyIPs, 100)

	for index := 0; index < len(ips); index++ {
		ip := models.NewIP()
		ip.Data = strings.TrimSpace(ips[index])
		ip.Type1 = "http"
		ip.Source = "iemoapi"
		clog.Info("[iemoapi] ip = %s, type = %s", ip.Data, ip.Type1)
		result = append(result, ip)
	}

	clog.Info("iemoapi done.")
	return 
}
