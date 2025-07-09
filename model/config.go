package model

type Config struct {
	Key   string `bson:"key" json:"key"`
	Value string `bson:"value" json:"value"`
} 