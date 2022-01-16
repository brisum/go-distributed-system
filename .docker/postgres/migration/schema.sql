-- CREATE DATABASE snapshot_sourcing
-- \i /migration/schema.sql

-- Entity
DROP TABLE IF EXISTS entity;
CREATE TABLE entity
(
    entity_type    varchar(1000) not null,
    entity_uuid    varchar(36)   not null,
    entity_version integer       not null
);
ALTER TABLE ONLY entity ADD CONSTRAINT entity_entity_uuid PRIMARY KEY (entity_uuid);

-- event
DROP TABLE IF EXISTS event;
CREATE TABLE event
(
    event_id         serial,
    event_type       varchar(1000) not null,
    entity_type      varchar(1000) not null,
    entity_uuid      varchar(36)   not null,
    event_data       varchar(1000) not null,
    promoter         varchar(1000) default null,
    triggering_event varchar(1000) default null,
    created_at       timestamp WITHOUT TIME ZONE DEFAULT timezone('utc' :: TEXT, now())
);
ALTER TABLE ONLY event ADD CONSTRAINT event_event_id PRIMARY KEY (event_id);
CREATE SEQUENCE event_id_sequence
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    MAXVALUE 99999999 CACHE 1
    OWNED BY event.event_id;

-- Snapshot
DROP TABLE IF EXISTS snapshot;
CREATE TABLE snapshot
(
    event_id         integer,
    entity_type      varchar(1000) not null,
    entity_uuid      varchar(36)   not null,
    snapshot_data    varchar(1000) not null,
    triggering_event varchar(1000) default null
);
ALTER TABLE ONLY snapshot ADD CONSTRAINT snapshot_entity_uuid PRIMARY KEY (entity_uuid);
