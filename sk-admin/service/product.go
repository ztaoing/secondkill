/**
* @Author:zhoutao
* @Date:2020/7/7 下午4:13
 */

package service

import (
	"github.com/gohouse/gorose/v2"
	"log"
	"secondkill/sk-admin/model"
)

type ProductService interface {
	CreateProduct(product *model.Product) error
	GetProductList() ([]gorose.Data, error)
}

type ProductServiceImpl struct {
}

func (p ProductServiceImpl) CreateProduct(product *model.Product) error {
	productEntity := model.NewProductModel()
	err := productEntity.CreateProduct(product)
	if err != nil {
		log.Printf("productEntity CreateProduct,ERROR:%v", err)
		return err
	}
	return nil
}

func (p ProductServiceImpl) GetProductList() ([]gorose.Data, error) {
	productEntity := model.NewProductModel()
	productList, err := productEntity.GetProductList()
	if err != nil {
		log.Printf("productEntity GetProductList,ERROR:%v", err)
		return nil, err
	}
	return productList, nil
}

type ProductServiceMiddleware func(ProductService) ProductService
