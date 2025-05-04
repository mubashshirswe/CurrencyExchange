create table if not exists currencies(
    id bigserial PRIMARY KEY,
    name varchar(255) not null,
    sell bigint default null,
    buy bigint default null,
    company_id bigint not null,
    created_at timestamp(0) with time zone not null default now()
)

