-- +goose Up
create table if not exists users
(
    id       serial primary key,
    email    varchar(255) not null,
    password varchar(255) not null,
    role     varchar(255) not null,
    enabled  boolean default true
);

insert into users (email, password, role)
values ('admin@admin.com', '$2a$12$2qhiZjWKMW5RInVFUBczfejgjcMT2fmBYxVI6rTEyBvHTR3rdcvEu', 'ROLE_ADMIN');