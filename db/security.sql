CREATE ROLE 'security_role';
GRANT USAGE on *.* to ROLE 'security_role' WITH GRANT OPTION;
GRANT CREATE USER on *.* to ROLE 'security_role';

CREATE GROUP 'security_ss_group';
GRANT ROLE 'security_role' to 'security_ss_group';

CREATE ROLE 'dba_role';
GRANT CREATE DATABASE, DROP DATABASE on *.* to ROLE 'dba_role';
GRANT RELOAD on *.* to ROLE 'dba_role';
GRANT SUPER on *.* to ROLE 'dba_role';
GRANT CLUSTER on *.* to ROLE 'dba_role';
GRANT SHOW METADATA on *.* to ROLE 'dba_role';
GRANT BACKUP, RELOAD on *.* to ROLE 'dba_role';

CREATE GROUP 'admin_ss_group';
GRANT ROLE 'dba_role' to 'admin_ss_group';
GRANT ROLE 'security_role' to 'admin_ss_group';

