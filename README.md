gotagmee
========

This is a small project make for filling some data to a Neo4J server from the
meetup.com API.

Quick example
-------------

If you want to quickly run this:

    go run main.go \
        -neo4j="http://user:pass@host:port/db/data/" \
        meetup_api_key \
        group_id # (ex: go-london-user-group)

Disclaimer
----------

Do you know TDD? I knew it as well! But sadly I forgot for this project :(

It's just a proof of concept and all the important testing should be mocking
http request which is tedious and not really important for what I wanted: **the
data**.

Of course, you can always use this under your own responsibility, and if your
think that it's somehow useful, go ahead and send me some pull request to make
it really reusable.
