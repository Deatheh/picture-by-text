INSERT INTO users (id, email, password, role_id, is_active)
SELECT 
    '11111111-1111-1111-1111-111111111111',
    'admin@example.com',
    '$2a$10$.tV6m8bEPAtpumOsxrFLcu6CDdz.1DWFXpQTvAc2CE/T2KQKqaAl2',
    1,
    true
WHERE NOT EXISTS (SELECT 1 FROM users WHERE email = 'admin@example.com');