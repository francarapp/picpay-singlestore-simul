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
