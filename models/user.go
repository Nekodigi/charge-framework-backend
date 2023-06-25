package models

import "time"

type User struct {
	Id           string    `json:"id" firestore:"id"`
	CustomerId   string    `json:"customerId" firestore:"customerId"`
	ServiceId    string    `json:"serviceId" firestore:"serviceId"`
	Plan         string    `json:"plan" firestore:"plan"`
	Status       string    `json:"status" firestore:"status"`
	AllocQuota   float64   `json:"allocQuota" firestore:"allocQuota"`
	RemainQuota  float64   `json:"remainQuota" firestore:"remainQuota"`
	Subscription string    `json:"subscription" firestore:"subscription"`
	Redirect     string    `json:"redirect" firestore:"redirect"`
	UpdateAt     time.Time `json:"updateAt" firestore:"updateAt"`
}
