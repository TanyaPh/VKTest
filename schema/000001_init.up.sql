CREATE TABLE users_checks
(
    id serial not null unique,
    checks int not null default 0
);