CREATE TABLE IF NOT EXISTS transactions (
    id bigserial PRIMARY KEY,
    amount bigint NOT NULL,
    service_fee bigint DEFAULT NULL,
    from_currency_type_id bigint,
    to_currency_type_id bigint,
    sender_id bigint,
    from_city_id bigint,
    to_city_id bigint,
    receiver_name varchar(255) NOT NULL,
    receiver_phone varchar(9) NOT NULL,
    status bigint NOT NULL,
    company_id bigint,
    balance_id bigint,
    details varchar(255) NOT NULL,
    type bigint NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT now()
);

ALTER TABLE transactions ADD CONSTRAINT fk_transactions_sender_id FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE SET NULL;
ALTER TABLE transactions ADD CONSTRAINT fk_transactions_from_currency_type_id FOREIGN KEY (from_currency_type_id) REFERENCES currencies(id) ON DELETE SET NULL;
ALTER TABLE transactions ADD CONSTRAINT fk_transactions_to_currency_type_id FOREIGN KEY (to_currency_type_id) REFERENCES currencies(id) ON DELETE SET NULL;
ALTER TABLE transactions ADD CONSTRAINT fk_transactions_from_city_id FOREIGN KEY (from_city_id) REFERENCES cities(id) ON DELETE SET NULL;
ALTER TABLE transactions ADD CONSTRAINT fk_transactions_to_city_id FOREIGN KEY (to_city_id) REFERENCES cities(id) ON DELETE SET NULL;
ALTER TABLE transactions ADD CONSTRAINT fk_transactions_company_id FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE SET NULL;
ALTER TABLE transactions ADD CONSTRAINT fk_transactions_balance_id FOREIGN KEY (balance_id) REFERENCES balances(id) ON DELETE SET NULL;
