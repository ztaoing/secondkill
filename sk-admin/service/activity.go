/**
* @Author:zhoutao
* @Date:2020/7/7 上午10:34
 */

package service

import (
	"fmt"
	"github.com/gohouse/gorose/v2"
	"github.com/unknwon/com"
	"log"
	"secondkill/sk-admin/model"
	"time"
)

type ActivityService interface {
	GetActivityList() ([]gorose.Data, error)
	CreateActivity(activity *model.Activity) bool
}

type ActivityServiceImpl struct {
}

//从mysql数据库中获取活动列表
func (a ActivityServiceImpl) GetActivityList() ([]gorose.Data, error) {
	activity := model.NewActivityModel()
	activityList, err := activity.GetActivityList()

	if err != nil {
		log.Printf("activity GetActivityList failed,ERROR:%v", err)
		return nil, err
	}

	for _, v := range activityList {
		startTime, _ := com.StrTo(fmt.Sprint(v["start_time"])).Int64()
		v["start_time_str"] = time.Unix(startTime, 0).Format("2006-01-02 15:04:05")

		endTime, _ := com.StrTo(fmt.Sprint(v["end_time"])).Int64()
		v["end_time_str"] = time.Unix(endTime, 0).Format("2006-01-02 15:04:05")

		nowTime := time.Now().Unix()
		if nowTime > endTime {
			v["status_str"] = "已结束"
			continue
		}

		status, _ := com.StrTo(fmt.Sprint(v["status"])).Int()
		if status == model.ActivityStatusNormal {
			v["status_str"] = "正常"
		} else if status == model.ActivityStatusDisable {
			v["status_str"] = "已禁用"
		}
	}
	log.Printf("get activity success")
	return activityList, nil
}
