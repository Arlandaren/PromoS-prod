CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       name VARCHAR(100) NOT NULL,
                       email VARCHAR(100) UNIQUE NOT NULL,
                       password VARCHAR(100) NOT NULL
);

CREATE TABLE companies (
                           id SERIAL PRIMARY KEY,
                           name VARCHAR(100) NOT NULL,
                           email VARCHAR(100) UNIQUE NOT NULL,
                           password VARCHAR(100) NOT NULL
);

CREATE INDEX idx_users_email ON users(email);

CREATE INDEX idx_companies_email ON companies(email);