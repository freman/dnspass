# dnsPass

Automagically bypass Australian internet filters

This comes with no warranty - it's also my first windows service

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

# Configuration file

The config file is called config.toml

It has few options, I'll go through them here

## Listen _"localhost:53"_

Probably best to keep this on localhost:53 unless you know what you're doing, this application has no built in security so keeping
it configured to only listen on localhost stops other people (both in your network, and outside) from getting up to no good

## AutoUpdatePoisonHosts _false_

Set this to true if you want to pull the latest version of poisoned hosts when the app starts and every 24 hours after that.
This is completely optional, the [source of the list](data/) is quite public so you can chose to manually update your config.toml file

## Untrust _[]_

This is a list of DNS servers you don't trust, usually your ISPs - in fact I have it by default set up for TPG [here](https://www.whatsmydns.net/dns/australia) is a list of DNS servers by ISP

## Trust _[]_

This is the list of DNS servers that you trust (if not completely, than more than your ISP) - I set them up for google

## PoisonHosts _[]_

This is the list of hosts that ISPs are known to poison answers with, this list is ignored if you use the AutoUpdatePoisonHosts option

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