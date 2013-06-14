REST Structure:

GET /machines											List of all machines
GET /machines/:Id									Details of one machine

GET /machines/:Id/setup						Detail on running setup procedure for one machine, if any. 404 if none

POST /machines/:Id/setup					Add setup procedure for one machine (identified by :Id in post data)

