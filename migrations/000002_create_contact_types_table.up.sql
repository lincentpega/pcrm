CREATE TABLE contact_types (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

INSERT INTO contact_types (name) VALUES 
    ('Email'),
    ('Phone'),
    ('Address'),
    ('Website'),
    ('Social Media');