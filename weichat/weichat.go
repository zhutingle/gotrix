package weichat

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
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

type WXPayNotifyReq struct {
	Return_code    string `xml:"return_code"`
	Return_msg     string `xml:"return_msg"`
	Appid          string `xml:"appid"`
	Mch_id         string `xml:"mch_id"`
	Nonce          string `xml:"nonce_str"`
	Sign           string `xml:"sign"`
	Result_code    string `xml:"result_code"`
	Openid         string `xml:"openid"`
	Is_subscribe   string `xml:"is_subscribe"`
	Trade_type     string `xml:"trade_type"`
	Bank_type      string `xml:"bank_type"`
	Total_fee      int    `xml:"total_fee"`
	Fee_type       string `xml:"fee_type"`
	Cash_fee       int    `xml:"cash_fee"`
	Cash_fee_Type  string `xml:"cash_fee_type"`
	Transaction_id string `xml:"transaction_id"`
	Out_trade_no   string `xml:"out_trade_no"`
	Attach         string `xml:"attach"`
	Time_end       string `xml:"time_end"`
}

func WxpayCallback(w http.ResponseWriter, r *http.Request) (returnMap map[string]interface{}, err error) {
	// body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("读取http body失败，原因!", err)
		http.Error(w.(http.ResponseWriter), http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return nil, err
	}
	defer r.Body.Close()

	fmt.Println("微信支付异步通知，HTTP Body:", string(body))
	var mr WXPayNotifyReq
	err = xml.Unmarshal(body, &mr)
	if err != nil {
		fmt.Println("解析HTTP Body格式到xml失败，原因!", err)
		http.Error(w.(http.ResponseWriter), http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return nil, err
	}

	var reqMap map[string]interface{}
	reqMap = make(map[string]interface{}, 0)

	reqMap["return_code"] = mr.Return_code
	reqMap["return_msg"] = mr.Return_msg
	reqMap["appid"] = mr.Appid
	reqMap["mch_id"] = mr.Mch_id
	reqMap["nonce_str"] = mr.Nonce
	reqMap["result_code"] = mr.Result_code
	reqMap["openid"] = mr.Openid
	reqMap["is_subscribe"] = mr.Is_subscribe
	reqMap["trade_type"] = mr.Trade_type
	reqMap["bank_type"] = mr.Bank_type
	reqMap["total_fee"] = mr.Total_fee
	reqMap["fee_type"] = mr.Fee_type
	reqMap["cash_fee"] = mr.Cash_fee
	reqMap["cash_fee_type"] = mr.Cash_fee_Type
	reqMap["transaction_id"] = mr.Transaction_id
	reqMap["out_trade_no"] = mr.Out_trade_no
	reqMap["attach"] = mr.Attach
	reqMap["time_end"] = mr.Time_end
	reqMap["sign"] = mr.Sign

	return reqMap, nil

	//	var resp WXPayNotifyResp
	//进行签名校验
	//	if wxpayVerifySign(reqMap, mr.Sign) {
	//这里就可以更新我们的后台数据库了，其他业务逻辑同理。
	//		resp.Return_code = "SUCCESS"
	//		resp.Return_msg = "OK"
	//	} else {
	//		resp.Return_code = "FAIL"
	//		resp.Return_msg = "failed to verify sign, please retry!"
	//	}

	//结果返回，微信要求如果成功需要返回return_code "SUCCESS"
	//	bytes, _err := xml.Marshal(resp)
	//	strResp := strings.Replace(string(bytes), "WXPayNotifyResp", "xml", -1)
	//	if _err != nil {
	//		fmt.Println("xml编码失败，原因：", _err)
	//		http.Error(w.(http.ResponseWriter), http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	//		return
	//	}
	//
	//	w.(http.ResponseWriter).WriteHeader(http.StatusOK)
	//	fmt.Fprint(w.(http.ResponseWriter), strResp)
}
