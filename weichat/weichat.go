package weichat

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"net/http"
	"sort"
	"strings"
)

type UnifyOrderReq struct {
	Appid            string `xml:"appid"`
	Body             string `xml:"body"`
	Mch_id           string `xml:"mch_id"`
	Nonce_str        string `xml:"nonce_str"`
	Notify_url       string `xml:"notify_url"`
	Trade_type       string `xml:"trade_type"`
	Spbill_create_ip string `xml:"spbill_create_ip"`
	Total_fee        string `xml:"total_fee"`
	Out_trade_no     string `xml:"out_trade_no"`
	OpenId           string `xml:"openid"`
	Sign             string `xml:"sign"`
}

type UnifyOrderResp struct {
	Return_code string `xml:"return_code"`
	Return_msg  string `xml:"return_msg"`
	Appid       string `xml:"appid"`
	Mch_id      string `xml:"mch_id"`
	Nonce_str   string `xml:"nonce_str"`
	Sign        string `xml:"sign"`
	Result_code string `xml:"result_code"`
	Prepay_id   string `xml:"prepay_id"`
	Trade_type  string `xml:"trade_type"`
}

func UnifiedOrder() {

	var yourReq UnifyOrderReq
	yourReq.Appid = "wx32e598477c7d1ef8"
	yourReq.Body = "1"
	yourReq.Mch_id = "1384831502" // 微信支付分配的商户号
	yourReq.Nonce_str = "5K8264ILTKCH16CQ2502SI8ZNMTM67VS"
	yourReq.Notify_url = "http://thribu.randiancx.com/weichat"
	yourReq.Trade_type = "JSAPI"
	yourReq.Spbill_create_ip = "183.14.251.198"
	yourReq.Total_fee = "1"
	yourReq.Out_trade_no = "20160905125346"
	yourReq.OpenId = "oDuXWv6YHNK34b-lhx5odUsLcyM4"

	var m map[string]interface{}
	m = make(map[string]interface{}, 0)
	m["appid"] = yourReq.Appid
	m["body"] = yourReq.Body
	m["mch_id"] = yourReq.Mch_id
	m["nonce_str"] = yourReq.Nonce_str
	m["notify_url"] = yourReq.Notify_url
	m["trade_type"] = yourReq.Trade_type
	m["spbill_create_ip"] = yourReq.Spbill_create_ip
	m["total_fee"] = yourReq.Total_fee
	m["out_trade_no"] = yourReq.Out_trade_no
	m["openid"] = yourReq.OpenId
	yourReq.Sign = wxpayCalcSign(m, "b840fc02d524045429941cc15f59e41c")

	bytes_req, err := xml.Marshal(yourReq)
	if err != nil {
		fmt.Println("以 XML 形式编码发送错误，原因：", err)
		return
	}

	str_req := string(bytes_req)
	str_req = strings.Replace(str_req, "UnifyOrderReq", "xml", -1)
	bytes_req = []byte(str_req)

	req, err := http.NewRequest("POST", "https://api.mch.weixin.qq.com/pay/unifiedorder", bytes.NewReader(bytes_req))
	if err != nil {
		fmt.Println("New Http Request 发生错误，原因：", err)
		return
	}
	req.Header.Set("Accpt", "application/xml")
	req.Header.Set("Content-type", "application/xml;charset=utf-8")

	c := http.Client{}
	resp, _err := c.Do(req)
	if _err != nil {
		fmt.Println("请求微信支付统一下单接口发生错误，原因：", _err)
		return
	}

	len := resp.ContentLength
	body := make([]byte, len)
	resp.Body.Read(body)
	defer resp.Body.Close()

	fmt.Println(string(body))

}

func HttpWeiChatPay(w http.ResponseWriter, r *http.Request) {
}
