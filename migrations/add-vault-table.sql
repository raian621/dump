CREATE TYPE vault_type AS ENUM ('S3_BUCKET', 'GCS_BUCKET', 'SELF_HOSTED');

CREATE TABLE vaults (
  id         SERIAL PRIMARY KEY,
  owner_id   INTEGER REFERENCES users(id) ON DELETE CASCADE,
  vault_name VARCHAR(200), -- Client vault name
  vault_type VAULT_TYPE
);