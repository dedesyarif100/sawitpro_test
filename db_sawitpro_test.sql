-- This is the SQL script that will be used to initialize the database schema.
-- We will evaluate you based on how well you design your database.
-- 1. How you design the tables.
-- 2. How you choose the data types and keys.
-- 3. How you name the fields.
-- In this assignment we will use PostgreSQL as the database.

-- This is test table. Remove this table and replace with your own tables. 

CREATE EXTENSION IF NOT EXISTS pgcrypto;

DROP TABLE IF EXISTS drone_simulations CASCADE;
DROP TABLE IF EXISTS trees CASCADE;
DROP TABLE IF EXISTS estates CASCADE;

CREATE TABLE estates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    name VARCHAR(100) NOT NULL,

    length_plot INT NOT NULL,

    width_plot INT NOT NULL,

    created_at TIMESTAMPTZ DEFAULT now(),

    updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS trees (

    id UUID PRIMARY KEY,

    estate_id UUID NOT NULL,

    plot_x INT NOT NULL,

    plot_y INT NOT NULL,

    height INT NOT NULL,

    created_at TIMESTAMPTZ DEFAULT now(),

    CONSTRAINT fk_estate
        FOREIGN KEY (estate_id)
        REFERENCES estates(id),

    CONSTRAINT unique_tree
        UNIQUE(estate_id, plot_x, plot_y)
);

CREATE TABLE IF NOT EXISTS drone_simulations (

    id UUID PRIMARY KEY,

    estate_id UUID NOT NULL,

    total_horizontal_distance INT,

    total_vertical_distance INT,

    total_distance INT,

    total_axis_turn INT,

    created_at TIMESTAMPTZ DEFAULT now(),

    CONSTRAINT fk_drone_simulation_estate
        FOREIGN KEY (estate_id)
        REFERENCES estates(id)
);

INSERT INTO estates (id, name, length_plot, width_plot)
SELECT '11111111-1111-1111-1111-111111111111'::uuid, 'Estate Alpha', 10, 10
WHERE NOT EXISTS (
    SELECT 1 FROM estates WHERE id = '11111111-1111-1111-1111-111111111111'::uuid
);

INSERT INTO estates (id, name, length_plot, width_plot)
SELECT '22222222-2222-2222-2222-222222222222'::uuid, 'Estate Beta', 12, 8
WHERE NOT EXISTS (
    SELECT 1 FROM estates WHERE id = '22222222-2222-2222-2222-222222222222'::uuid
);

INSERT INTO trees (id, estate_id, plot_x, plot_y, height)
SELECT '33333333-3333-3333-3333-333333333333'::uuid, '11111111-1111-1111-1111-111111111111'::uuid, 1, 1, 12
WHERE NOT EXISTS (
    SELECT 1 FROM trees WHERE id = '33333333-3333-3333-3333-333333333333'::uuid
);

INSERT INTO trees (id, estate_id, plot_x, plot_y, height)
SELECT '44444444-4444-4444-4444-444444444444'::uuid, '11111111-1111-1111-1111-111111111111'::uuid, 2, 3, 15
WHERE NOT EXISTS (
    SELECT 1 FROM trees WHERE id = '44444444-4444-4444-4444-444444444444'::uuid
);

INSERT INTO trees (id, estate_id, plot_x, plot_y, height)
SELECT '55555555-5555-5555-5555-555555555555'::uuid, '22222222-2222-2222-2222-222222222222'::uuid, 4, 2, 9
WHERE NOT EXISTS (
    SELECT 1 FROM trees WHERE id = '55555555-5555-5555-5555-555555555555'::uuid
);

INSERT INTO drone_simulations (id, estate_id, total_horizontal_distance, total_vertical_distance, total_distance, total_axis_turn)
SELECT '66666666-6666-6666-6666-666666666666'::uuid, '11111111-1111-1111-1111-111111111111'::uuid, 20, 30, 50, 4
WHERE NOT EXISTS (
    SELECT 1 FROM drone_simulations WHERE id = '66666666-6666-6666-6666-666666666666'::uuid
);

INSERT INTO drone_simulations (id, estate_id, total_horizontal_distance, total_vertical_distance, total_distance, total_axis_turn)
SELECT '77777777-7777-7777-7777-777777777777'::uuid, '22222222-2222-2222-2222-222222222222'::uuid, 14, 22, 36, 3
WHERE NOT EXISTS (
    SELECT 1 FROM drone_simulations WHERE id = '77777777-7777-7777-7777-777777777777'::uuid
);
