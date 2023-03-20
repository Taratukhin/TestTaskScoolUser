Start MySQL server: `docker-compose up -d`

Manual connect to MySQL server is: `mysql -h localhost -P 5439 -u username -puserpassword -D schooldb --protocol=tcp`
Please wait 1-2 seconds for start container.

Stop MySQL server: `docker-compose down -v --rmi all`
	option `-v` delete all volumes
	option `-rmi all` delete image
