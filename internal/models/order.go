package models

type OrderStatus string

const (
	Pending   OrderStatus = "pending"
	Confirmed OrderStatus = "confirmed"
	Accepted  OrderStatus = "accepted"
	PickedUp  OrderStatus = "pickedUp"
	Delivered OrderStatus = "delivered"
	Cancelled OrderStatus = "cancelled"
	Deleted   OrderStatus = "deleted"
)

func (r OrderStatus) String() string {
	return string(r)
}

type DeliveryOption string

const (
	Express  DeliveryOption = "express"
	Today    DeliveryOption = "today"
	Tomorrow DeliveryOption = "tomorrow"
	Later    DeliveryOption = "later"
)

func (r DeliveryOption) String() string {
	return string(r)
}
