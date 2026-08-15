package adapter

import (
	"context"
	"fmt"
	"sync"
)

// MultiSearchQuery specifies criteria for querying multiple provider categories in parallel.
type MultiSearchQuery struct {
	HospitalCriteria  *HospitalSearchCriteria
	FerryCriteria     []FerrySearchCriteria
	HotelCriteria     *HotelSearchCriteria
	TransportCriteria []TransportSearchCriteria
}

// MultiSearchResult collects offers across provider categories along with non-fatal warnings.
type MultiSearchResult struct {
	HospitalOffers  []Offer
	FerryOffers     []Offer
	HotelOffers     []Offer
	TransportOffers []Offer
	Warnings        []string
}

// Aggregator coordinates parallel provider searches with error isolation.
type Aggregator struct {
	hospital  *HospitalAdapter
	ferry     *FerryAdapter
	hotel     *HotelAdapter
	transport *TransportAdapter
}

// NewAggregator constructs a new multi-provider search Aggregator.
func NewAggregator(
	hospital *HospitalAdapter,
	ferry *FerryAdapter,
	hotel *HotelAdapter,
	transport *TransportAdapter,
) *Aggregator {
	return &Aggregator{
		hospital:  hospital,
		ferry:     ferry,
		hotel:     hotel,
		transport: transport,
	}
}

// SearchAll executes independent provider queries concurrently, isolating failures per category.
func (a *Aggregator) SearchAll(ctx context.Context, reqID string, q MultiSearchQuery) *MultiSearchResult {
	result := &MultiSearchResult{
		HospitalOffers:  make([]Offer, 0),
		FerryOffers:     make([]Offer, 0),
		HotelOffers:     make([]Offer, 0),
		TransportOffers: make([]Offer, 0),
		Warnings:        make([]string, 0),
	}

	var mu sync.Mutex
	var wg sync.WaitGroup

	// 1. Hospital Search
	if q.HospitalCriteria != nil {
		wg.Add(1)
		go func(crit HospitalSearchCriteria) {
			defer wg.Done()
			offers, err := a.hospital.Search(ctx, reqID, crit)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				result.Warnings = append(result.Warnings, fmt.Sprintf("hospital search failed: %v", err))
			} else {
				result.HospitalOffers = append(result.HospitalOffers, offers...)
			}
		}(*q.HospitalCriteria)
	}

	// 2. Ferry Searches
	for _, fc := range q.FerryCriteria {
		wg.Add(1)
		go func(crit FerrySearchCriteria) {
			defer wg.Done()
			offers, err := a.ferry.Search(ctx, reqID, crit)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				result.Warnings = append(result.Warnings, fmt.Sprintf("ferry search failed (%s->%s): %v", crit.OriginPortCode, crit.DestinationPortCode, err))
			} else {
				result.FerryOffers = append(result.FerryOffers, offers...)
			}
		}(fc)
	}

	// 3. Hotel Search
	if q.HotelCriteria != nil {
		wg.Add(1)
		go func(crit HotelSearchCriteria) {
			defer wg.Done()
			offers, err := a.hotel.Search(ctx, reqID, crit)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				result.Warnings = append(result.Warnings, fmt.Sprintf("hotel search failed: %v", err))
			} else {
				result.HotelOffers = append(result.HotelOffers, offers...)
			}
		}(*q.HotelCriteria)
	}

	// 4. Transport Searches
	for _, tc := range q.TransportCriteria {
		wg.Add(1)
		go func(crit TransportSearchCriteria) {
			defer wg.Done()
			offers, err := a.transport.Search(ctx, reqID, crit)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				result.Warnings = append(result.Warnings, fmt.Sprintf("transport search failed (%s->%s): %v", crit.PickupLocationCode, crit.DropoffLocationCode, err))
			} else {
				result.TransportOffers = append(result.TransportOffers, offers...)
			}
		}(tc)
	}

	wg.Wait()
	return result
}
