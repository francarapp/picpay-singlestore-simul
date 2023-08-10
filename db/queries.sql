select format(count(*), 0) from event;
select * from event limit 10;
select event_name, format(count(*), 0) from event
group by event_name;

select created_minute(dt_created), event_name, count(*) from event 
where 
  created_hour(dt_created) between '2023-08-10 01:00:00' and '2023-08-10 12:00:00'
  and event_name  in ("ev_bus_498", "ev_bus_499", "ev_bus_500", "ev_bus_501", "ev_bus_502")
  and  match(labels) against ('bus5')
group by created_minute(dt_created), event_name;

select created_minute(dt_created), event_name, sum(payload::value) from event 
where 
  created_hour(dt_created) between '2023-08-10 01:00:00' and '2023-08-10 12:00:00'
  and event_name  in ("ev_bus_498", "ev_bus_499", "ev_bus_500", "ev_bus_501", "ev_bus_502")
  and  match(labels) against ('bus5')
  and payload::key_10 is not null
group by created_minute(dt_created), event_name;

select user_id, event_name, format(count(*), 0) from event
where
  event_name  in ("ev_bus_498", "ev_bus_499", "ev_bus_500", "ev_bus_501", "ev_bus_502")
group by user_id, event_name;

select event_name, format(sum(payload::value), 2) from event
group by event_name;

select format(count(*), 0) from event
where payload::$key_40 is not null;

select format(count(*), 0) from event
where MATCH(labels) AGAINST ('bus9');

select event_name, format(sum(payload::value), 2) from event
where MATCH(labels) AGAINST ('bus9')
group by event_name;