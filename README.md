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

### Usage

#### NOTE: The current working branch is 'develop'.

To use the server, first start and/or populate a neo4j database with [wikinodes-preprocessing](https://github.com/crunchypi/wikinodes-preprocessing). Population can be done with anything else, of course, but that will likely require some additional work because the names of node labels and fields are hardcoded [here](https://github.com/crunchypi/wikinodes-server/blob/develop/db/neo4j/nodeconstruct.go), as is necessary for this server to do its job.

<br>

When the Neo4j db is up, simply start the server with the line below (uri, usr and pwd is for neo4j). Note, at the moment, the IP & port is static (localhost:1234) and there is no TLS. This will be fixed soon (by 2021).
```
sudo go run main.go <uri> <usr> <pwd> # Starts listening at localhost:1234
```

<br>

### API

The API is simple and has three endpoints, each of them listed below (more information further down).
* (**1**) ```ip:port/data/search/node ```  (search wiki data)
* (**2**) ```ip:port/data/search/neigh```  (search neighbours/related articles of a specific article)
* (**3**) ```ip:port/data/search/rand ```  (search random nodes)

<br>

(**1**) This endpoint accepts a json with format ```{"title":string,"brief":bool}```, where "title" represents a wiki article title, while "brief" specifies how much data to request (true gives id&title, false gives id&title&html).
* Example 1:
  * curl: ```localhost:1234/data/search/node -d "{\"title\":\"Art\",\"brief\":true}"```
  * resp: ```[{"id":306,"title":"Art"}, ... ]```
* Example 2:
  * curl: ```localhost:1234/data/search/node -d "{\"title\":\"Art\",\"brief\":false}"```
  * resp: ```[{"id":306,"title":"Art","html":"<some long html string>"}, ... ]```
  
<br>
  
(**2**) This endpoint also accepts a json, though with a bit different options: ```{"title":string,"exclude":[string], "limit":int}```, where "title" represents a wiki article title, "exclude" is a list of article names that are not wanted (these are filtered from the response), while "limit" specifies the max amount of objects that are wanted in the response. Do note that objects included in the response are always "brief", similar to **1, example 1**.
* Example 1:
  * curl: ```ip:port/data/search/neigh -d "{\"title\":\"Art\",\"exclude\":[\"Art history\"],\"limit\":2}"```
  * resp: ```[{"id":2,"title":"Abstract animation"}{...not 'Art history'...}]```
* Example 2:
  * curl: ```ip:port/data/search/neigh -d "{\"title\":\"Art\",\"exclude\":[],\"limit\":3}"```
  * resp: ```[{"id":2,"title":"Abstract animation"},{...},{...}]```

<br>

(**3**) This endpoint accepts a json as well, with format ```{"amount":int}```, where "amount" specifies how many random nodes/objects to get. As with the previous endpoint, the return data is also "brief".
* Example 1:
  * curl: ```ip:port/data/search/rand -d "{\"amount\":1}```
  * resp: ```[{"id":10,"title":"Albert Gleizes"}]```
* Example 2:
  * curl: ```ip:port/data/search/rand -d "{\"amount\":3}```
  * resp: ```[{"id":7,"title":"Abstraction"},{...},{...}]```

