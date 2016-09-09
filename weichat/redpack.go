package weichat

import (
	"bytes"
	"crypto/md5"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type Bonus struct {
	XMLName      xml.Name `xml:"xml"`
	Act_name     string   `xml:"act_name"`
	Client_ip    string   `xml:"client_ip"`
	Mch_billno   string   `xml:"mch_billno"`
	Mch_id       string   `xml:"mch_id"`
	Nonce_str    string   `xml:"nonce_str"`
	Re_openid    string   `xml:"re_openid"`
	Remark       string   `xml:"remark"`
	Send_name    string   `xml:"send_name"`
	Total_amount int      `xml:"total_amount"`
	Total_num    int      `xml:"total_num"`
	Wishing      string   `xml:"wishing"`
	Wxappid      string   `xml:"wxappid"`
	Sign         string   `xml:"sign"`
}

const (
	WECHATCERTPATH    = "D:/apiclient_cert.pem"                                       //客户端证书存放绝对路径
	WECHATKEYPATH     = "D:/apiclient_key.pem"                                        //客户端私匙存放绝对路径
	WECHATCAPATH      = "D:/wwwroot/rootca.pem"                                       //服务端证书存放绝对路径
	WECHATURL         = "https://api.mch.weixin.qq.com/mmpaymkttransfers/sendredpack" //微信红包发送链接
	REDISURL          = "127.0.0.1:6379"                                              //redis地址
	ENDTIME           = "201603221657"                                                //红包活动结束时间
	OPENIDS           = "openids"                                                     //中奖用户redis openid列表名
	BONUSPRICEPOSTFIX = "_bonusprice"                                                 //redis get对应用户红包价格的后缀 例：set openid+BONUSPRICEPOSTFIX price

	//微信红包内容XML设置
	ACT_NAME  = ""           //活动名
	CLIENT_IP = ""           //服务器IP
	MCH_ID    = "1231911402" //商户号
	REMARK    = ""
	SEND_NAME = ""
	TOTAL_NUM = 1 //红包个数
	WISHING   = ""
	WXAPPID   = ""
	KEY       = ""
)

// var _tlsConfig *tls.Config

func getTLSConfig() (*tls.Config, error) {
	if _tlsConfig != nil {
		return _tlsConfig, nil
	}

	// load cert
	cert, err := tls.LoadX509KeyPair(WECHATCERTPATH, WECHATKEYPATH)
	if err != nil {
		fmt.Println("load wechat keys fail", err)
		return nil, err
	}

	// load root ca
	caData, err := ioutil.ReadFile(WECHATCAPATH)
	if err != nil {
		fmt.Println("read wechat ca fail", err)
		return nil, err
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caData)

	_tlsConfig = &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      pool,
	}
	return _tlsConfig, nil
}

func SecurePost(url string, xmlContent []byte) (*http.Response, error) {
	tlsConfig, err := getTLSConfig()
	if err != nil {
		return nil, err
	}

	tr := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: tr}

	return client.Post(
		url,
		"text/xml",
		bytes.NewBuffer(xmlContent))
}

func bonus(oid string, price int) {

	openid := oid               //openid
	total_amount := price * 100 //金额

	//生成随机数并md5加密随机数
	rand.Seed(time.Now().UnixNano())
	md5Ctx1 := md5.New()
	md5Ctx1.Write([]byte(strconv.Itoa(rand.Intn(1000))))
	nonce_str := hex.EncodeToString(md5Ctx1.Sum(nil)) //随机数
	//订单号
	mch_billno := MCH_ID + time.Now().Format("20060102") + strconv.FormatInt(time.Now().Unix(), 10)
	// 生成签名
	s1 := "act_name=" + ACT_NAME + "&client_ip=" + CLIENT_IP + "&mch_billno=" + mch_billno + "&mch_id=" + MCH_ID + "&nonce_str=" + nonce_str + "&re_openid=" + openid + "&remark=" + REMARK + "&send_name=" + SEND_NAME + "&total_amount=" + strconv.Itoa(total_amount) + "&total_num=" + strconv.Itoa(TOTAL_NUM) + "&wishing=" + WISHING + "&wxappid=" + WXAPPID + "&key=" + KEY

	md5Ctx2 := md5.New()
	md5Ctx2.Write([]byte(s1))
	s1 = hex.EncodeToString(md5Ctx2.Sum(nil))
	sign := strings.ToUpper(s1) //签名

	v := &Bonus{
		Act_name:     ACT_NAME,
		Client_ip:    CLIENT_IP,
		Mch_billno:   mch_billno,
		Mch_id:       MCH_ID,
		Nonce_str:    nonce_str,
		Remark:       REMARK,
		Send_name:    SEND_NAME,
		Re_openid:    openid,
		Total_amount: total_amount,
		Total_num:    TOTAL_NUM,
		Wishing:      WISHING,
		Wxappid:      WXAPPID,
		Sign:         sign}

	output, err := xml.MarshalIndent(v, "  ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	// os.Stdout.Write(output)
	//POST数据
	SecurePost(WECHATURL, output)

}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	ticker := time.NewTicker(3 * time.Second)

	for {
		if time.Now().Format("200601021504") >= ENDTIME {
			log.Printf("over")
			break
			ticker.Stop()
			os.Exit(1)
		} else {
			<-ticker.C
			//redis
			redisDataBase, _ := redis.Dial("tcp", REDISURL)
			openidsLen, _ := redis.Int(redisDataBase.Do("llen", OPENIDS))
			if openidsLen > 0 {
				log.Printf("Start Bonus")
				openid, _ := redis.String(redisDataBase.Do("lpop", OPENIDS))
				log.Printf(openid)
				//获取红包价格
				bonusPrice, _ := redis.Int(redisDataBase.Do("get", openid+BONUSPRICEPOSTFIX))
				redisDataBase.Do("del", openid+BONUSPRICEPOSTFIX)
				redisDataBase.Close()
				bonus(openid, bonusPrice)
				log.Printf("End Bonus")
			} else {
				log.Printf("Openids Is Empty")
				redisDataBase.Close()
			}
		}

	}

}
