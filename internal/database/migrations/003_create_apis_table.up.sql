CREATE TABLE apis (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title VARCHAR(255),
    description VARCHAR(255),
    method VARCHAR(255) NOT NULL,
    path VARCHAR(255) NOT NULL,
    token_id INTEGER NOT NULL,
    body JSON NOT NULL,
    private TINYINT NOT NULL DEFAULT 0,
    entity VARCHAR(255),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    FOREIGN KEY (token_id) REFERENCES tokens (id) 
        ON UPDATE NO ACTION 
        ON DELETE NO ACTION
);