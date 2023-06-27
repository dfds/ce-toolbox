package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
)

type Stats struct {
	TotalSignins             int
	UniqueUsers              int
	AverageUserSignIn        int
	MeanUserSignIn           float64
	SignInCountByUser        map[string]int
	SignInCountryTotalCount  map[string]int
	SignInCountryUniqueUsers map[string]map[string]int
}

func NewStats() Stats {
	return Stats{
		SignInCountByUser:        map[string]int{},
		SignInCountryTotalCount:  map[string]int{},
		SignInCountryUniqueUsers: map[string]map[string]int{},
	}
}

func (s *Stats) calcMedianSignInCount() {
	data := []float64{}
	for _, entry := range s.SignInCountByUser {
		data = append(data, float64(entry))
	}

	s.MeanUserSignIn = CalcMedian(data)
}

func (s *Stats) OutputTotalAsCsv() string {
	buf := bytes.Buffer{}
	fmt.Fprintf(&buf, "totalSignins,uniqueUsers,averageUserSignIn,MeanUserSignIn\n")
	fmt.Fprintf(&buf, "%d,%d,%d,%.2f", s.TotalSignins, s.UniqueUsers, s.AverageUserSignIn, s.MeanUserSignIn)

	return buf.String()
}

func (s *Stats) OutputUserActivityAsCsv() string {
	buf := bytes.Buffer{}
	fmt.Fprintf(&buf, "user,totalSignIns\n")
	for user, count := range s.SignInCountByUser {
		fmt.Fprintf(&buf, "%s,%d\n", user, count)
	}

	return buf.String()
}

func (s *Stats) OutputUserCountryAsCsv() string {
	buf := bytes.Buffer{}
	fmt.Fprintf(&buf, "country,totalSignIns,uniqueUsers\n")

	for country, count := range s.SignInCountryTotalCount {
		fmt.Fprintf(&buf, "%s,%d,%d\n", country, count, len(s.SignInCountryUniqueUsers[country]))
	}

	return buf.String()
}

func (s *Stats) CalcGeneralStats(data InputData) {
	// Gather stats
	for _, entry := range data {
		s.TotalSignins = s.TotalSignins + 1
		s.SignInCountByUser[entry.UserPrincipalName] = s.SignInCountByUser[entry.UserPrincipalName] + 1
	}

	// Calc stats
	s.UniqueUsers = len(s.SignInCountByUser)
	s.AverageUserSignIn = s.TotalSignins / s.UniqueUsers
	s.calcMedianSignInCount()
}

func (s *Stats) CalcCountries(data InputData) {
	for _, entry := range data {
		s.SignInCountryTotalCount[entry.Location.CountryOrRegion] = s.SignInCountryTotalCount[entry.Location.CountryOrRegion] + 1

		if _, ok := s.SignInCountryUniqueUsers[entry.Location.CountryOrRegion]; !ok {
			s.SignInCountryUniqueUsers[entry.Location.CountryOrRegion] = map[string]int{}
		}

		s.SignInCountryUniqueUsers[entry.Location.CountryOrRegion][entry.UserPrincipalName] = 0
	}
}

func CalcMedian(n []float64) float64 {
	sort.Float64s(n) // sort the numbers

	mNumber := len(n) / 2

	if IsOdd(n) {
		return n[mNumber]
	}

	return (n[mNumber-1] + n[mNumber]) / 2
}

func IsOdd(n []float64) bool {
	if len(n)%2 == 0 {
		return false
	}

	return true
}

func LoadInputData(path string) InputData {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Unable to load data from %s\n", path)
	}

	var payload InputData

	err = json.Unmarshal(data, &payload)
	if err != nil {
		log.Println("Unable to deserialise json from input data")
		log.Fatal(err)
	}

	return payload
}
