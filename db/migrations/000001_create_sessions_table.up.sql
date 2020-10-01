CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS sessions(
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id integer NOT NULL,
  created_at TIMESTAMP,
  last_seen TIMESTAMP
)