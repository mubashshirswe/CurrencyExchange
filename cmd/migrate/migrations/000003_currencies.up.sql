CREATE TABLE IF NOT EXISTS currencies (
    id bigserial PRIMARY KEY,
    name varchar(255) NOT NULL UNIQUE,
    sell bigint DEFAULT NULL,
    buy bigint DEFAULT NULL,
    company_id bigint,
    created_at timestamp(0) with time zone NOT NULL DEFAULT now()
);

ALTER TABLE currencies
    ADD CONSTRAINT fk_currencies_company_id FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE SET NULL;
