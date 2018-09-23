CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS moddatetime;

CREATE TABLE posts (
    uid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(80) NOT NULL,
    url VARCHAR(80),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp,
    modified_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp
);

CREATE TRIGGER posts_moddatetime
    BEFORE UPDATE ON posts
    FOR EACH ROW
    EXECUTE PROCEDURE moddatetime(modified_at);
