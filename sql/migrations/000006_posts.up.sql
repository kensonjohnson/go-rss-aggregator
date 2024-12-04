CREATE TABLE posts (
    id uuid PRIMARY KEY,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    title text NOT NULL,
    description text,
    published_at timestamp NOT NULL,
    url text NOT NULL UNIQUE,
    feed_id uuid NOT NULL REFERENCES feeds (id) ON DELETE CASCADE
);

