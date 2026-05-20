-- =============================================================================
-- seed.sql  –  Initial data for maintenance application
-- Executed automatically by MySQL docker-entrypoint-initdb.d after init.sql
-- =============================================================================

-- Use the correct database (set by MYSQL_DATABASE env var)
USE maintenance_db;

-- ── Users ─────────────────────────────────────────────────────────────────────
INSERT INTO users (username, role) VALUES
    ('admin_sa',       'SA'),
    ('budi_approval',  'APPROVAL'),
    ('satria_approval', 'APPROVAL');

-- ── Vehicles ──────────────────────────────────────────────────────────────────
INSERT INTO vehicles (license_plate, model) VALUES
    ('AB 1234 CD', 'Toyota Avanza 2021'),
    ('AB 5678 EF', 'Honda Brio 2022'),
    ('AB 9012 GH', 'Mitsubishi Xpander 2023');

-- ── Master Items ──────────────────────────────────────────────────────────────
INSERT INTO master_items (item_name, type, price) VALUES
    ('Oli Mesin 5W-30 1L',       'PART',    75000.00),
    ('Filter Oli',               'PART',    45000.00),
    ('Ban Radial 185/65 R15',    'PART',   650000.00),
    ('Ganti Oli & Filter',       'SERVICE',  85000.00),
    ('Spooring & Balancing',     'SERVICE', 150000.00);
