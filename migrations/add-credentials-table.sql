CREATE TABLE credentials (
  user_id  INTEGER REFERENCES users(id) ON DELETE CASCADE,
  username VARCHAR(500),
  passhash CHAR(100)
);