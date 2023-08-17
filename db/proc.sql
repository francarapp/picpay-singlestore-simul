
-- --------------------------------------------------------------------------------;
-- --------------------------------------------------------------------------------;

DELIMITER //
CREATE FUNCTION as_date(dt_created VARCHAR(24)) RETURNS datetime(6) AS
  BEGIN
    RETURN str_to_date(dt_created, '%Y-%m-%d %H:%i:%s.%f');
  END //
DELIMITER ;

DELIMITER //
CREATE FUNCTION created_day(dt_created VARCHAR(24)) RETURNS datetime(6) AS
  BEGIN
    RETURN date_trunc('day', str_to_date(dt_created, '%Y-%m-%d %H:%i:%s.%f'));
  END //
DELIMITER ;

DELIMITER //
CREATE FUNCTION created_hour(dt_created VARCHAR(24)) RETURNS datetime(6) AS
  BEGIN
    RETURN date_trunc('hour', str_to_date(dt_created, '%Y-%m-%d %H:%i:%s.%f'));
  END //
DELIMITER ;

DELIMITER //
CREATE FUNCTION as_created_hour(dt_created DATETIME(6)) RETURNS datetime(6) AS
  BEGIN
    RETURN date_trunc('hour', dt_created);
  END //
DELIMITER ;

DELIMITER //
CREATE FUNCTION created_minute(dt_created VARCHAR(24)) RETURNS datetime(6) AS
  BEGIN
    RETURN date_trunc('minute', str_to_date(dt_created, '%Y-%m-%d %H:%i:%s.%f'));
  END //
DELIMITER ;

DELIMITER //
CREATE FUNCTION as_created_minute(dt_created DATETIME(6)) RETURNS datetime(6) AS
  BEGIN
    RETURN date_trunc('minute', dt_created);
  END //
DELIMITER ;

select names(1,10);

DELIMITER //
CREATE OR REPLACE FUNCTION names(st INT, ed INT) RETURNS VARCHAR(255) AS
  DECLARE
    s VARCHAR(255) = "";
  BEGIN
    FOR i IN st .. ed LOOP
      if i > st THEN
        s = concat(s, ', ');
      END IF;
      s = concat(s, '\'ev_bus_', i, '\'') ;
    END LOOP;
    RETURN s;
  END // 
DELIMITER ;

DELIMITER //
CREATE OR REPLACE PROCEDURE insert_event(st INT,  ed INT) AS
DECLARE
  qry TEXT;
BEGIN
  qry = concat(
    'insert into v_event( event_name, event_id, correlation_id, user_id,  dt_created, dt_received, dt_ingested, labels, payload) ',
    'select event_name, event_id, correlation_id, user_id,  as_date(dt_created), as_date(dt_received), as_date(dt_ingested), labels, payload from event where event_name in (',
      names(st,ed), 
    ')');
  EXECUTE IMMEDIATE qry;
END //
DELIMITER ;

DELIMITER //
CREATE OR REPLACE FUNCTION array_as_string(a ARRAY(VARCHAR(255)))
RETURNS VARCHAR(255) AS
  DECLARE
    result VARCHAR(255) = "";
  BEGIN
    FOR i IN 0 .. LENGTH(a) - 1 LOOP
      result = concat(result, a[i]);
      if length(a) > 1 and i < LENGTH(a) - 1 THEN
        result = concat(result, " ");
      END IF;
    END LOOP;
    RETURN result;
  END //
DELIMITER ;

update event set labels = array_as_string(split(labels, ","));

select array_as_string(split("a1, a2", ","));

select count(*) from event where event_name in (names(90, 100));

call insert_event(91, 100);
select count(*) from v_event;
-- --------------------------------------------------------------------------------
-- --------------------------------------------------------------------------------

