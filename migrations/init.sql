-- CREATE EXTENSIONS
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Table: employees
CREATE TABLE employees (
    employee_id UUID NOT NULL DEFAULT uuid_generate_v4(),
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE employees ADD PRIMARY KEY (employee_id);

-- Table: wallets
CREATE TABLE wallets (
    employee_id UUID PRIMARY KEY,
    balance INTEGER DEFAULT 1000 CHECK (balance >= 0),
    FOREIGN KEY (employee_id) REFERENCES employees(employee_id)
);

-- Table: merch_items
CREATE TABLE merch_items (
    item_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price INTEGER NOT NULL
);

-- Table: purchases
CREATE TABLE purchases (
    purchase_id SERIAL PRIMARY KEY,
    employee_id UUID NOT NULL,
    item_id INTEGER NOT NULL,
    quantity INTEGER NOT NULL,
    purchase_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (employee_id) REFERENCES employees(employee_id),
    FOREIGN KEY (item_id) REFERENCES merch_items(item_id)
);

-- Table: transactions
CREATE TABLE transactions (
    transaction_id SERIAL PRIMARY KEY,
    sender_id UUID NOT NULL,
    receiver_id UUID NOT NULL,
    amount INTEGER NOT NULL,
    transaction_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (sender_id) REFERENCES employees(employee_id),
    FOREIGN KEY (receiver_id) REFERENCES employees(employee_id)
);


-- Indexes
CREATE INDEX idx_employees_email ON employees (email);
CREATE INDEX idx_employees_username ON employees (username);
CREATE INDEX idx_wallets_employee_id ON wallets (employee_id);
CREATE INDEX idx_merch_items_name ON merch_items (name);
CREATE INDEX idx_purchases_employee_id ON purchases (employee_id);
CREATE INDEX idx_transactions_sender_id ON transactions (sender_id);
CREATE INDEX idx_transactions_receiver_id ON transactions (receiver_id);

-- INIT metch data
INSERT INTO merch_items (name, price)
VALUES ('t-shirt', 80), ('cup', 20), ('book', 50), ('pen', 10), ('powerbank', 200),
    ('hoody', 300),('umbrella', 200),('socks', 10),('wallet', 50), ('pink-hoody', 500);


