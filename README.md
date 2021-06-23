# Installation
## Install redis
###
    brew install redis
    redis-server /usr/local/etc/redis.conf
### Ubuntu
    sudo apt-get install redis-server
    systemctl start redis-server
### Other
    Following https://redis.io/topics/quickstart

## Run cache-server
```git checkout https://github.com/johnzhanghua/cache-server```

```cd cache-server && go build ./cmd/cache-server/...```

```REDISSRV=127.0.0.1:6379 REDISUSER="" REDISPASSWORD="" ./cache-server```

## Test

### Unit test
    go test ./...

### Integration test
```curl -XPOST http://127.0.0.1:8080/v1/contact -d @testdata/contact.json```

```{"contact_id":"person_87F8C32A-B357-4C6C-B683-1CD0A344DC46","Email":""}```

```curl -XGET http://127.0.0.1:8080/v1/contact/person_87F8C32A-B357-4C6C-B683-1CD0A344DC46```

```"contact_id":"person_87F8C32A-B357-4C6C-B683-1CD0A344DC46","Email":"chris@autopilothq.com","FirstName":"Chris","LastName":"Sharkey","Company":"Magpie API","type":"Contact","Phone":"4159945916","LeadSource":"Autopilot","Status":"Testing"}```