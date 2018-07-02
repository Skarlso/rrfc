[![Build Status](https://travis-ci.org/Skarlso/rrfc.svg?branch=master)](https://travis-ci.org/Skarlso/rrfc)

# Random RFC

What is Random RFC? Random RFC will provide a randomly selected RFC to work on and discuss about with others.

The workflow is as follows:

* Download a list of available RFCS from IETF
* Select one randomly
* Display it and add a disqus link for discussion about that RFC
* View all previously appeared RFCS in a list if you are interested in a particular one

That's about it

# Backend

The backend is powered by Go. It downloads the list and stores it in a postgres database for later parsing.

The randomly selected rfc is then used to generate the index file for the project which is a simple static html file.

# Frontend

The front-end is a lightweight static HTML app. All files are generated using Go html template files. The disqus links are the only one that are dynamic. Each rfc should have its own disqus link.

# Deplyoment

RRFC is deployed via [docker stack](https://docs.docker.com/get-started/part5/). You can find the deployment yaml [here](docker-cloud.yml).

To deploy the stack use:

```bash
docker stack deploy -c docker-cloud.yml rrfc
```

The stack will create a postgres database and the rrfc application by running rrfc once. That will initialize the database, create the necessary tables and fill it up with the current available rfcs and generate the index and previous html files.

It's worth noting that postgres will do a bind mount outside into `/data` in order to persist it's database.

Also, the rrfc app is using HTTPS which requires letsencrypt in the host to be present for now under `/etc/letsencrypt`. Because for obvious reasons the docker image does not have the certificates included.

# Contributions

Are welcomed and appreciated!
