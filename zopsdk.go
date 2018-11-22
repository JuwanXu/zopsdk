package zopsdk

import (
	"sort"
	"encoding/json"
	"encoding/base64"
	"crypto/md5"
	"net/http"
	"strings"
	"io/ioutil"
)

type Client struct {
	CompanyId    string  `json:"company_id"`
	Key          string  `json:"key"`
}

type ZopRequest struct {
	Url			string 					`json:"url"`
	MsgType		string					`json:"msg_type"`
	ZopParams	map[string]interface{} 	`json:"zop_params"`
}


func (client *Client)Execute(zopRequest ZopRequest) (map[string]interface{},error) {
	requestParams:=zopRequest.ZopParams
	bodyParams:=map[string]string{}
	for k := range requestParams {
		value:=requestParams[k]
		valueJson,err:=json.Marshal(value)
		if err!=nil{
			panic(err)
		}
		valueStr:=string(valueJson[:])
		bodyParams[k]=valueStr
	}
	bodyParams=sortMapByKey(bodyParams)
	digest,err:=client.Sign(zopRequest.MsgType,bodyParams)
	if err!=nil{
		return nil,err
	}

	bodyJson,err:=json.Marshal(bodyParams)
	if err!=nil{
		panic(err)
	}
	resqBody:=string(bodyJson[:])
	httpClient:= &http.Client{}
	req, err := http.NewRequest("POST", zopRequest.Url, strings.NewReader(resqBody))
	if err != nil {
		return nil,err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("x-companyid", client.CompanyId)
	req.Header.Set("x-datadigest", digest)
	resp, err := httpClient.Do(req)
	defer resp.Body.Close()
	respBody, err:= ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil,err
	}
	resultMap:=map[string]interface{}{}
	err=json.Unmarshal(respBody,&resultMap)
	if err!=nil{
		return nil,err
	}
	return resultMap,nil

}

func (client *Client)Sign( msgType string, bodyParams  map[string]string)(string,error){
	bytesToDigest, err := json.Marshal(bodyParams)
	if err!=nil{
		return "",err
	}
	strToDigest := string(bytesToDigest[:])
	strToDigest="msg_type="+msgType+"&data="+strToDigest+"&company_id="+client.CompanyId+client.Key
	h := md5.New()
	h.Write([]byte(strToDigest))
	digest := h.Sum(nil)
	result:=base64.StdEncoding.EncodeToString(digest)
	return result,nil
}

func sortMapByKey(input map[string]string) map[string]string {
	var SortString []string
	output:=map[string]string{}
	for k := range input {
		SortString = append(SortString, k)
	}
	sort.Strings(SortString)  //会根据字母的顺序进行排序
	for _, k := range SortString {
		output[k]=input[k]
	}
	return  output
}

