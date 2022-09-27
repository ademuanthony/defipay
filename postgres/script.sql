drop table if exists account_transaction;
drop table if exists wallet;
drop table if exists deposit;
drop table if exists account;
drop table if exists user_setting;
drop table if exists login_info;
drop table if exists security_code;
drop table if exists notification;
drop table if exists transaction;
drop table if exists payment_link;
drop table if exists beneficiary;

CREATE TABLE IF NOT EXISTS account 
(
    id character varying(64) NOT NULL PRIMARY KEY,
    email character varying(256) not null UNIQUE,
    referral_code character varying(256) not null UNIQUE,
    password character varying(256) NOT NULL,
    phone_number character varying(32) not null,
    created_at bigint NOT NULL,
    first_name character varying(256) NOT NULL default '',
    last_name character varying(256) NOT NULL default '',
    referral_id character varying(256) default '',
    withdrawal_addresss character varying(256) NOT NULL default '',
    balance bigint NOT NULL default 0,
    role int default 0
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
    bank_name character varying(64) not null,
    account_number character varying(64) not null,
    account_name character varying(64) not null,
    amount int8 not null,
    amount_paid character varying(64) not null,
    email character varying(64) not null,
    network character varying(64) not null,
    currency character varying(64) not null,
    wallet_address character varying(64) not null,
    private_key character varying(64) not null,
    payment_link character varying(64) not null,
    type character varying(28) not null,
    status character varying(28) not null
);

create table if not exists payment_link ( 
    permalink character(64) not null primary key,
    account_id character varying(64) references account(id),
    email  character varying(64) not null,
    accountName character varying(64) not null,
    accountNumber character varying(64) not null,
    bankName character varying(64) not null,
    fixAmount int8 not null,
    title character varying(100) not null,
    description character varying(500) not null
);

create table if not exists beneficiary (
    id uuid not null default gen_random_uuid() primary key,
    account_id character varying(64) references account(id),
    bank character varying(64) not null,
    account_number character varying(64) not null,
    account_name character varying(64) not null,
    country character varying(64) not null,
    beneficial_email character varying(64) not null
);

create table if not exists agent (
    id serial not null primary key,
    slack_username character varying(64) not null unique,
    name character varying(64) not null,
    balance int8 not null,
    status int not null
);

create table if not exists transaction_assignment(
    id serial not null primary key,
    agent_id int not null references agent(id),
    transaction_id uuid not null references transaction(id),
    amount int8 not null,
    date int8 not null,
    status int noy null
);

alter table account add referral_code character varying(256) not null UNIQUE;
alter table transaction add status character varying(28) not null;
alter table transaction add amount_paid character varying(64) not null;
alter table transaction add token_amount character varying(64) not null;
