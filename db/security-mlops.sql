
CREATE ROLE 'db_mlops_owner_role';
GRANT CREATE, ALTER, DROP on mlops.* to ROLE 'db_mlops_owner_role';
GRANT CREATE VIEW, ALTER VIEW, DROP VIEW on mlops.* to ROLE 'db_mlops_owner_role';
GRANT SHOW VIEW on mlops.* to ROLE 'db_mlops_owner_role';
GRANT CREATE TEMPORARY TABLES on mlops.* to ROLE 'db_mlops_owner_role';
GRANT SELECT, INSERT, UPDATE, DELETE on mlops.* to ROLE 'db_mlops_owner_role';
GRANT LOCK TABLES on mlops.* to ROLE 'db_mlops_owner_role';

CREATE GROUP 'admin_mlops_group';
GRANT ROLE 'db_mlops_owner_role' to 'admin_mlops_group';
GRANT ROLE 'monitoring_role' to 'admin_mlops_group';


CREATE ROLE 'db_mlops_analytics_role';
GRANT SELECT on mlops.* to ROLE 'db_mlops_analytics_role';
GRANT SHOW VIEW on mlops.* to ROLE 'db_mlops_analytics_role';
GRANT LOCK TABLES on mlops.* to ROLE 'db_mlops_analytics_role';
GRANT CREATE TEMPORARY TABLES on mlops.* to ROLE 'db_mlops_analytics_role';

CREATE GROUP 'analytics_mlops_group';
GRANT ROLE 'db_mlops_analytics_role' to 'analytics_mlops_group';

CREATE ROLE 'db_mlops_ingest_role';
GRANT INSERT, UPDATE, DELETE, SELECT on mlops.* to ROLE 'db_mlops_ingest_role';
GRANT SHOW VIEW on mlops.* to ROLE 'db_mlops_ingest_role';
GRANT LOCK TABLES on mlops.* to ROLE 'db_mlops_ingest_role';
GRANT CREATE TEMPORARY TABLES on mlops.* to ROLE 'db_mlops_ingest_role';


CREATE GROUP 'ingest_mlops_group';
GRANT ROLE 'db_mlops_ingest_role' to 'ingest_mlops_group';
