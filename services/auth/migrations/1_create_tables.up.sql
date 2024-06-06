CREATE TABLE roles (
    id serial primary key,
    name varchar(255) unique
);

INSERT INTO roles (id, name) VALUES (1, 'user'), (2, 'admin');

CREATE TABLE users (
    id bigserial primary key,
    role_id int not null,
    login varchar(255) unique not null,
    password varchar(255) not null,

    foreign key(role_id) references roles(id)
);