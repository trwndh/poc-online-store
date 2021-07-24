package model

type Cart struct {
	ID         int64      `json:"cart_id"`
	UserID     int64      `json:"user_id"`
	Items      []CartItem `json:"items"`
	TotalPrice int64      `json:"total_price"`
}

type CartItem struct {
	ProductID   int64  `json:"product_id"`
	ProductName string `json:"product_name"`
	Quantity    int32  `json:"quantity"`
	Price       int32  `json:"price"`
	SubTotal    int32  `json:"sub_total"`
}

func (c Cart) IsUserIDNotValid() bool {
	return c.UserID == 0
}

func (c Cart) IsEmptyItems() bool {
	return len(c.Items) == 0
}

func (ci CartItem) IsProductIDValid() bool {
	return ci.ProductID > 0
}

func (ci CartItem) IsProductNameValid() bool {
	return ci.ProductName != ""
}

func (ci CartItem) IsQuantityValid() bool {
	return ci.Quantity > 0
}

/* checkout json:
{
	"data":{
		"cart_id":1,
		"user_id":2,
		"items":[
			product_id : 1,
			product_name: "Permen Minyak",
			quantity: 10,
			price: 2500,
			sub_total: 25000
		],
		"total_price":30000
	}
}
*/
