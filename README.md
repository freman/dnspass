# dnsPass

Automagically bypass Australian internet filters

This is crudely hacked together and comes with no warranty - it's also my first windows service

# How it works

Magic

# Seriously, how it works?

Well, you send it your dns queries, it forwards them on to your ISP `Untrust` - this means

 * you get all your regional goodies
 * you get your your isp freebies

However, if your ISP replies with an address we know is a poisoned address then
this application simply goes and asks google `Trust` for it's opinion.

It also remembers that this has happened for a while so it doesn't waste time asking your
ISP again and just goes directly to google.

## Wtf _Trust_ and _Untrust_?

Well, _Untrust_ means we can't trust the results coming from those servers - they've broken our trust and we know they lie to us.
_Trust_ means that we trust those servers more so we'll belive them for a second opinion.

You don't have to use google's servers as your trust, you can put any servers in there you trust, opendns for example.

# Installing

## Windows

([More detailed instructions](WINDOWS.md))

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

### Uninstalling

edit your network adapters settings to use your router, or your isp settings (or just flick the dns setting back to dhcp)

start cmd as administrator

```
dnspass stop
dnspass remove
```

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