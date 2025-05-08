CREATE TABLE IF NOT EXISTS cities (
    id bigserial PRIMARY KEY,
    name varchar(255) NOT NULL,
    parent_id bigint DEFAULT NULL,
    company_id bigint NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT now()
);

ALTER TABLE cities ADD CONSTRAINT fk_cities_parent_id FOREIGN KEY (parent_id) REFERENCES cities(id) ON DELETE SET NULL;
ALTER TABLE cities ADD CONSTRAINT fk_cities_company_id FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE SET NULL;