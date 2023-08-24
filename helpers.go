package main

import (
	"encoding/json"
	"math"
	"net/http"
)

func calculateDistance(coord1, coord2 [2]float64) float64 {
	const earthRadiusKm = 6371

	lat1 := degToRad(coord1[0])
	lon1 := degToRad(coord1[1])
	lat2 := degToRad(coord2[0])
	lon2 := degToRad(coord2[1])

	deltaLat := lat2 - lat1
	deltaLon := lon2 - lon1

	a := math.Pow(math.Sin(deltaLat/2), 2) + math.Cos(lat1)*math.Cos(lat2)*math.Pow(math.Sin(deltaLon/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := earthRadiusKm * c
	return distance
}

func degToRad(deg float64) float64 {
	return deg * (math.Pi / 180)
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data any) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "json/application")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	return json.NewDecoder(r.Body).Decode(dst)
}

func getErrorFromChan(channel chan error) error {
	data := <-channel
	return data
}

func getItemOrError[T any](itemCh <-chan T, errCh <-chan error) (T, error) {
	var item T
	var err error
	select {
	case item = <-itemCh:
	case err = <-errCh:
	}
	return item, err
}
