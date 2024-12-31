CREATE TABLE apis (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title VARCHAR(255),
    description VARCHAR(255),
    method VARCHAR(255) NOT NULL,
    token_id INTEGER NOT NULL,
    body JSON NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    FOREIGN KEY (token_id) REFERENCES tokens (id) 
        ON UPDATE NO ACTION 
        ON DELETE NO ACTION
);