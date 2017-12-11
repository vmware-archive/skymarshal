
# security

Concourse security is responsible for issuing tokens that will be consumed by the `atc`. They should be statically verifiable and contain all required openid connect claims, as well as a `groups` claim for verifying concourse team membership. 

This library is a thin wrapper around [coreos/dex](http://github.com/coreos/dex).

### Future considerations

We want to investigate the ability to grant personal access tokens.
