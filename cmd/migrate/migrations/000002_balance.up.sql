create table if not exists balances(
    id bigserial PRIMARY KEY,
    balance bigint default null,
    user_id bigint not null,
    in_out_lay bigint default null,
    out_in_lay bigint default null,
    company_id bigint not null,
    created_at timestamp(0) with time zone not null default now(),
    updated_at timestamp(0) with time zone not null default now()
)

