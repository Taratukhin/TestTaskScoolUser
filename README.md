#  Test task: Users - schoolchildren
## Task description
1. Write service. Stack: `Go, MySQL`
2. Starting a service with a port parameter. (port)
   - endpoints: `GET /profile`
   - headers: Content-type: application/json
3. Create a test base and import the structure and data.
   - `scheme.sql` - structure
   - `data.sql` - data
4. Add constraints to tables.
5. Write an authentication middleware, checking the header `Api-key` in the auth table. In case of incorrect `Api-key`, error 403.
6. `GET /profile` - give the data of all users.
    - If the `username` parameter is present, return one object.
7. At the output, the object should contain: `id, username, first_name, last_name, city, school`, taken from the tables `user, user_profile, user_data`.
## HOWTO 
### Compile
Run `make build`
### Run
Run `./schoolserver -h` to show the manual. You will see:
```
Usage of ./schoolserver:
   -db string
         Database connect URL (default "username:userpassword@tcp(127.0.0.1:5439)/schooldb")
   -port int
         Port of HTTP server 8080 (default 8080)
   -sql value
         List of SQL files to execute, default: --sql ./sql/create.sql --sql ./sql/scheme.sql --sql ./sql/truncate.sql --sql ./sql/data.sql --sql ./sql/constraints.sql
```
### Test database function and api
The tests use MySQL. Please start MySQL early. After that run `make test`
### Start MySQL
`cd mysql-db` and `docker-compose up -d`
`cd ..` and run `./schoolserver`
After tests run `docker-compose down -v --rmi all` it will erase volume and image 

### Work
- Run `curl -i -H "Api-key: www-dfq92-sqfwf" -X GET http://localhost:8080/profile` Result must be `200 OK`
- Run `curl -i -H "Api-key: www-dfq92-sqfwf" -X GET http://localhost:8080/profile?username=test` Result must be `200 OK`
- Run `curl -i -H "Api-key: www-dfq92-sqfwf" -X GET http://localhost:8080/profile?username=testtttt` Result must be `400 Badrequest`
- Run `curl -i -H "Api-key: wrong-key" -X GET http://localhost:8080/profile?username=test` Result must be `403 Forbidden`
