package model

import "encoding/json"

type Crontab struct {
	Id        string `json:"_id"`
	Name      string `json:"name"`
	Command   string `json:"command"`
	Schedule  string `json:"schedule"`
	Timestamp string `json:"timestamp"`
	Stopped   bool   `json:"stopped"`
}

func (c *Crontab) Marshal() string {
	js, _ := json.Marshal(c)
	return string(js)
}
