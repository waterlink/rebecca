#!/usr/bin/env bash

# By providing 'PG_USER' and ('PG_PASS' or `PG_ASK_PASS`) you can
# control how this script will authenticate to local pg server.
PARAMS="-U ${PG_USER:-postgres}"
[[ -z "$PG_PASS" ]] || PGPASSWORD="$PG_PASS"
[[ -z "$PG_ASK_PASS" ]] || PARAMS="$PARAMS -W"

psql $PARAMS -c "create database rebecca_pg_test"
psql $PARAMS -c "create user rebecca_pg with superuser password 'rebecca_pg'"

psql $PARAMS rebecca_pg_test -c "drop table if exists people; create table people( id serial primary key, name varchar(50), age int )"
psql $PARAMS rebecca_pg_test -c "drop table if exists posts; create table posts( id serial primary key, title varchar(50), content text, created_at timestamp )"
