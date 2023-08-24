package main

type Shop struct {
	ID                    string     `bson:"_id,omitempty"`
	Name                  string     `bson:"name"`
	Location              string     `bson:"location"`
	OperationHours        string     `bson:"operation_hours"`
	ServiceableProductsId []string   `bson:"products"`
	Coordinates           [2]float64 `bson:"coordinates"`
}

func (s Shop) Id() string {
	return s.ID
}

type Product struct {
	ID          string  `bson:"_id,omitempty"`
	Name        string  `bson:"name"`
	Description string  `bson:"description"`
	Price       float64 `bson:"price"`
}

func (s Product) Id() string {
	return s.ID
}

type Inventory struct {
	ID        string `bson:"_id,omitempty"`
	ShopID    string `bson:"shop_id"`
	ProductID string `bson:"product_id"`
	Quantity  int    `bson:"quantity"`
}

func (s Inventory) Id() string {
	return s.ID
}

type ServiceableProduct struct {
	ID        string `bson:"_id,omitempty"`
	ProductID string `bson:"product_id"`
	ShopID    string `bson:"shop_id"`
}

func (s ServiceableProduct) Id() string {
	return s.ID
}

type ShopsServiceableProducts struct {
	ID       string    `bson:"_id,omitempty"`
	ShopID   string    `bson:"shop_id"`
	Products []Product `bson:"products"`
}

func (s ShopsServiceableProducts) Id() string {
	return s.ID
}

type User struct {
	ID          string     `bson:"_id,omitempty"`
	Name        string     `bson:"name"`
	Location    string     `bson:"location"`
	Coordinates [2]float64 `bson:"coordinates"`
}

func (s User) Id() string {
	return s.ID
}
