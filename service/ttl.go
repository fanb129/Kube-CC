package service

import (
	"Kube-CC/dao"
	"errors"
	"gorm.io/gorm"
	"time"
)

func CreateOrUpdateTtl(ns string, expiredTime time.Time) error {
	ttl, _ := dao.GetTtlByNs(ns)
	if ttl != nil {
		ttl.ExpiredTime = expiredTime
		row, err := dao.UpdateTtl(ttl)
		if err != nil || row == 0 {
			return errors.New("ttl更新失败")
		}
		return nil
	}
	// 是否被软删除
	ttl, _ = dao.GetDeletedTtlByNs(ns)
	if ttl != nil {
		ttl.ExpiredTime = expiredTime
		ttl.DeletedAt = gorm.DeletedAt{}
		ttl.CreatedAt = time.Now()
		row, err := dao.UpdateUnscopedTtl(ttl)
		if err != nil || row == 0 {
			return errors.New("ttl添加失败")
		}
	} else {
		row, err := dao.CreateTtl(ns, expiredTime)
		if err != nil || row == 0 {
			return errors.New("ttl添加失败")
		}
	}
	return nil
}

func DeleteTtl(ns string) error {
	ttl, err := dao.GetTtlByNs(ns)
	if ttl != nil {
		err = dao.DeleteTtl(ttl)
		if err != nil {
			return errors.New("删除ttl失败")
		}
	}
	return nil
}

//func updateTtl(ns string, expiredTime time.Time) error {
//	ttl, err := dao.GetTtlByNs(ns)
//	if err != nil {
//		return errors.New("获取ttl失败")
//	}
//	ttl.ExpiredTime = expiredTime
//	row, err := dao.UpdateTtl(ttl)
//	if err != nil || row == 0 {
//		return errors.New("ttl更新失败")
//	}
//	return nil
//}
