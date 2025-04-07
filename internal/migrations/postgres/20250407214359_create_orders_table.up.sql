CREATE TABLE orders (
    order_id SERIAL PRIMARY KEY,
    order_num TEXT,
    user_id NUMERIC NOT NULL,
    order_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    total_amount DECIMAL(10, 2) NOT NULL,
    status VARCHAR(50) NOT NULL,
    bonus_sum DECIMAL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user
        FOREIGN KEY (user_id)
        REFERENCES users(user_id)
);

CREATE INDEX idx_orders_user_id ON orders(user_id);