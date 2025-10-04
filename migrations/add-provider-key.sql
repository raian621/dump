CREATE TYPE provider_type_enum AS ENUM ('AWS', 'GCP', 'SELF');

CREATE TABLE provider_keys (
  id SERIAL PRIMARY KEY,
  user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
  provider_type PROVIDER_TYPE_ENUM
);