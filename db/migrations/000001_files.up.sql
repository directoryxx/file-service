CREATE TABLE IF NOT EXISTS files (
	id serial PRIMARY KEY,
    uuid uuid,
	name varchar,
	url varchar,
    user_id integer,
    is_temp integer,
    created_at timestamp,
    updated_at timestamp
);