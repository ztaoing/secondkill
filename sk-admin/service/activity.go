/**
* @Author:zhoutao
* @Date:2020/7/7 上午10:34
 */

package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gohouse/gorose/v2"
	"github.com/samuel/go-zookeeper/zk"
	"github.com/unknwon/com"
	"log"
	pkgConfig "secondkill/pkg/config"
	"secondkill/sk-admin/model"
	"time"
)

type ActivityService interface {
	GetActivityList() ([]gorose.Data, error)
	CreateActivity(activity *model.Activity) error
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

//创建活动
func (a ActivityServiceImpl) CreateActivity(activity *model.Activity) error {
	//写入数据库
	activityEntity := model.NewActivityModel()
	err := activityEntity.CreateActivity(activity)
	if err != nil {
		log.Printf("CreateActivity failed,ERROR:%v", err)
		return err
	}
	//写入到zk或etcd
	log.Printf("sync to zk")
	err = a.syncToZk(activity)
	if err != nil {
		log.Printf("sync to zk failed,ERROR:%v", err)
		return err
	}
	return nil
}

//同步到zk中
//首先从zk中拉取存储的数据，如果数据为空，则将其转换为secProductInfoList,然后将新创建的activity添加到列表中，再更新到zk
//先调用conn的exist方法判断是否存在，如果存在则调用set方法更新，否则调用create方法创建新数据路径
func (a ActivityServiceImpl) syncToZk(activity *model.Activity) error {
	//数据路径
	zkPath := pkgConfig.Zk.SecProductKey
	//加载数据
	secProductInfoList, err := a.loadProductFromZk(zkPath)
	if err != nil {
		//空
		secProductInfoList = []*model.SecProductInfoConf{}
	}
	var secProductInfo = &model.SecProductInfoConf{}
	secProductInfo.EndTime = activity.Endtime
	secProductInfo.OnePersonBuyLimit = activity.BuyLimit
	secProductInfo.ProductId = activity.ProductId
	secProductInfo.SoldMaxLimit = activity.Speed
	secProductInfo.StartTime = activity.StartTime
	secProductInfo.Status = activity.Status
	secProductInfo.Total = activity.Total
	secProductInfo.BuyRate = activity.BuyRate
	secProductInfoList = append(secProductInfoList, secProductInfo)

	data, err := json.Marshal(secProductInfoList)
	if err != nil {
		log.Printf("json marshal failed,ERROR:%v", err)
		return err
	}

	conn := pkgConfig.Zk.ZkConn
	var bytedata = []byte(string(data))
	var flags int32 = 0
	//permission todo
	var acls = zk.WorldACL(zk.PermAll)

	exists, _, _ := conn.Exists(zkPath)
	if exists {
		//设置
		_, err_set := conn.Set(zkPath, bytedata, flags)
		if err_set != nil {
			fmt.Println(err_set)
		}
	} else {
		//创建
		_, err_create := conn.Create(zkPath, bytedata, flags, acls)
		if err_create != nil {
			fmt.Println(err_create)
		}
	}
	log.Println("save to zk success,data:%v", string(data))
	return nil

}

func (a ActivityServiceImpl) loadProductFromZk(zkpath string) ([]*model.SecProductInfoConf, error) {
	_, cancel := context.WithTimeout(context.Background(), time.Second*10)
	//超时取消
	defer cancel()
	v, stat, err := pkgConfig.Zk.ZkConn.Get(zkpath)
	if err != nil {
		log.Printf("get [%s] from zk failed,ERROR:%v", zkpath, err)
		return nil, err
	}

	log.Printf("get from zk success,resp:%v", stat)

	var secProductInfoList []*model.SecProductInfoConf
	err = json.Unmarshal(v, &secProductInfoList)
	if err != nil {
		log.Printf("unmarshal SecProductInfoConf failed,ERROR:%v", err)
		return nil, err
	}
	return secProductInfoList, nil
}

//将活动同步到etcd
func (a ActivityServiceImpl) syncToEtcd(activity *model.Activity) error {
	etcdKey := pkgConfig.Etcd.EtcdSecProductKey
	secProductInfoList, err := a.loadProductFromEtcd(etcdKey)
	if err != nil {
		//空
		secProductInfoList = []*model.SecProductInfoConf{}
	}
	var secProductInfo = &model.SecProductInfoConf{}
	secProductInfo.EndTime = activity.Endtime
	secProductInfo.OnePersonBuyLimit = activity.BuyLimit
	secProductInfo.ProductId = activity.ProductId
	secProductInfo.SoldMaxLimit = activity.Speed
	secProductInfo.StartTime = activity.StartTime
	secProductInfo.Status = activity.Status
	secProductInfo.Total = activity.Total
	secProductInfo.BuyRate = activity.BuyRate
	secProductInfoList = append(secProductInfoList, secProductInfo)
	data, err := json.Marshal(secProductInfoList)
	if err != nil {
		log.Printf("json marshal failed,ERROR:%v", err)
		return err
	}
	conn := pkgConfig.Etcd.EtcdConn
	_, err = conn.Put(context.Background(), etcdKey, string(data))
	if err != nil {
		log.Printf("put to etcd failed,ERROR:%v", err)
		return err
	}
	return nil
}

func (a ActivityServiceImpl) loadProductFromEtcd(etcdKey string) ([]*model.SecProductInfoConf, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := pkgConfig.Etcd.EtcdConn.Get(ctx, etcdKey)
	if err != nil {
		log.Printf("get [%s] from etcd failed ,ERROR:%v", etcdKey, err)
		return nil, err
	}

	log.Printf("get [%s] from etcd success,RESP:%v", etcdKey, resp)

	var secProductInfoList []*model.SecProductInfoConf
	for _, v := range resp.Kvs {
		err := json.Unmarshal(v.Value, &secProductInfoList)
		if err != nil {
			log.Printf("unmarshal product info failed,ERROR:%v", err)
			return nil, err
		}
	}
	return secProductInfoList, nil
}

//装饰
type ActivityServiceMiddleware func(ActivityService) ActivityService
