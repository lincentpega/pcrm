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