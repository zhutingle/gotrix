package weichat

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"net/http"
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

func UnifiedOrder(orderReq UnifyOrderReq) (returnMap map[string]interface{}) {

	bytes_req, err := xml.Marshal(orderReq)
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

	var orderResp UnifyOrderResp
	xml.Unmarshal(body, &orderResp)

	returnMap = make(map[string]interface{}, 0)
	returnMap["return_code"] = orderResp.Return_code
	returnMap["return_msg"] = orderResp.Return_msg
	returnMap["appid"] = orderResp.Appid
	returnMap["mch_id"] = orderResp.Mch_id
	returnMap["nonce_str"] = orderResp.Nonce_str
	returnMap["sign"] = orderResp.Sign
	returnMap["result_code"] = orderResp.Result_code
	returnMap["prepay_id"] = orderResp.Prepay_id
	returnMap["trade_type"] = orderResp.Trade_type

	return returnMap

}
