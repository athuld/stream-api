# Steam API
A Golang Gin-Gonic Server to manage stream links from https://github.com/athuld/stream-bot

## Setting Up Server

* Prerequisite
    * Golang (Obviously)
    * MySQL

* Setting up database

  First the MySQL database needs to be setup. Create a database with the name `streamdb` and create a table `streamdata` with following columns

  ```
  +---------------+--------------+------+-----+-------------------+-------------------+
  | Field         | Type         | Null | Key | Default           | Extra             |
  +---------------+--------------+------+-----+-------------------+-------------------+
  | id            | bigint       | NO   | PRI | NULL              | auto_increment    |
  | hash          | varchar(255) | YES  |     | NULL              |                   |
  | filename      | varchar(255) | YES  |     | NULL              |                   |
  | stream_link   | varchar(255) | YES  |     | NULL              |                   |
  | download_link | varchar(255) | YES  |     | NULL              |                   |
  | created_at    | timestamp    | NO   |     | CURRENT_TIMESTAMP | DEFAULT_GENERATED |
  +---------------+--------------+------+-----+-------------------+-------------------+
  ```
* Environmental variables

  Create a `.env` file in the root of the project and add the following variables

  ```
  DB_IP=<ip address of the mysql server: along with port>
  USERNAME=<username for db>
  PASSWORD=<password for db>
  ```

* Running the project

  * Build the project to install dependencies and also generate binary
    ```
    go build
    ```
  * Running the server

    run the generated binary

    ```
    ./streamapi
    ```
    or run the `main.go` file

    ```
    go run main.go
    ```
