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
        'hospital-demo-v1', 'db0e409c75d1271be13c172285e6bd060c7ecc6710cf35e4965dbad963015b6f',
        'SHA256', 'ACTIVE', NULL, true, '2026-08-15T00:00:00Z', '2026-08-15T00:00:00Z'),
    ('00000000-0000-4000-8000-000000000202', '00000000-0000-4000-8000-000000000102',
        'ferry-demo-v1', 'e07f7fd6b8f901af3fc022555f6f8a54425afd194b58e00a03d8610ad93ebaea',
        'SHA256', 'ACTIVE', NULL, true, '2026-08-15T00:00:00Z', '2026-08-15T00:00:00Z'),
    ('00000000-0000-4000-8000-000000000203', '00000000-0000-4000-8000-000000000103',
        'hotel-demo-v1', '5fcdf35b2ebf1ac6b1a3c75d06796733e3712a72c33e31fa128128c93a3b0d9d',
        'SHA256', 'ACTIVE', NULL, true, '2026-08-15T00:00:00Z', '2026-08-15T00:00:00Z'),
    ('00000000-0000-4000-8000-000000000204', '00000000-0000-4000-8000-000000000104',
        'transport-demo-v1', '7e9cafe0d0fcd3bf10e1995134bab6557ed2e11abb30fbe89095e099ee5a98fa',
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
