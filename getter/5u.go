package getter

import (
	"github.com/Aiicy/htmlquery"
	"github.com/henson/proxypool/pkg/models"
	clog "unknwon.dev/clog/v2"
	//"os"
	//"fmt"
	"strings"
)

//Data5u is not work now
// Data5u get ip from data5u.com
func Data5u() (result []*models.IP) {
	pollURL := "https://www.proxy-list.download/SOCKS5"
		
	doc, _ := htmlquery.LoadURL(pollURL)
	//fmt.Println(doc)
	//os.Exit(0)
	trNode, err := htmlquery.Find(doc, "//table[@id='example1']//tbody//tr")
	
	if err != nil {
		clog.Warn(err.Error())
	}
	for i := 0; i < len(trNode); i++ {
		tdNode, _ := htmlquery.Find(trNode[i], "//td")


		ip := strings.TrimSpace(htmlquery.InnerText(tdNode[0]))
		port := strings.TrimSpace(htmlquery.InnerText(tdNode[1]))
		Type := "Socks5"
		speed := strings.TrimSpace(htmlquery.InnerText(tdNode[4]))

		IP := models.NewIP()
		IP.Data = ip + ":" + port
		IP.Type1 = Type
		IP.Source = "proxy-list"
		IP.Speed = extractSpeed(speed)
		//fmt.Println(IP.Data)
		//os.Exit(0)
		result = append(result, IP)
	}

	clog.Info("[fanqiedl] done")
	return
}