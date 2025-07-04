CREATE TABLE IF NOT EXISTS transactions (
    id bigserial PRIMARY KEY,
    marked_service_fee bigint DEFAULT 0,
    received_company_id bigint,
    received_user_id bigint,
    received_amount bigint,
    received_currency varchar(255),
    delivered_amount bigint,
    delivered_currency varchar(255),
    delivered_company_id bigint,
    delivered_user_id bigint,
    delivered_service_fee bigint DEFAULT 0,
    phone varchar(9) DEFAULT NULL,
    details varchar(255) DEFAULT NULL,
    status bigint,
    type bigint,
    created_at timestamp(0) with time zone NOT NULL DEFAULT now()
);
ALTER TABLE transactions ADD CONSTRAINT fk_transactions_delivered_user_id FOREIGN KEY (delivered_user_id) REFERENCES users(id) ON DELETE SET NULL;
ALTER TABLE transactions ADD CONSTRAINT fk_transactions_received_user_id FOREIGN KEY (received_user_id) REFERENCES users(id) ON DELETE SET NULL;
ALTER TABLE transactions ADD CONSTRAINT fk_transactions_received_company_id FOREIGN KEY (received_company_id) REFERENCES companies(id) ON DELETE SET NULL;
ALTER TABLE transactions ADD CONSTRAINT fk_transactions_delivered_company_id FOREIGN KEY (delivered_company_id) REFERENCES companies(id) ON DELETE SET NULL;