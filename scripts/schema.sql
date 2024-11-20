CREATE DATABASE cryptkeeper;

-- schema.sql

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS groups (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS user_groups (
    user_id INT NOT NULL,
    group_id INT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (group_id) REFERENCES groups (id),
    PRIMARY KEY (user_id, group_id)
);

CREATE TABLE IF NOT EXISTS secrets (
    id SERIAL PRIMARY KEY,
    path VARCHAR(255) NOT NULL,
    version INT NOT NULL,
    value TEXT NOT NULL,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS secret_access (
    secret_id INT NOT NULL,
    user_id INT,
    group_id INT,
    access_level VARCHAR(50) NOT NULL,
    FOREIGN KEY (secret_id) REFERENCES secrets (id),
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (group_id) REFERENCES groups (id),
    PRIMARY KEY (secret_id, user_id, group_id)
);


CREATE TABLE IF NOT EXISTS secret_deletions (
    id SERIAL PRIMARY KEY,
    secret_id INT NOT NULL,
    path VARCHAR(255) NOT NULL,
    version INT NOT NULL,
    value TEXT NOT NULL,
    metadata JSONB,
    deleted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);




ALTER TABLE secrets ADD CONSTRAINT unique_path_version UNIQUE (path, version);
