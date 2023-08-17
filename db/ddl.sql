CREATE TABLE event(
  event_name VARCHAR(40),
  event_id VARCHAR(200),
  correlation_id VARCHAR(200),
  user_id VARCHAR(200),
  dt_created DATETIME(6),
  dt_received DATETIME(6), 
  dt_ingested DATETIME(6),
  labels VARCHAR(200),
  payload JSON,
  SORT KEY(dt_created),
  SHARD(event_name),
  FULLTEXT (labels)
);

CREATE INDEX idx_event_id ON event(event_id) USING HASH;
CREATE INDEX idx_event_correlation ON event(correlation_id) USING HASH;



drop table n_event;
CREATE TABLE n_event(
  event_name VARCHAR(40),
  event_id VARCHAR(200),
  correlation_id VARCHAR(200),
  user_id VARCHAR(200),
  dt_created DATETIME(6),
  dt_received DATETIME(6),
  dt_ingested DATETIME(6),
  labels VARCHAR(200),
  payload JSON,
  SORT KEY (dt_created),
  SHARD KEY (user_id),
  FULLTEXT (labels)
);

CREATE UNIQUE INDEX idx_n_event_id ON n_event(event_id) USING HASH;
CREATE INDEX idx_n_event_user ON n_event(user_id) USING HASH;
CREATE INDEX idx_n_event_correlation ON n_event(correlation_id) USING HASH;

ALTER TABLE n_event ADD value AS payload::$value PERSISTED DOUBLE;
ALTER TABLE n_event ADD dt_created_min AS date_trunc('minute', dt_created) PERSISTED DATETIME(6);

CREATE INDEX idx_n_event_dt_min ON n_event(event_name, dt_created_min) USING HASH;

INSERT INTO n_event SELECT * FROM event WHERE event_name '' amd '';

insert into n_event(
  event_name, event_id, correlation_id, user_id, 
  dt_created, dt_received, dt_ingested,
  labels, payload
)
select 
  event_name, event_id, correlation_id, user_id, 
  as_date(dt_created), as_date(dt_received), as_date(dt_ingested),
  array_as_string(split(labels, ",")), payload
from event
where
  to_number(substring(event_name, 8)) between 496 and 500;
  