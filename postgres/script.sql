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
    principal bigint NOT NULL default 0,
    matured_principal bigint not null default 0
);

CREATE TABLE IF NOT EXISTS wallet 
(
    id character varying(64) NOT NULL PRIMARY KEY,
    address character varying(64) NOT NULL UNIQUE,
    private_key character varying(124) NOT NULL UNIQUE,
    coin_symbol character varying(32) NOT NULL,
    account_id character varying(64) NOT NULL REFERENCES account(id)
);

CREATE TABLE IF NOT EXISTS package 
(
    id character varying(64) NOT NULL PRIMARY KEY,
    name character varying(64) NOT NULL UNIQUE,
    price BIGINT NOT NULL,
    min_return_per_month INT NOT NULL,
    max_return_per_month INT NOT NULL,
    trades_per_day INT NOT NULL,
    accuracy INT NOT NULL
);

CREATE TABLE IF NOT EXISTS subscription
(
    id character varying(64) not null primary key,
    account_id character varying(64) not null references account(id),
    package_id character varying(64) not null references package(id),
    start_date bigint not null,
    end_date bigint not null
);

create table if not exists daily_earning
(
    id serial not null primary key,
    account_id character varying(64) not null references account(id),
    date bigint not null,
    percentage int not null,
    principal bigint not null,
    unique(account_id, date)
);

create table if not exists deposit
(
    id character varying(64) not null primary key,
    amount bigint not null,
    account_id character varying(64) not null references account(id),
    date bigint not null,
    ref character varying(256) not null unique
);

create table if not exists transfer
(
    id character varying(64) not null primary key,
    amount bigint not null,
    sender_id character varying(64) not null references account(id),
    receiver_id character varying(64) not null references account(id),
    date bigint not null
);

create table if not exists withdrawal 
(
    id character varying(64) not null primary key,
    account_id character varying(64) not null references account(id),
    amount bigint not null,
    date bigint not null,
    destination character varying(256) not null,
    ref character varying(256) not null,
    status character varying(32) not null
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

create table if not exists investment (
    id character varying(64) not null primary key,
    account_id character varying(64) not null references account(id),
    amount bigint not null,
    date bigint not null,
    activation_date bigint not null,
    status int not null default 0
);

create table if not exists weekly_payout (
    id character varying(63) not null primary key,
    date bigint not null,
    amount bigint not null
);

create table if not exists referral_payout (
    id character varying(64) not null primary key,
    account_id character varying(64) not null references account(id),
    from_account_id character varying(64) not null references account(id),
    subscription_id character varying(64) not null references subscription(id),
    generation int not null,
    amount bigint not null,
    date bigint not null,
    payment_method int not null,
    payment_status int not null,
    payment_ref string not null,
    unique(subscription_id, generation)
);

create table if not exists notification (
    id character varying(64) not null primary key,
    account_id character varying(64) not null references account(id),
    status int not null,
    title character varying(128) not null,
    content character varying(500) not null,
    date bigint not null
);

create table if not exists trade_schedule (
    id uuid not null default gen_random_uuid() primary key,
    account_id character varying(64) not null references account(id),
    trade_no int not null,
    total_trades int not null,
    date bigint not null,
    target_profit_percentage int not null default 0,
    start_date int not null default 0,
    unique(account_id, date, trade_no)
);

create table if not exists trade (
    id uuid not null default gen_random_uuid() primary key,
    account_id character varying(64) not null references account(id), 
    trade_no int not null,
    date bigint not null,
    start_date bigint not null,
    end_date bigint not null,
    amount bigint not null,
    profit bigint not null,
    unique(account_id, date, trade_no)
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

alter table account add referral_id_2 character varying(256) default '';
alter table account add referral_id_3 character varying(256) default '';
alter table account add role int default 0;
alter table package add icon character varying(256) default '';
alter table notification add type int not null default 0;
alter table notification add action_link character varying(64) not null default 0;
alter table notification add action_text character varying(32) not null default 0;

alter table trade add last_view_profit bigint not null default 0;
alter table trade add last_view_time bigint not null default 0;
