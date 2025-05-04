create table if not exists cities(
    id bigserial PRIMARY KEY,
    name varchar(255) not null,
    sub_name varchar(255) default null,
    company_id bigint not null,
    created_at timestamp(0) with time zone not null default now()
)

