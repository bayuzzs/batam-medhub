INSERT INTO room_types (
    id,
    code,
    property_name,
    name,
    max_guests_per_room,
    accessibility,
    active,
    synthetic,
    source,
    created_at,
    updated_at
) VALUES (
    '30000000-0000-4000-8000-000000000001',
    'DELUXE_TWIN',
    'Synthetic Batam Centre Hotel',
    'DELUXE_TWIN',
    2,
    '[]'::jsonb,
    true,
    true,
    'MOCK',
    '2026-08-15T08:00:00Z',
    '2026-08-15T08:00:00Z'
);

INSERT INTO room_inventory_days (
    id,
    room_type_id,
    stay_date,
    rooms_total,
    price_amount_minor,
    price_currency,
    status,
    created_at,
    updated_at
) VALUES (
    '30000000-0000-4000-8000-000000000101',
    '30000000-0000-4000-8000-000000000001',
    '2026-08-22',
    8,
    85000000,
    'IDR',
    'AVAILABLE',
    '2026-08-15T08:00:00Z',
    '2026-08-15T08:00:00Z'
);

INSERT INTO hotel_offers (
    id,
    external_offer_id,
    external_inventory_id,
    room_type_id,
    check_in_date,
    check_out_date,
    service_starts_at,
    service_ends_at,
    start_time_zone,
    end_time_zone,
    valid_until,
    status,
    synthetic,
    source,
    created_at,
    updated_at
) VALUES (
    '30000000-0000-4000-8000-000000000201',
    'hotel-offer-batam-centre-20260822-1n',
    'hotel-room-deluxe-20260822',
    '30000000-0000-4000-8000-000000000001',
    '2026-08-22',
    '2026-08-23',
    '2026-08-22T07:00:00Z',
    '2026-08-23T05:00:00Z',
    'Asia/Jakarta',
    'Asia/Jakarta',
    '2026-08-22T03:00:00Z',
    'AVAILABLE',
    true,
    'MOCK',
    '2026-08-15T08:00:00Z',
    '2026-08-15T08:00:00Z'
);

INSERT INTO hotel_offer_nights (offer_pk, inventory_day_id, room_type_id)
VALUES (
    '30000000-0000-4000-8000-000000000201',
    '30000000-0000-4000-8000-000000000101',
    '30000000-0000-4000-8000-000000000001'
);
