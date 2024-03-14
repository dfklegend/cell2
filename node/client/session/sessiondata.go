package session

import (
	"encoding/json"
)

// SessionData
// 方便打包传到后端服务器构建backsession
// 注意，整数传过来变成了浮点数(json的原因)
type SessionData struct {
	data map[string]interface{}
}

func NewSessionData() *SessionData {
	return &SessionData{
		data: make(map[string]interface{}),
	}
}

func (d *SessionData) Reset() {
	if len(d.data) == 0 {
		return
	}
	d.data = make(map[string]interface{})
}

func (d *SessionData) GetMap() map[string]interface{} {
	return d.data
}

func (d *SessionData) Set(k string, v interface{}) {
	d.data[k] = v
}

func (d *SessionData) Get(k string, def interface{}) interface{} {
	v, ok := d.data[k]
	if ok {
		return v
	}
	return def
}

func (d *SessionData) Has(k string) bool {
	_, ok := d.data[k]
	return ok
}

func (d *SessionData) ToJson() []byte {
	v, _ := json.Marshal(d.data)
	return v
}

func (d *SessionData) FromJson(data []byte) {
	json.Unmarshal(data, &d.data)
}

func (d *SessionData) ToJsonStr() string {
	return string(d.ToJson())
}

func (d *SessionData) FromJsonStr(s string) {
	json.Unmarshal([]byte(s), &d.data)
}

func (d *SessionData) UpdateFromJson(data []byte) {
	newData := make(map[string]interface{})
	json.Unmarshal(data, &newData)

	// 更新过去
	for k, v := range newData {
		d.data[k] = v
	}
}
