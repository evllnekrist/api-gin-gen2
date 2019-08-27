package helpers

import(
	"encoding/json"
	"time"
)

type EveBasicHelper struct{}

var start time.Time

func (hlpr EveBasicHelper) Panics(err error){ //error handling
	if err != nil{
		panic(err.Error())
	}
}

func (hlpr EveBasicHelper) GetCurrentTime() string{
	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)
	return now.String()
}

//Create the desc & merge with original json. Code adapted from --> https://www.sohamkamani.com/blog/2017/10/18/parsing-json-in-golang/
func (hlpr EveBasicHelper) GenerateJsonDesc(oldJson string, status bool) ([]byte, error){
	type gen struct{
		ListData string `json:"listdata"`
		GenTime string `json:"generatetime"`
		GenStatus bool `json:"generatestatus"`
	}
	desc_value	:= &gen{GenTime: hlpr.GetCurrentTime(), GenStatus: status}
	temp, newJson := map[string]interface{}{}, map[string]interface{}{}
	json.Unmarshal([]byte(oldJson), &temp)
	newJson["listdata"] = temp
	newJson["generatetime"] = desc_value.GenTime
	newJson["generatestatus"] = desc_value.GenStatus
	return json.Marshal(newJson) 
}
