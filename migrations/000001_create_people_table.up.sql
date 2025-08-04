CREATE TABLE people (
    id BIGSERIAL PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    second_name VARCHAR(255),
    middle_name VARCHAR(255),
    birth_year INTEGER,
    birth_month INTEGER CHECK (birth_month >= 1 AND birth_month <= 12),
    birth_day INTEGER CHECK (birth_day >= 1 AND birth_day <= 31),
    approximate_age INTEGER,
    approximate_age_updated_at DATE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

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

CREATE TABLE contacts (
    id BIGSERIAL PRIMARY KEY,
    person_id BIGINT NOT NULL REFERENCES people(id) ON DELETE CASCADE,
    contact_type_id BIGINT NOT NULL REFERENCES contact_types(id) ON DELETE RESTRICT,
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_contacts_person_id ON contacts(person_id);
CREATE INDEX idx_contacts_contact_type_id ON contacts(contact_type_id);