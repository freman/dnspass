# dnsPass

Automagically bypass Australian internet filters

This is crudely hacked together and comes with no warranty - it's also my first windows service

# Installing

## Windows

You'll need 2 files.

 * dnspass.exe
 * config.toml

You'll need to edit config.toml to include your ISP's dns servers under Untrust look [here](https://www.whatsmydns.net/dns/australia) if you don't know then

Place these two files in a directory together, start cmd as administrator

```
dnspass install
dnspass start
```

edit your network adapters settings to use 127.0.0.1 as your primary dns

## Linux

You'll need 2 files.

 * dnspass
 * config.toml

You'll need to edit config.toml to include your ISP's dns servers under Untrust look [here](https://www.whatsmydns.net/dns/australia) if you don't know then

Place these two files into a directory together

```
sudo dnspass
```

edit your network settings (probably just vim /etc/resolv.conf) to have nameserver 127.0.0.1