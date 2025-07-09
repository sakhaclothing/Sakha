package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Tracker struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Hostname         string             `bson:"hostname" json:"hostname"`
	URL              string             `bson:"url" json:"url"`
	Browser          string             `bson:"browser" json:"browser"`
	BrowserLanguage  string             `bson:"browser_language" json:"browser_language"`
	ScreenResolution string             `bson:"screen_resolution" json:"screen_resolution"`
	Timezone         string             `bson:"timezone" json:"timezone"`
	Ontouchstart     bool               `bson:"ontouchstart" json:"ontouchstart"`
	TanggalAmbil     time.Time          `bson:"tanggal_ambil" json:"tanggal_ambil"`
	ISP              ISPInfo            `bson:"isp" json:"isp"`
}

type ISPInfo struct {
	IP        string  `bson:"ip" json:"ip"`
	City      string  `bson:"city" json:"city"`
	Region    string  `bson:"region" json:"region"`
	Country   string  `bson:"country_name" json:"country_name"`
	Postal    string  `bson:"postal" json:"postal"`
	Latitude  float64 `bson:"latitude" json:"latitude"`
	Longitude float64 `bson:"longitude" json:"longitude"`
	Timezone  string  `bson:"timezone" json:"timezone"`
	ASN       string  `bson:"asn" json:"asn"`
	Org       string  `bson:"org" json:"org"`
}
