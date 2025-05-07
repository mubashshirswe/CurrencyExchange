create table if not exists cities(
    id bigserial PRIMARY KEY,
    name varchar(255) not null,
    parent_id bigint default null,
    company_id bigint not null,
    created_at timestamp(0) with time zone not null default now()
)

