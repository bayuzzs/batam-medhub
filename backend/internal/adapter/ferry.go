package adapter

import (
	"context"
	"net/http"
	"time"
)

// FerryAdapter provides typed HTTP integration with the mock ferry provider service.
type FerryAdapter struct {
	client *Client
}

// NewFerryAdapter constructs a new FerryAdapter.
func NewFerryAdapter(baseURL, integrationKey string, timeout time.Duration) *FerryAdapter {
	return &FerryAdapter{
		client: NewClient(baseURL, integrationKey, timeout),
	}
}

// Health queries the provider service health.
func (a *FerryAdapter) Health(ctx context.Context, reqID string) (*HealthResponse, error) {
	var resp HealthResponse
	err := a.client.Do(ctx, http.MethodGet, "/healthz", reqID, "", nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// Search queries available ferry sailing offers.
func (a *FerryAdapter) Search(ctx context.Context, reqID string, criteria FerrySearchCriteria) ([]Offer, error) {
	req := SearchRequest{
		ProviderType: ProviderTypeFerry,
		Criteria:     criteria,
	}
	var resp SearchResponse
	err := a.client.Do(ctx, http.MethodPost, "/v1/offers/search", reqID, "", req, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Offers, nil
}

// CreateHold creates a capacity hold for a selected ferry offer.
func (a *FerryAdapter) CreateHold(ctx context.Context, reqID, idemKey string, req CreateHoldRequest) (*Hold, error) {
	req.ProviderType = ProviderTypeFerry
	var hold Hold
	err := a.client.Do(ctx, http.MethodPost, "/v1/holds", reqID, idemKey, req, &hold)
	if err != nil {
		return nil, err
	}
	return &hold, nil
}

// ConfirmHold confirms a live ferry hold into a reservation.
func (a *FerryAdapter) ConfirmHold(ctx context.Context, reqID, idemKey, holdID string) (*Reservation, error) {
	var res Reservation
	err := a.client.Do(ctx, http.MethodPost, "/v1/holds/"+holdID+"/confirm", reqID, idemKey, nil, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ReleaseHold releases an unconfirmed ferry hold.
func (a *FerryAdapter) ReleaseHold(ctx context.Context, reqID, idemKey, holdID string) (*ReleaseResult, error) {
	var res ReleaseResult
	err := a.client.Do(ctx, http.MethodPost, "/v1/holds/"+holdID+"/release", reqID, idemKey, nil, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// GetReservation retrieves the status of a confirmed ferry reservation.
func (a *FerryAdapter) GetReservation(ctx context.Context, reqID, reservationID string) (*Reservation, error) {
	var res Reservation
	err := a.client.Do(ctx, http.MethodGet, "/v1/reservations/"+reservationID, reqID, "", nil, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ReleaseReservation cancels or compensates a confirmed ferry reservation.
func (a *FerryAdapter) ReleaseReservation(ctx context.Context, reqID, idemKey, reservationID string) (*ReleaseResult, error) {
	var res ReleaseResult
	err := a.client.Do(ctx, http.MethodPost, "/v1/reservations/"+reservationID+"/release", reqID, idemKey, nil, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
