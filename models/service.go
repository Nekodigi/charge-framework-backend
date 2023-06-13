package models

import "time"

type Plan struct {
	PriceId   string  `json:"priceId" firestore:"priceId"`
	Quota     float64 `json:"quota" firestore:"quota"`
	QuotaLeak float64 `json:"quotaLeak" firestore:"quotaLeak"`
}

type Service struct {
	Id          string          `json:"id" firestore:"id"`
	AllocQuota  float64         `json:"allocQuota" firestore:"allocQuota"`
	RemainQuota float64         `json:"remainQuota" firestore:"remainQuota"`
	Plan        map[string]Plan `json:"plan" firestore:"plan"`
	UpdateAt    time.Time       `json:"updateAt" firestore:"updateAt"`
}
