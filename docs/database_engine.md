
# Database Engine

```json
{"password": "mysecretpassword", "username": "postgres", "role_template": "CREATE ROLE {{name}} WITH LOGIN PASSWORD '{{password}}'; GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO {{name}}; ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO {{name}};", "connection_string": "postgresql://postgres:mysecretpassword@localhost:5432/cryptkeeper?sslmode=disable"}
```


CREATE ROLE miriam WITH LOGIN PASSWORD 'jw8s0F4' VALID UNTIL '2005-01-01';
