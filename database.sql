-- This is the SQL script that will be used to initialize the database schema.
-- We will evaluate you based on how well you design your database.
-- 1. How you design the tables.
-- 2. How you choose the data types and keys.
-- 3. How you name the fields.
-- In this assignment we will use PostgreSQL as the database.

-- This is test table. Remove this table and replace with your own tables. 

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS estates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    name VARCHAR(100) NOT NULL,

    length_plot INT NOT NULL,

    width_plot INT NOT NULL,

    created_at TIMESTAMP DEFAULT now(),

    updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE IF NOT EXISTS trees (

    id BIGSERIAL PRIMARY KEY,

    estate_id BIGINT NOT NULL,

    plot_x INT NOT NULL,

    plot_y INT NOT NULL,

    height INT NOT NULL,

    created_at TIMESTAMP DEFAULT now(),

    CONSTRAINT fk_estate
        FOREIGN KEY (estate_id)
        REFERENCES estates(id),

    CONSTRAINT unique_tree
        UNIQUE(estate_id, plot_x, plot_y)
);

CREATE TABLE IF NOT EXISTS drone_simulations (

    id BIGSERIAL PRIMARY KEY,

    estate_id BIGINT NOT NULL,

    total_horizontal_distance INT,

    total_vertical_distance INT,

    total_distance INT,

    total_axis_turn INT,

    created_at TIMESTAMP DEFAULT now(),

    CONSTRAINT fk_drone_simulation_estate
        FOREIGN KEY (estate_id)
        REFERENCES estates(id)
);

INSERT INTO estates (id, name, length_plot, width_plot)
SELECT 1, 'Estate Alpha', 10, 10
WHERE NOT EXISTS (
    SELECT 1 FROM estates WHERE id = 1
);

INSERT INTO estates (id, name, length_plot, width_plot)
SELECT 2, 'Estate Beta', 12, 8
WHERE NOT EXISTS (
    SELECT 1 FROM estates WHERE id = 2
);

INSERT INTO trees (id, estate_id, plot_x, plot_y, height)
SELECT 1, 1, 1, 1, 12
WHERE NOT EXISTS (
    SELECT 1 FROM trees WHERE id = 1
);

INSERT INTO trees (id, estate_id, plot_x, plot_y, height)
SELECT 2, 1, 2, 3, 15
WHERE NOT EXISTS (
    SELECT 1 FROM trees WHERE id = 2
);

INSERT INTO trees (id, estate_id, plot_x, plot_y, height)
SELECT 3, 2, 4, 2, 9
WHERE NOT EXISTS (
    SELECT 1 FROM trees WHERE id = 3
);

INSERT INTO drone_simulations (id, estate_id, total_horizontal_distance, total_vertical_distance, total_distance, total_axis_turn)
SELECT 1, 1, 20, 30, 50, 4
WHERE NOT EXISTS (
    SELECT 1 FROM drone_simulations WHERE id = 1
);

INSERT INTO drone_simulations (id, estate_id, total_horizontal_distance, total_vertical_distance, total_distance, total_axis_turn)
SELECT 2, 2, 14, 22, 36, 3
WHERE NOT EXISTS (
    SELECT 1 FROM drone_simulations WHERE id = 2
);

SELECT setval(
    pg_get_serial_sequence('estates', 'id'),
    GREATEST(COALESCE((SELECT MAX(id) FROM estates), 1), 1)
);

SELECT setval(
    pg_get_serial_sequence('trees', 'id'),
    GREATEST(COALESCE((SELECT MAX(id) FROM trees), 1), 1)
);

SELECT setval(
    pg_get_serial_sequence('drone_simulations', 'id'),
    GREATEST(COALESCE((SELECT MAX(id) FROM drone_simulations), 1), 1)
);
