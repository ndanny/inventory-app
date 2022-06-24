package models

import "math"

type Analytics struct {
	CompletedOrders int     `json:"completedOrders"`
	RejectedOrders  int     `json:"rejectedOrders"`
	TotalRevenue    float64 `json:"totalRevenue"`
}

// Combine adds the fields from two Analytics events and returns
// a new "merged" Analytics event
func Combine(this, that Analytics) Analytics {
	return Analytics{
		CompletedOrders: this.CompletedOrders + that.CompletedOrders,
		RejectedOrders:  this.RejectedOrders + that.RejectedOrders,
		TotalRevenue:    math.Round((this.TotalRevenue+that.TotalRevenue)*100) / 100,
	}
}
