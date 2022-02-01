CREATE TABLE IF NOT EXISTS users_balances (
  user_id bigint not null,
  amount bigint not null
);

INSERT INTO users_balances VALUES (11, 5000), (21, 3000), (31, 1000);


CREATE TABLE IF NOT EXISTS users_transactions (
  id bigserial primary key,
  user_id bigint not null,
  amount bigint not null,
  event int not null,
  transfer_id bigint not null,
  message text not null default '',
  created_at timestamp not null
);

INSERT INTO users_transactions (user_id, amount, event, transfer_id, created_at)
VALUES (11, 5000, 2, 11, now()), (21, 3000, 2, 21, now()), (31, 1000, 2, 31, now());