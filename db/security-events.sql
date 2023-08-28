
CREATE ROLE 'db_events_owner_role';
GRANT CREATE, ALTER, DROP on events.* to ROLE 'db_events_owner_role';
GRANT CREATE VIEW, ALTER VIEW, DROP VIEW on events.* to ROLE 'db_events_owner_role';
GRANT SHOW VIEW on events.* to ROLE 'db_events_owner_role';
GRANT CREATE TEMPORARY TABLES on events.* to ROLE 'db_events_owner_role';
GRANT SELECT, INSERT, UPDATE, DELETE on events.* to ROLE 'db_events_owner_role';
GRANT LOCK TABLES on events.* to ROLE 'db_events_owner_role';

CREATE GROUP 'admin_events_group';
GRANT ROLE 'db_events_owner_role' to 'admin_events_group';


CREATE ROLE 'db_events_analytics_role';
GRANT SELECT, INSERT, UPDATE, DELETE on events.* to ROLE 'db_events_analytics_role';
GRANT SHOW VIEW on events.* to ROLE 'db_events_analytics_role';
GRANT LOCK TABLES on events.* to ROLE 'db_events_analytics_role';
GRANT CREATE TEMPORARY TABLES on events.* to ROLE 'db_events_analytics_role';


CREATE GROUP 'admin_events_group';
GRANT ROLE 'db_events_owner_role' to 'admin_events_group';