CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS moddatetime;

CREATE TABLE comments (
    uid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    post_uid UUID NOT NULL,
    body TEXT NOT NULL,
    parent_uid UUID,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp,
    modified_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp
);

CREATE TRIGGER posts_moddatetime
    BEFORE UPDATE ON comments
    FOR EACH ROW
    EXECUTE PROCEDURE moddatetime(modified_at);
