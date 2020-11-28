# wikinodes-server
web-server for the wikinodes project (API)
<br>

This repo contains a simple server which is mean to support the 'wikinodes' project. <br>
It (the server) acts as middleware between an app and a backing database.
<br>
App & Database population repos are found at:
  * [wikinodes-app](https://github.com/crunchypi/wikinodes-app) (front-end)
  * [wikinodes-preprocessing](https://github.com/crunchypi/wikinodes-preprocessing) (creating/populating db)

<br>
<br>

### Usage

#### NOTE: The current working branch is 'develop'.

To use the server, first start and/or populate a neo4j database with [wikinodes-preprocessing](https://github.com/crunchypi/wikinodes-preprocessing). Population can be done with anything else, of course, but that will likely require some additional work because the names of node labels and fields are hardcoded [here](https://github.com/crunchypi/wikinodes-server/blob/develop/db/neo4j/neo4j.go), as is necessary for this server to do its job.

<br><br>

When the Neo4j db is up, simply start the server with the line below (uri, usr and pwd is for neo4j)
```
sudo go run main.go <uri> <usr> <pwd>
```

<br> <br>
The API is simple and can do only two things, as shown below. Both endpoints expect a string that specifies the title of a wikipedia node and return an appropriate json. Note, at the moment, the IP & port is static (localhost:42040) and there is no TLS. This will be fixed soon (by 2021).

* Fetch a Wiki Node:
  * curl: ```curl localhost:42040/data/read -d someWikiNodeTitle```
  * resp: ```[{"ID":569,"Title":"someWikiNodeTitle", "HTML":"..."}]```
  
* Fetch a Wiki Node, but only the bare minimum of data:
  * curl: ```curl localhost:42040/data/read-brief -d someWikiNodeTitle```
  * resp: ```[{"ID":569,"Title":"someWikiNodeTitle"}]```
