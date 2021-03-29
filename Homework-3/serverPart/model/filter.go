package model

import (
	"log"
	"serverPart/dao"
	ut "serverPart/utilities"
)

func FiltersUpdate(f *ut.Filters) {
	err := dao.DB.Select("id", "content").Find(f).Error
	if err != nil {
		log.Println(err)
	}
}

func BannersUpdate(b *map[string][]string) {
	var (
		name    []string
		banName []string
	)
	err := dao.DB.Table("ban_information").Distinct("id").Select("users_information.name").
		Joins("left join users_information on users_information.id = ban_information.user_id").Find(&name)
	if err != nil {
		//log.Println("BannersUpdate Error1:", err)
	}
	for _, v := range name {
		err = dao.DB.Table("ban_information").Distinct("ban_user_id").Select("users_information.name").
			Joins("left join users_information on users_information.id = ban_information.ban_user_id").Find(&banName)
		if err != nil {
			//log.Println("BannersUpdate Error2:", err)
		}
		(*b)[v] = banName
	}
}
