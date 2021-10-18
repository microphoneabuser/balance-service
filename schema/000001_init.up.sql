CREATE TABLE accounts
(
    id serial not null unique,
    balance int
);

CREATE TABLE transactions
(
    id serial not null unique,
    sender_id int references accounts(id) on delete cascade,
    recipient_id int references accounts(id) on delete cascade,
    amount int,
    description varchar(255),
    timestamp timestamp default current_timestamp
);

