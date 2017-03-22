# How to run those test script

1. Change session or Insert valid session
  * `Apitoken := {"name": "root", "sig": "233fdb00f99811e68a5c001500c6ca5a"}`
2. cat alarms.sql | mysql -u user -p [password]
3. go build -> build single binary file `api`
4. configure `cfg.json` & run `./api`
5. cd module/api/test/api/alarm
6. `go test`
