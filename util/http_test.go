package util

import (
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"testing"
)

func TestHttpPost(t *testing.T) {
	HttpPost("http://127.0.0.1", "test")
}

type ArcAigcParam struct {
	TaskType            int32  `json:"algo_id"`
	TaskStyle           int32  `json:"algo_style"`
	ImageUrl            string `json:"image"`
	VideoUrl            string `json:"video"`
	MaskUrl             string `json:"mask"`
	Text                string `json:"text"`
	Title               string `json:"title"`
	Tags                string `json:"tags"`
	MakeInstrumental    bool   `json:"makeInstrumental"`
	SynthesizeVoiceName string `json:"synthesizeVoiceName"`
	Ssml                string `json:"ssml"`
	// params 3dVideo {"arcAigc":{"algo_id":8,"algo_style":101,"video":"/mnt/data/remote/2450MAM/bc/video/bag.mp4","params":{"projectType":0,"modelType":0,"useGS":true,"useBgRemove":true}}}
	Params interface{} `json:"params"`
}
type extractArcAigcParam struct {
	ArcAigc ArcAigcParam `json:"arcAigc"`
}

func DecodeArcAigcParam(extraData string) (*ArcAigcParam, error) {
	a := extractArcAigcParam{}
	err := json.Unmarshal([]byte(extraData), &a)
	if err != nil {
		return nil, err
	}
	return &a.ArcAigc, nil
}
func TestUrl(t *testing.T) {
	//str := "http://unknowhost.com/person.php"
	//str, _ = HttpAddParam(str, "taskId", "34")
	//t.Log(str)
	//str = "http://unknowhost.com/person.php?taskId=434&dafs=3"
	//str, _ = HttpAddParam(str, "taskId", "34")
	//t.Log(str)

	extraData := "{\"arcAigc\":{\"algo_id\":8,\"algo_style\":101,\"video\":\"/mnt/data/remote/2450MAM/bc/video/bag.mp4\",\"params\":{\"projectType\":0,\"modelType\":0,\"useGS\":true,\"useBgRemove\":true}}}"
	param, err := DecodeArcAigcParam(extraData)
	if err != nil {
		return
	}
	jParams, _ := convertSimpleJsonParam(param.Params)

	useGS, _ := jParams.CheckGet("useGS")
	if v, ok := useGS.Bool(); ok == nil {
		fmt.Println(v)
	}
}

func convertSimpleJsonParam(params interface{}) (*simplejson.Json, error) {
	m, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	p, err := simplejson.NewJson(m)
	if err != nil {
		return nil, err
	}
	return p, nil
}
