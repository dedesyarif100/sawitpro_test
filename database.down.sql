DELETE FROM drone_simulations
WHERE id IN (1, 2);

DELETE FROM trees
WHERE id IN (1, 2, 3);

DELETE FROM estates
WHERE id IN (1, 2);

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