/**
* @Author:zhoutao
* @Date:2020/7/7 上午8:46
 */

package model

import (
	"github.com/gohouse/gorose/v2"
	"log"
	"secondkill/pkg/mysql"
)

//活动
type Activity struct {
	ActivityId   int    `json:"activity_id"`   //活动id
	ActivityName string `json:"activity_name"` //活动名称
	ProductId    int    `json:"product_id"`    //商品id
	StartTime    int64  `json:"start_time"`    //开始时间
	Endtime      int64  `json:"endtime"`       //结束时间
	Total        int    `json:"total"`         //商品总数
	Status       int    `json:"status"`        //状态

	StartTimeStr string  `json:"start_time_str"`
	EndTimeStr   string  `json:"end_time_str"`
	StatusStr    string  `json:"status_str"`
	Speed        int     `json:"speed"`
	BuyLimit     int     `json:"buy_limit"`
	BuyRate      float64 `json:"buy_rate"`
}

//秒杀商品
type SecProductInfoConf struct {
	ProductId         int     `json:"product_id"`           //商品id
	StartTime         int64   `json:"start_time"`           //开始时间
	EndTime           int64   `json:"end_time"`             //结束时间
	Status            int     `json:"status"`               //状态
	Total             int     `json:"total"`                //商品总数
	Left              int     `json:"left"`                 //剩余商品数
	OnePersonBuyLimit int     `json:"one_person_buy_limit"` //单用户购买限制
	BuyRate           float64 `json:"buy_rate"`             //买中几率
	SoldMaxLimit      int     `json:"sold_max_limit"`       //每秒最多能卖多少
}

type ActivityModel struct {
}

func NewActivityModel() *ActivityModel {
	return &ActivityModel{}
}

func (p *ActivityModel) getTableName() string {
	return "activity"
}

//获取活动列表
func (p *ActivityModel) GetActivityList() ([]gorose.Data, error) {
	conn := mysql.DB()
	list, err := conn.Table(p.getTableName()).Order("activity_id desc").Get()
	if err != nil {
		log.Printf("ERROR:%v", err)
		return nil, err
	}
	return list, nil
}

//创建活动
func (p *ActivityModel) CreateActivity(activity *Activity) error {
	conn := mysql.DB()
	_, err := conn.Table(p.getTableName()).Data(
		map[string]interface{}{
			"activity_name": activity.ActivityName,
			"product_id":    activity.ProductId,
			"start_time":    activity.StartTime,
			"end_time":      activity.Endtime,
			"total":         activity.Total,
			"sec_speed":     activity.Speed,
			"buy_limit":     activity.BuyLimit,
			"buy_rate":      activity.BuyRate,
		},
	).Insert()
	if err != nil {
		return err
	}
	return nil
}
