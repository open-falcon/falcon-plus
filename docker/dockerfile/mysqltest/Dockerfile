FROM mysql:5.7

COPY ./scripts/mysql/db_schema/ /db_schema/
COPY ./modules/api/test_utils/sqltestset/ /sqltestset/
RUN cp /db_schema/*.sql /docker-entrypoint-initdb.d/
RUN cp /sqltestset/*.sql /docker-entrypoint-initdb.d/
