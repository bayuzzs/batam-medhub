package adapter

import (
	"context"
	"net/http"
	"time"
)

// HospitalAdapter provides typed HTTP integration with the mock hospital provider service.
type HospitalAdapter struct {
	client *Client
}

// NewHospitalAdapter constructs a new HospitalAdapter.
func NewHospitalAdapter(baseURL, integrationKey string, timeout time.Duration) *HospitalAdapter {
	return &HospitalAdapter{
		client: NewClient(baseURL, integrationKey, timeout),
	}
}

// Health queries the provider service health.
func (a *HospitalAdapter) Health(ctx context.Context, reqID string) (*HealthResponse, error) {
	var resp HealthResponse
	err := a.client.Do(ctx, http.MethodGet, "/healthz", reqID, "", nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// Search queries available hospital appointment offers.
func (a *HospitalAdapter) Search(ctx context.Context, reqID string, criteria HospitalSearchCriteria) ([]Offer, error) {
	req := SearchRequest{
		ProviderType: ProviderTypeHospital,
		Criteria:     criteria,
	}
	var resp SearchResponse
	err := a.client.Do(ctx, http.MethodPost, "/v1/offers/search", reqID, "", req, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Offers, nil
}

// CreateHold creates a capacity hold for a selected hospital offer.
func (a *HospitalAdapter) CreateHold(ctx context.Context, reqID, idemKey string, req CreateHoldRequest) (*Hold, error) {
	req.ProviderType = ProviderTypeHospital
	var hold Hold
	err := a.client.Do(ctx, http.MethodPost, "/v1/holds", reqID, idemKey, req, &hold)
	if err != nil {
		return nil, err
	}
	return &hold, nil
}

// ConfirmHold confirms a live hospital hold into a reservation.
func (a *HospitalAdapter) ConfirmHold(ctx context.Context, reqID, idemKey, holdID string) (*Reservation, error) {
	var res Reservation
	err := a.client.Do(ctx, http.MethodPost, "/v1/holds/"+holdID+"/confirm", reqID, idemKey, nil, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ReleaseHold releases an unconfirmed hospital hold.
func (a *HospitalAdapter) ReleaseHold(ctx context.Context, reqID, idemKey, holdID string) (*ReleaseResult, error) {
	var res ReleaseResult
	err := a.client.Do(ctx, http.MethodPost, "/v1/holds/"+holdID+"/release", reqID, idemKey, nil, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// GetReservation retrieves the status of a confirmed hospital reservation.
func (a *HospitalAdapter) GetReservation(ctx context.Context, reqID, reservationID string) (*Reservation, error) {
	var res Reservation
	err := a.client.Do(ctx, http.MethodGet, "/v1/reservations/"+reservationID, reqID, "", nil, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ReleaseReservation cancels or compensates a confirmed hospital reservation.
func (a *HospitalAdapter) ReleaseReservation(ctx context.Context, reqID, idemKey, reservationID string) (*ReleaseResult, error) {
	var res ReleaseResult
	err := a.client.Do(ctx, http.MethodPost, "/v1/reservations/"+reservationID+"/release", reqID, idemKey, nil, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
