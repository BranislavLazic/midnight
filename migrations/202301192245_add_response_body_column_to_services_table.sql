-- +goose Up
alter table services
    add column response_body text;