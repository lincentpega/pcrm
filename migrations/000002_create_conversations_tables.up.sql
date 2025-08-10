CREATE TABLE conversation_types (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

INSERT INTO conversation_types (name) VALUES 
    ('Phone Call'),
    ('Video Call'),
    ('In-Person Meeting'),
    ('Text Message'),
    ('Email'),
    ('Social Media');

CREATE TABLE conversations (
    id BIGSERIAL PRIMARY KEY,
    person_id BIGINT NOT NULL REFERENCES people(id) ON DELETE CASCADE,
    conversation_type_id BIGINT NOT NULL REFERENCES conversation_types(id) ON DELETE RESTRICT,
    initiator VARCHAR(16) NOT NULL CHECK (initiator IN ('owner','person')),
    notes TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_conversations_person_id ON conversations(person_id);
CREATE INDEX idx_conversations_conversation_type_id ON conversations(conversation_type_id);

