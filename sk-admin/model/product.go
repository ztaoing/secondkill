/**
* @Author:zhoutao
* @Date:2020/7/7 上午9:30
 */

package model

import (
	"github.com/gohouse/gorose/v2"
	"log"
	"secondkill/pkg/mysql"
)

type Product struct {
	ProductId   int    `json:"product_id"`   //商品id
	ProductName string `json:"product_name"` //商品名称
	Total       int    `json:"total"`        //商品数量
	Status      int    `json:"status"`       //商品状态
}

type ProductModel struct {
}

func NewProductModel() *ProductModel {
	return &ProductModel{}
}

func (p *ProductModel) getTabName() string {
	return "product"
}

//获得商品列表
func (p *ProductModel) GetProductList() ([]gorose.Data, error) {
	conn := mysql.DB()
	list, err := conn.Table(p.getTabName()).Get()
	if err != nil {
		log.Printf("ERROR:%v", err)
		return nil, err
	}
	return list, nil
}

//创建商品
func (p *ProductModel) CreateProduct(product *Product) error {
	conn := mysql.DB()
	_, err := conn.Table(p.getTabName()).Data(
		map[string]interface{}{
			"product_name": product.ProductName,
			"total":        product.Total,
			"stauts":       product.Status,
		},
	).Insert()
	if err != nil {
		log.Printf("ERROR:%v", err)
		return err
	}
	return nil
}
