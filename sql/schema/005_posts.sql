-- +goose Up
create table posts (
    id serial primary key,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    title text not null,
    url text unique not null,
    description text,
    published_at timestamp,
    feed_id int not null,
    foreign key (feed_id) references feeds(id)
);

-- +goose Down
drop table posts;