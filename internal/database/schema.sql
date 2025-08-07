create table if not exists messages
(
    id         numeric(20, 0) not null primary key,
    thread_id  numeric(20, 0),
    account_id numeric(20, 0),
    created_at timestamptz    not null default now(),
    mentions   jsonb          not null default '[]',
    content    text           not null
);