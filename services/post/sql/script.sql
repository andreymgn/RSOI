CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE posts (
    uid UUID PRIMARY KEY,
    title VARCHAR(80) NOT NULL,
    url VARCHAR(80),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    modified_at TIMESTAMP WITH TIME ZONE NOT NULL
);
