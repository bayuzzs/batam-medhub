package adapter

import (
	"context"
	"net/http"
	"time"
)

// TransportAdapter provides typed HTTP integration with the mock ground transport provider service.
type TransportAdapter struct {
	client *Client
}

// NewTransportAdapter constructs a new TransportAdapter.
func NewTransportAdapter(baseURL, integrationKey string, timeout time.Duration) *TransportAdapter {
	return &TransportAdapter{
		client: NewClient(baseURL, integrationKey, timeout),
	}
}

// Health queries the provider service health.
func (a *TransportAdapter) Health(ctx context.Context, reqID string) (*HealthResponse, error) {
	var resp HealthResponse
	err := a.client.Do(ctx, http.MethodGet, "/healthz", reqID, "", nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// Search queries available ground transport offers.
func (a *TransportAdapter) Search(ctx context.Context, reqID string, criteria TransportSearchCriteria) ([]Offer, error) {
	req := SearchRequest{
		ProviderType: ProviderTypeTransport,
		Criteria:     criteria,
	}
	var resp SearchResponse
	err := a.client.Do(ctx, http.MethodPost, "/v1/offers/search", reqID, "", req, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Offers, nil
}

// CreateHold creates a capacity hold for a selected transport offer with booking requirements.
func (a *TransportAdapter) CreateHold(ctx context.Context, reqID, idemKey string, req CreateHoldRequest) (*Hold, error) {
	req.ProviderType = ProviderTypeTransport
	var hold Hold
	err := a.client.Do(ctx, http.MethodPost, "/v1/holds", reqID, idemKey, req, &hold)
	if err != nil {
		return nil, err
	}
	return &hold, nil
}

// ConfirmHold confirms a live transport hold into a reservation.
func (a *TransportAdapter) ConfirmHold(ctx context.Context, reqID, idemKey, holdID string) (*Reservation, error) {
	var res Reservation
	err := a.client.Do(ctx, http.MethodPost, "/v1/holds/"+holdID+"/confirm", reqID, idemKey, nil, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ReleaseHold releases an unconfirmed transport hold.
func (a *TransportAdapter) ReleaseHold(ctx context.Context, reqID, idemKey, holdID string) (*ReleaseResult, error) {
	var res ReleaseResult
	err := a.client.Do(ctx, http.MethodPost, "/v1/holds/"+holdID+"/release", reqID, idemKey, nil, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// GetReservation retrieves the status of a confirmed transport reservation.
func (a *TransportAdapter) GetReservation(ctx context.Context, reqID, reservationID string) (*Reservation, error) {
	var res Reservation
	err := a.client.Do(ctx, http.MethodGet, "/v1/reservations/"+reservationID, reqID, "", nil, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ReleaseReservation cancels or compensates a confirmed transport reservation.
func (a *TransportAdapter) ReleaseReservation(ctx context.Context, reqID, idemKey, reservationID string) (*ReleaseResult, error) {
	var res ReleaseResult
	err := a.client.Do(ctx, http.MethodPost, "/v1/reservations/"+reservationID+"/release", reqID, idemKey, nil, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
