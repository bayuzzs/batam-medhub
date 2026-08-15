package adapter

import (
	"context"
	"net/http"
	"time"
)

// HotelAdapter provides typed HTTP integration with the mock hotel provider service.
type HotelAdapter struct {
	client *Client
}

// NewHotelAdapter constructs a new HotelAdapter.
func NewHotelAdapter(baseURL, integrationKey string, timeout time.Duration) *HotelAdapter {
	return &HotelAdapter{
		client: NewClient(baseURL, integrationKey, timeout),
	}
}

// Health queries the provider service health.
func (a *HotelAdapter) Health(ctx context.Context, reqID string) (*HealthResponse, error) {
	var resp HealthResponse
	err := a.client.Do(ctx, http.MethodGet, "/healthz", reqID, "", nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// Search queries available hotel room inventory offers.
func (a *HotelAdapter) Search(ctx context.Context, reqID string, criteria HotelSearchCriteria) ([]Offer, error) {
	req := SearchRequest{
		ProviderType: ProviderTypeHotel,
		Criteria:     criteria,
	}
	var resp SearchResponse
	err := a.client.Do(ctx, http.MethodPost, "/v1/offers/search", reqID, "", req, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Offers, nil
}

// CreateHold creates a capacity hold for a selected hotel offer.
func (a *HotelAdapter) CreateHold(ctx context.Context, reqID, idemKey string, req CreateHoldRequest) (*Hold, error) {
	req.ProviderType = ProviderTypeHotel
	var hold Hold
	err := a.client.Do(ctx, http.MethodPost, "/v1/holds", reqID, idemKey, req, &hold)
	if err != nil {
		return nil, err
	}
	return &hold, nil
}

// ConfirmHold confirms a live hotel hold into a reservation.
func (a *HotelAdapter) ConfirmHold(ctx context.Context, reqID, idemKey, holdID string) (*Reservation, error) {
	var res Reservation
	err := a.client.Do(ctx, http.MethodPost, "/v1/holds/"+holdID+"/confirm", reqID, idemKey, nil, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ReleaseHold releases an unconfirmed hotel hold.
func (a *HotelAdapter) ReleaseHold(ctx context.Context, reqID, idemKey, holdID string) (*ReleaseResult, error) {
	var res ReleaseResult
	err := a.client.Do(ctx, http.MethodPost, "/v1/holds/"+holdID+"/release", reqID, idemKey, nil, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// GetReservation retrieves the status of a confirmed hotel reservation.
func (a *HotelAdapter) GetReservation(ctx context.Context, reqID, reservationID string) (*Reservation, error) {
	var res Reservation
	err := a.client.Do(ctx, http.MethodGet, "/v1/reservations/"+reservationID, reqID, "", nil, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ReleaseReservation cancels or compensates a confirmed hotel reservation.
func (a *HotelAdapter) ReleaseReservation(ctx context.Context, reqID, idemKey, reservationID string) (*ReleaseResult, error) {
	var res ReleaseResult
	err := a.client.Do(ctx, http.MethodPost, "/v1/reservations/"+reservationID+"/release", reqID, idemKey, nil, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
