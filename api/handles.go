package api

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"advertising-system/db"
	"advertising-system/models"

	"github.com/gorilla/mux"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

type AdResponse struct {
	AdvertisementID int64 `json:"advertisement_id"`
	Adm             struct {
		Title         string `json:"title"`
		Content       string `json:"content"`
		TrackingLinks struct {
			Impression string `json:"impression"`
			Click      string `json:"click"`
		} `json:"tracking_links"`
	} `json:"adm"`
}

type Response struct {
	Message string `json:"message"`
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, Response{Message: message})
}

func GetAdvertiser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid advertiser ID")
		return
	}

	var advertiser models.Advertiser
	if err := db.DB.First(&advertiser, id).Error; err != nil {
		writeError(w, http.StatusNotFound, "Advertiser not found")
		return
	}

	writeJSON(w, http.StatusOK, advertiser)
}

func GetAdvertisement(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid advertisement ID")
		return
	}

	var ad models.Advertisement
	if err := db.DB.First(&ad, id).Error; err != nil {
		writeError(w, http.StatusNotFound, "Advertisement not found")
		return
	}

	writeJSON(w, http.StatusOK, ad)
}

func RecordImpression(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	adID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid advertisement ID")
		return
	}

	key := fmt.Sprintf("ad:%d:counts", adID)
	field := "cpm_count"

	// use HIncrBy to increment the impression count
	if err := db.Rdb.HIncrBy(context.Background(), key, field, 1).Err(); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to record impression in Redis")
		return
	}

	writeJSON(w, http.StatusOK, Response{Message: "Impression recorded successfully"})
}

func RecordClick(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	adID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid advertisement ID")
		return
	}

	key := fmt.Sprintf("ad:%d:counts", adID)
	field := "cpc_count"
	// use HIncrBy to increment the click count
	if err := db.Rdb.HIncrBy(context.Background(), key, field, 1).Err(); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to record click in Redis")
		return
	}

	writeJSON(w, http.StatusOK, Response{Message: "Click recorded successfully"})
}

func RequestAd(w http.ResponseWriter, r *http.Request) {
	var ads []models.Advertisement

	if err := db.DB.Joins("JOIN advertisers ON advertisements.advertiser_id = advertisers.id").
		Where("advertisers.budget > 0").
		Find(&ads).Error; err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to fetch advertisements")
		return
	}

	if len(ads) == 0 {
		writeError(w, http.StatusNotFound, "No advertisements available")
		return
	}

	selectedAd := ads[rng.Intn(len(ads))]

	var fullAd models.Advertisement
	if err := db.DB.First(&fullAd, selectedAd.ID).Error; err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to fetch advertisement details")
		return
	}

	response := AdResponse{}
	response.AdvertisementID = fullAd.ID
	response.Adm.Title = fullAd.Title
	response.Adm.Content = fullAd.Content
	response.Adm.TrackingLinks.Impression = fmt.Sprintf("/api/advertisements/%d/impression", fullAd.ID)
	response.Adm.TrackingLinks.Click = fmt.Sprintf("/api/advertisements/%d/click", fullAd.ID)

	writeJSON(w, http.StatusOK, response)
}
