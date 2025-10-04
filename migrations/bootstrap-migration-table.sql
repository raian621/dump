CREATE TABLE migrations (
  id SERIAL PRIMARY KEY,
  migration VARCHAR(500)
);

INSERT INTO migrations (migration) VALUES ('bootstrap-migration-table.sql');
