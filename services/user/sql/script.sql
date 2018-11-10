CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    uid UUID PRIMARY KEY,
    username VARCHAR(30) NOT NULL,
    password_hash CHAR(60) NOT NULL
);

CREATE TABLE user_posts (
  uid UUID REFERENCES users(uid),
  post_uid UUID NOT NULL UNIQUE,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL,
);
