package MySQL

import (
	"gorm.io/gorm"
)

type ClientInfo struct {
	gorm.Model
	ClientID     int
	ClientSecret string
}

func (c ClientInfo) TableName() string {
	return "oauth_client_information"
}

func (c *ClientInfo) CheckClientID() bool {
	DB.Select("id").Where(c, "client_id").Find(c)
	return !(c.ID == 0)
}

func (c *ClientInfo) CheckClientInfo() bool {
	DB.Select("id").Where(c, "client_id", "client_secret").Find(c)
	return !(c.ID == 0)
}
