# wikinodes-server
web-server/api for the wikinodes project.
<br>

This repo contains a server which is mean to support the 'wikinodes' project. <br>
It (the server) essentially acts as middleware between an app and a backing database.
<br>
App & Database population repos are found at:
  * [wikinodes-app](https://github.com/crunchypi/wikinodes-app) (front-end)
  * [wikinodes-preprocessing](https://github.com/crunchypi/wikinodes-preprocessing) (creating/populating db)

<br>

### Usage

Dependencies:
- Neo4j (used for content, specifically Wikipedia articles)
- Redis (caching of IPs for spam mitigation & a feature related to article recommendation)

To use the server, first start and/or populate a neo4j database with [wikinodes-preprocessing](https://github.com/crunchypi/wikinodes-preprocessing). Population can be done with anything else, of course, but that will likely require some additional work because the db schema might not match up  (labels, field, etc). Still, that can be found [here](https://github.com/crunchypi/wikinodes-preprocessing/blob/master/src/typehelpers.py). Redis has to be started as well.

<br>

When the Neo4j & Redis services are up, configure the server in root/config/config.go, the most important things to look for are uri,usr&pwd for databases and perhaps the ip+port for this server. 

<br>

### API

The API has 7 endpoints, all of which are JSON over POST. They're all read-only in the sense that you can't directly change any data
but the 4th one below (../byneigh) is used with Redis to cache searches and use that data to update a relationship weight between
linked articles in Neo4j for the purpose of article recommendation.

- [```ip:port/data/search/articles/byid```](#ipportdatasearcharticlesbyid)
- [```ip:port/data/search/articles/bytitle```](#ipportdatasearcharticlesbytitle)
- [```ip:port/data/search/articles/bycontent```](#ipportdatasearcharticlesbycontent)
- [```ip:port/data/search/articles/byneigh```](#ipportdatasearcharticlesbyneigh)
- [```ip:port/data/html/byid```](#ipportdatahtmlbyid)
- [```ip:port/data/check/relsexist```](#ipportdatacheckrelsexist)
- [```ip:port/data/random/articles```](#ipportdatarandomarticles)

----
#### ip:port/data/search/articles/byid
This endpoint searches the data layer for Wikipedia content (article(s)) by article ID and accepts a JSON of form `{id:int}`.
<br>
curl(v7.68.0) example:
```
curl http://ip:port/data/search/articles/byid -d "{\"id\":4279}"
# Return might be [{"id":4279,"title":"1853"}] if that article exists.
```
----
#### ip:port/data/search/articles/bytitle
This endpoint searches the data layer for Wikipedia content (article(s)) by title and accepts a JSON of form `{title:string}`
<br>
curl(v7.68.0) example:
```
curl http://ip:port/data/search/articles/bytitle -d "{\"title\":\"Last Thursdayism\"}"
# Return might be [{"id":5,"title":"Last Thursdayism"}] if that article exists.
```
----
#### ip:port/data/search/articles/bycontent
This endpoint searches the data layer for Wikipedia content (article(s)) by looking through the body/content, using
a JSON of form `{str:string, limit:int}`. Note this relies on DB indexing for performance but should not
be a problem if the DB is populated with [wikinodes-preprocessing](https://github.com/crunchypi/wikinodes-preprocessing).
<br>
curl(v7.68.0) example:
```
curl http://ip:port/data/search/articles/bycontent -d "{\"str\":\"the\", \"limit\":1}"
# Return might be [{"id":8,"title":"Last Thursdayism"}] if the search is ok.
```
----
#### ip:port/data/search/articles/byneigh
This endpoint searches the data layer for Wikipedia content (article(s)) for neighbours of a given article id (i.e
articles hyperlinked from the article with the provided ID), using a JSON of form `{id:int, limit:int}`. **Note**,
this endpoint isn't deterministic, it uses a markov-chain-like recommendation (relying on Redis).
<br>
curl(v7.68.0) example:
```
curl http://ip:port/data/search/articles/byneigh -d "{\"id\":4394, \"limit\":1}"
# Return might be [{"id":8,"title":"Last Thursdayism"}] if the relationship is true.
```
----
#### ip:port/data/html/byid
This endpoint searches the data layer for the *HTML* of a Wikipedia content with a article given ID, using a
JSON of form `{id:int}`
<br>
curl(v7.68.0) example:
```
curl http://ip:port/data/search/html/byid -d "{\"id\":8}"
# Might return a HTML string if that article exists.
```
----
#### ip:port/data/check/relsexist
This endpoint checks the data layer for whether or not relationships exist between articles, using a JSON
where the key is 'rels' and value is expected to be a nested list, where the inner ones are of length 2, like
`{rels:[[4394, 4395]]}`. This checks if there is an article with ID 4395 in the database that connects to another
article with ID 4395 (order matters).
<br>
curl(v7.68.0) example:
```
curl http://ip:port/data/check/relsexist -d "{\"rels\":[[4394, 4395]]}"
# Returns [true] if that relationship exists.
```
----
#### ip:port/data/random/articles
This endpoint searches the data layer for a specified amount of random articles, accepting a JOSN of form `{limit:int}`.
<br>
curl(v7.68.0) example:
```
curl http://ip:port/data/random/articles -d "{\"limit\":1}"
# Might return [{"id":9,"title":"2010"}]
```

