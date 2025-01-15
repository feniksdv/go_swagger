CREATE TABLE entity_fields (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    entity_id INTEGER NOT NULL,
    field_name VARCHAR(255) NOT NULL,
    field_desc VARCHAR(255),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    FOREIGN KEY (entity_id) REFERENCES entities (id) 
        ON UPDATE NO ACTION 
        ON DELETE NO ACTION
);