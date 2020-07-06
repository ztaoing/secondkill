/**
* @Author:zhoutao
* @Date:2020/7/6 下午4:40
 */

package svc_user

import "sync"

//用户购买记录
type UserBuyHistory struct {
	History map[int]int
	Lock    sync.RWMutex
}

//读取商品的购买数量
func (u *UserBuyHistory) GetProductBuyCount(productId int) int {
	u.Lock.RLock()
	defer u.Lock.RUnlock()
	count, _ := u.History[productId]
	return count
}

//增加
func (u *UserBuyHistory) Add(productId, count int) {
	u.Lock.Lock()
	defer u.Lock.Unlock()

	cur, ok := u.History[productId]
	if !ok {
		cur = count
	} else {
		cur += count
	}
	u.History[productId] = cur

}
