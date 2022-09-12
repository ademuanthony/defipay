drop table if exists account_transaction;
drop table if exists daily_earning;
drop table if exists investment;
drop table if exists subscription;
drop table if exists transfer;
drop table if exists wallet;
drop table if exists weekly_payout;
drop table if exists withdrawal;
drop table if exists deposit;
drop table if exists account;


CREATE TABLE IF NOT EXISTS account 
(
    id character varying(64) NOT NULL PRIMARY KEY,
    username character varying(256) NOT NULL UNIQUE,
    password character varying(256) NOT NULL,
    email character varying(256) not null,
    phone_number character varying(32) not null,
    created_at bigint NOT NULL,
    first_name character varying(256) NOT NULL default '',
    last_name character varying(256) NOT NULL default '',
    referral_id character varying(256) default '',
    withdrawal_addresss character varying(256) NOT NULL default '',
    balance bigint NOT NULL default 0,
    role int default 0
);

CREATE TABLE IF NOT EXISTS wallet 
(
    id character varying(64) NOT NULL PRIMARY KEY,
    address character varying(64) NOT NULL UNIQUE,
    private_key character varying(124) NOT NULL UNIQUE,
    coin_symbol character varying(32) NOT NULL,
    account_id character varying(64) NOT NULL REFERENCES account(id)
);

CREATE TABLE IF NOT EXISTS account_transaction (
	id serial,
	account_id character varying(64) NOT NULL REFERENCES account(id),
	amount INT8 NOT NULL,
	tx_type VARCHAR(32) NOT NULL,
	opening_balance INT8 NOT NULL,
	closing_balance INT8 NOT NULL,
	date INT8 NOT NULL,
	description VARCHAR(256) NOT NULL,
	PRIMARY KEY(id)
);

create table if not exists user_setting (
    id uuid not null default gen_random_uuid() primary key,
    account_id character varying(64) not null references account(id),
    config_key character varying(64) not null,
    config_value character varying(500) not null,
    unique(account_id, config_key)
);

create table if not exists login_info (
    id uuid not null default gen_random_uuid() primary key,
    account_id character varying(64) not null references account(id),
    date bigint not null,
    ip character varying(128) not null,
    platform character varying(128) not null
);

CREATE TABLE IF NOT EXISTS security_code (
	id uuid not null default gen_random_uuid() primary key,
    account_id character varying(64) not null references account(id),
	code VARCHAR(3200) NOT NULL,
	date bigint NOT NULL
);

create table if not exists notification (
    id character varying(64) not null primary key,
    account_id character varying(64) not null references account(id),
    status int not null,
    title character varying(128) not null,
    content character varying(500) not null,
    date bigint not null,
    type int not null default 0,
    action_link character varying(64) not null default 0,
    action_text character varying(32) not null default 0
);

create table if not exists transaction (
    id uuid not null default gen_random_uuid() primary key,
    account_id character varying(64) not null references account(id),
)
