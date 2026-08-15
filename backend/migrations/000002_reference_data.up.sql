INSERT INTO providers (
    id, provider_type, code, display_name, status, synthetic, source, created_at, updated_at
) VALUES
    ('00000000-0000-4000-8000-000000000101', 'HOSPITAL', 'hospital-demo-01',
        'Synthetic Batam MedHub Hospital', 'ACTIVE', true, 'MOCK',
        '2026-08-15T00:00:00Z', '2026-08-15T00:00:00Z'),
    ('00000000-0000-4000-8000-000000000102', 'FERRY', 'ferry-demo-01',
        'Synthetic Cross-Strait Ferry', 'ACTIVE', true, 'MOCK',
        '2026-08-15T00:00:00Z', '2026-08-15T00:00:00Z'),
    ('00000000-0000-4000-8000-000000000103', 'HOTEL', 'hotel-demo-01',
        'Synthetic Batam Centre Hotel', 'ACTIVE', true, 'MOCK',
        '2026-08-15T00:00:00Z', '2026-08-15T00:00:00Z'),
    ('00000000-0000-4000-8000-000000000104', 'TRANSPORT', 'transport-demo-01',
        'Synthetic Batam MedHub Transport', 'ACTIVE', true, 'MOCK',
        '2026-08-15T00:00:00Z', '2026-08-15T00:00:00Z');

INSERT INTO provider_credentials (
    id, provider_id, key_prefix, secret_hash, hash_algorithm, status,
    expires_at, synthetic, created_at, updated_at
) VALUES
    ('00000000-0000-4000-8000-000000000201', '00000000-0000-4000-8000-000000000101',
        'hospital-demo-v1', '461dfce091195ef876ea27496743a176c4deacc8e50ba3d37126448869f88814',
        'SHA256', 'ACTIVE', NULL, true, '2026-08-15T00:00:00Z', '2026-08-15T00:00:00Z'),
    ('00000000-0000-4000-8000-000000000202', '00000000-0000-4000-8000-000000000102',
        'ferry-demo-v1', 'b291d2694b0540849da60bf10f26bfc4a2c325aeec93c5353f197deaddc941df',
        'SHA256', 'ACTIVE', NULL, true, '2026-08-15T00:00:00Z', '2026-08-15T00:00:00Z'),
    ('00000000-0000-4000-8000-000000000203', '00000000-0000-4000-8000-000000000103',
        'hotel-demo-v1', '29a39b089fe0653a10d3b282f8e6370f01ac95ef31f18c2dce735749f36070c1',
        'SHA256', 'ACTIVE', NULL, true, '2026-08-15T00:00:00Z', '2026-08-15T00:00:00Z'),
    ('00000000-0000-4000-8000-000000000204', '00000000-0000-4000-8000-000000000104',
        'transport-demo-v1', 'e01da70ce71e8788bde2456362b1b7d75ce3ed2171d3cde5509b6fe083e1c5fd',
        'SHA256', 'ACTIVE', NULL, true, '2026-08-15T00:00:00Z', '2026-08-15T00:00:00Z');

INSERT INTO medical_services (
    id, code, name, category, description, default_duration_minutes,
    active, synthetic, source, created_at, updated_at
) VALUES
    ('00000000-0000-4000-8000-000000000301', 'MCU_BASIC', 'Basic Medical Check-up',
        'PREVENTIVE_CHECKUP',
        'Synthetic basic planned check-up package for the hackathon demo.',
        120, true, true, 'MOCK', '2026-08-15T00:00:00Z', '2026-08-15T00:00:00Z'),
    ('00000000-0000-4000-8000-000000000302', 'MCU_COMPREHENSIVE', 'Comprehensive Medical Check-up',
        'PREVENTIVE_CHECKUP',
        'Synthetic comprehensive planned check-up package for the hackathon demo.',
        240, true, true, 'MOCK', '2026-08-15T00:00:00Z', '2026-08-15T00:00:00Z'),
    ('00000000-0000-4000-8000-000000000303', 'DENTAL_CHECKUP', 'Dental Check-up',
        'DENTAL_SCREENING',
        'Synthetic planned dental screening package for the hackathon demo.',
        60, true, true, 'MOCK', '2026-08-15T00:00:00Z', '2026-08-15T00:00:00Z'),
    ('00000000-0000-4000-8000-000000000304', 'EYE_SCREENING', 'Eye Screening',
        'VISION_SCREENING',
        'Synthetic planned vision screening package for the hackathon demo.',
        60, true, true, 'MOCK', '2026-08-15T00:00:00Z', '2026-08-15T00:00:00Z');

INSERT INTO provider_capabilities (
    id, provider_id, medical_service_id, external_service_id,
    active, synthetic, source, created_at, updated_at
) VALUES
    ('00000000-0000-4000-8000-000000000401', '00000000-0000-4000-8000-000000000101',
        '00000000-0000-4000-8000-000000000301', 'MCU_BASIC',
        true, true, 'MOCK', '2026-08-15T00:00:00Z', '2026-08-15T00:00:00Z'),
    ('00000000-0000-4000-8000-000000000402', '00000000-0000-4000-8000-000000000101',
        '00000000-0000-4000-8000-000000000302', 'MCU_COMPREHENSIVE',
        true, true, 'MOCK', '2026-08-15T00:00:00Z', '2026-08-15T00:00:00Z'),
    ('00000000-0000-4000-8000-000000000403', '00000000-0000-4000-8000-000000000101',
        '00000000-0000-4000-8000-000000000303', 'DENTAL_CHECKUP',
        true, true, 'MOCK', '2026-08-15T00:00:00Z', '2026-08-15T00:00:00Z'),
    ('00000000-0000-4000-8000-000000000404', '00000000-0000-4000-8000-000000000101',
        '00000000-0000-4000-8000-000000000304', 'EYE_SCREENING',
        true, true, 'MOCK', '2026-08-15T00:00:00Z', '2026-08-15T00:00:00Z');

INSERT INTO fx_rates (
    id, base_currency, quote_currency, rate, source, effective_at,
    estimated, synthetic, created_at
) VALUES
    ('00000000-0000-4000-8000-000000000501', 'SGD', 'SGD', 1.000000000000,
        'DEMO_STATIC_2026_08', '2026-08-01T00:00:00Z', true, true, '2026-08-01T00:00:00Z'),
    ('00000000-0000-4000-8000-000000000502', 'IDR', 'IDR', 1.000000000000,
        'DEMO_STATIC_2026_08', '2026-08-01T00:00:00Z', true, true, '2026-08-01T00:00:00Z'),
    ('00000000-0000-4000-8000-000000000503', 'IDR', 'SGD', 0.000084388200,
        'DEMO_STATIC_2026_08', '2026-08-01T00:00:00Z', true, true, '2026-08-01T00:00:00Z'),
    ('00000000-0000-4000-8000-000000000504', 'SGD', 'IDR', 11850.000000000000,
        'DEMO_STATIC_2026_08', '2026-08-01T00:00:00Z', true, true, '2026-08-01T00:00:00Z');
