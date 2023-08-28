GRANT USAGE ON *.* TO 'adm_francara' IDENTIFIED BY 'singlestore';
GRANT GROUP 'admin_ss_group' TO 'adm_francara';
GRANT GROUP 'admin_events_group' TO 'adm_francara';

GRANT USAGE ON *.* TO 'adm_ronie' IDENTIFIED BY '12345';
GRANT GROUP 'admin_events_group' TO 'adm_ronie';

GRANT USAGE ON *.* TO 'zuki' IDENTIFIED BY '12345';
GRANT GROUP 'admin_mlops_group' TO 'zuki';
GRANT GROUP 'analytics_events_group' TO 'zuki';
GRANT GROUP 'ingest_events_group' TO 'zuki';

GRANT USAGE ON *.* TO 'nicole.oliveira' IDENTIFIED BY '12345';
GRANT GROUP 'analytics_events_group' TO 'nicole.oliveira';
GRANT GROUP 'analytics_mlops_group' TO 'nicole.oliveira';

CREATE USER ingest_events WITH DEFAULT RESOURCE POOL = ingest;
ALTER USER ingest_events IDENTIFIED BY '12345'
GRANT GROUP 'ingest_events_group' TO 'ingest_events';

CREATE USER rt_pp1min WITH DEFAULT RESOURCE POOL = rt;
ALTER USER rt_pp1min IDENTIFIED BY '12345'
GRANT GROUP 'analytics_events_group' TO 'rt_pp1min';

CREATE USER nrt_pp1min WITH DEFAULT RESOURCE POOL = nrt;
ALTER USER nrt_pp1min IDENTIFIED BY '12345'
GRANT GROUP 'analytics_events_group' TO 'nrt_pp1min';
