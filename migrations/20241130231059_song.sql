-- +goose Up
-- +goose StatementBegin
create table if not exists song (
    "id" serial primary key,
    "name" varchar(255) not null,
    "group" varchar(255) not null,
    "text" text not null,
    "link" varchar(255) not null,
    "release_date" date not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table song;
-- +goose StatementEnd
