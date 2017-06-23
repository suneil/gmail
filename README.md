# Gmail

Downloads email from gmail account and stores it in mongo for further analysis by other tools

This will work but is not production ready as it is missing some features like resuming from last email


## Install

```
go get golang.org/x/net/context
go get golang.org/x/oauth2/google
go get google.golang.org/api/gmail/v1
go get gopkg.in/mgo.v2
go get github.com/juju/ratelimit
go get github.com/spf13/cobra
go get github.com/spf13/viper

git clone https://github.com/suneil/gmail
```

## License

See LICENSE for details
