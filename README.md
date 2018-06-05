# rrfc

Random RFC

`pg_ctl -D /usr/local/var/postgres start`

# Running the container

```bash
docker run -itd --name rrfc -e PGDATA=/data \
                            -e POSTGRES_PASSWORD=password123 \
                            -e POSTGRES_USER=rrfc \
                            -e POSTGRES_DB=rfcs \
                            -v `pwd`:/data \
                            -p 80:80 \
                            skarlso/rrfc /bin/bash
```