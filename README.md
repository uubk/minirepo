# Minirepo
A small tool and library to quickly set up binary repositories for your applications.

[![Build Status](https://travis-ci.com/uubk/minirepo.svg?branch=master)](https://travis-ci.com/uubk/minirepo)
[![Go Report Card](https://goreportcard.com/badge/github.com/uubk/minirepo?style=flat)](https://goreportcard.com/report/github.com/uubk/minirepo)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](http://godoc.org/github.com/uubk/minirepo)
[![Release](https://img.shields.io/github/release/uubk/minirepo.svg?style=flat)](https://github.com/uubk/minirepo/releases/latest)

## Motivation
Minirepo assumes that you want to serve something like the following structure to your application:
```
repo/
repo/platform
repo/platform/component
repo/platform/component/version
repo/platform/component/version/file
repo/meta.yml
repo/meta.asc
```
In the above example, the repository root contains files sorted by platform, component
and version, although this is completely left up to you. Minrepo will then generate an index (meta.yml) which contains sha384sums
of all files and then sign this metadata. To reproduce the structure from this example,
simply use
`minirepo -repo <PATH>`

A metadata file for example can look like this:
```
contents:
- name: foo
  children:
  - name: bar
    children:
    - name: test
      hash: e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
name: test
timestamp: 2018-08-16T13:05:04.343849113+02:00
```

To use this repository in a client, you would then do:
```
client := NewRepoClient("/tmp", "http://127.0.0.1:8080", "<content of pub.asc>")
file, err := client.GetFile("foo", "bar", "test")
```

### Debugging hints
To verify the detached signature manually when using a new-ish GPG release, you'll need
to create a keyring with the public key:
```
# gpg -v --no-default-keyring --keyring $(pwd)/pub.gpg --import pub.asc
```
You can then use
```
# gpg -v --no-default-keyring --keyring $(pwd)/pub.gpg --verify meta.asc meta.yml
```
to check the signature. The output should look like this:
```
gpg: Signature made Don 16 Aug 2018 13:05:04 CEST
gpg:                using RSA key 0xC17FEE56EF1ED674
gpg: using pgp trust model
gpg: Good signature from "minirepo (Autogenerated)" [unknown]
gpg: WARNING: This key is not certified with a trusted signature!
gpg:          There is no indication that the signature belongs to the owner.
Primary key fingerprint: A6EB D286 E13E 9DE0 FA49  4BEA C17F EE56 EF1E D674
gpg: binary signature, digest algorithm SHA256, key algorithm rsa2048
```
The warning is expected as we didn't trust the key.

## License
Apache 2.0