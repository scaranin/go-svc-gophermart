CREATE TABLE withdraws (
    withdraws_id SERIAL PRIMARY KEY,
    order_num TEXT,
    user_id NUMERIC NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);