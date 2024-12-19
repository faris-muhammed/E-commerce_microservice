// package models

// import "gorm.io/gorm"


// type ProductDetails struct {
// 	gorm.Model
// 	Id          uint    `gorm:"primarykey"`
// 	ProductName string  `gorm:"not null" json:"productname"`
// 	Price       float64 `json:"price"`
// 	Quantity    uint    `json:"quantity"`
// 	Size        string  `json:"size"`
// 	Brand       string  `json:"brand"`
// 	Barcode     string  `json:"barcode"`
// 	IsDeleted   bool    `gorm:"default:false"`
// 	SellerId    uint    `json:"sellerid"`
// 	Seller      SellerModel
// 	CategoryId  uint `json:"categoryid"`
// 	Category    Category
// }
package models

type Product struct {
	ID          string  `gorm:"primary_key"`
	Name        string
	Description string
	Price       float64
	Category    string
}