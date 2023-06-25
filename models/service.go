package models

import "time"

type Plan struct {
	Id        string  `json:"id" firestore:"id"`
	PriceId   string  `json:"priceId" firestore:"priceId"`
	Currency  string  `json:"currency" firestore:"currency"`
	Price     float64 `json:"price" firestore:"price"`
	Mode      string  `json:"mode" firestore:"mode"`
	Quota     float64 `json:"quota" firestore:"quota"`
	QuotaLeak float64 `json:"quotaLeak" firestore:"quotaLeak"`
}

type Service struct {
	Id          string          `json:"id" firestore:"id"`
	Name        string          `json:"name" firestore:"name"`
	AllocQuota  float64         `json:"allocQuota" firestore:"allocQuota"`
	RemainQuota float64         `json:"remainQuota" firestore:"remainQuota"`
	Plan        map[string]Plan `json:"plan" firestore:"plan"`
	UpdateAt    time.Time       `json:"updateAt" firestore:"updateAt"`
}
