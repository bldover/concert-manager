package util

import (
	"concert-manager/data"
	"sort"
)

const (
	ExactTolerance    = 0.0
	StrictTolerance   = 0.1
	ModerateTolerance = 0.25
	LenientTolerance  = 0.4
)

const NoMaxResults = 0

type optionDistance[T any] struct {
	Option   T
	Distance int
}

func SearchOptions[T any](term string, options []T, maxResults int, tolerance float64, computeDistance func(string, T) int) []T {
	threshold := int(float64(len(term)) * tolerance)
	var optionDistances []optionDistance[T]

	for _, option := range options {
		distance := computeDistance(term, option)
		if distance <= threshold {
			optionDistances = append(optionDistances, optionDistance[T]{Option: option, Distance: distance})
		}
	}

	sort.Slice(optionDistances, func(i, j int) bool {
		return optionDistances[i].Distance < optionDistances[j].Distance
	})

	var results []T
	for i := 0; i < len(optionDistances) && (maxResults == NoMaxResults || i < maxResults); i++ {
		results = append(results, optionDistances[i].Option)
	}

	return results
}

func SearchArtists(term string, options []data.Artist, maxResults int, tolerance float64) []data.Artist {
	return SearchOptions(term, options, maxResults, tolerance, computeArtistDistance)
}

func SearchVenues(term string, options []data.Venue, maxResults int, tolerance float64) []data.Venue {
	return SearchOptions(term, options, maxResults, tolerance, computeVenueDistance)
}

func SearchEventsByArtists(term string, options []data.Event, maxResults int, tolerance float64) []data.Event {
	return SearchOptions(term, options, maxResults, tolerance, computeEventDistanceByArtists)
}

func SearchEventsByVenue(term string, options []data.Event, maxResults int, tolerance float64) []data.Event {
	return SearchOptions(term, options, maxResults, tolerance, func(term string, option data.Event) int {
		return computeVenueDistance(term, option.Venue)
	})
}

func SearchEventDetailsByArtist(term string, options []data.EventDetails, maxResults int, tolerance float64) []data.EventDetails {
	return SearchOptions(term, options, maxResults, tolerance, func(term string, option data.EventDetails) int {
		return computeEventDistanceByArtists(term, option.Event)
	})
}

func SearchEventDetailsByVenue(term string, options []data.EventDetails, maxResults int, tolerance float64) []data.EventDetails {
	return SearchOptions(term, options, maxResults, tolerance, func(term string, option data.EventDetails) int {
		return computeVenueDistance(term, option.Event.Venue)
	})
}

func SearchStrings(term string, options []string, maxResults int, tolerance float64) []string {
	return SearchOptions(term, options, maxResults, tolerance, getLevenshteinDistance)
}

func computeVenueDistance(term string, option data.Venue) int {
	return getLevenshteinDistance(term, option.Name)
}

func computeArtistDistance(term string, option data.Artist) int {
	return getLevenshteinDistance(term, option.Name)
}

func computeEventDistanceByArtists(term string, option data.Event) int {
	minDistance := getLevenshteinDistance(term, option.MainAct.Name)
	for _, opener := range option.Openers {
		d := getLevenshteinDistance(term, opener.Name)
		if d < minDistance {
			minDistance = d
		}
	}
	return minDistance
}
