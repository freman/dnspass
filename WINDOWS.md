# More detailed windows instructions

## Important note

The config.toml that comes in the zip is pre-configured for TPG because that's what I have to test with, you might have to edit it for your ISP.
I recommend notepad++ or wordpad (because it's probably not got the right line endings to just work in notepad.exe)

## Windows 7

One of my first users who patiently waited for me to modify the travis build script to *actually* build the windows binaries offered up [these instructions](https://forums.whirlpool.net.au/forum-replies.cfm?t=2591386&p=53#r1045).

For those interested in running it on Windows 7, the best way is to do this: (I dont have Windows 8 or Windows 10 so I dont know all the steps there)

1. Copy the files from the zip into c:\windows\system32 (if you are on a 64 bit version of Windows, use the 64 bit zip, otherwise use the 32 bit zip)
2. Go to c:\windows\system32 in explorer, find cmd.exe, right click on it and select "run as administrator" and a command window will appear.
3. Type dnspass install then dnspass start to install and start the service.
4. Type services.msc to bring up the services interface
5. Find "my service" and right click on it and select "properties". It should say "c:\windows\system32\dnspass.exe" as the path to the executable if you have done it right.
6. Under "startup type" select "automatic" which will start the service along with Windows.
7. To change your DNS to point to the new server, follow these instructions from Google:
https://developers.google.com/speed/public-dns/docs/using#windows
except you want to use 127.0.0.1 and 127.0.0.0 as the DNS server addresses instead of the Google addresses. You want to change the settings for Internet Protocol Version 4 only.

This app is better than using Google directly since all the CDNs and other things work properly and pick up the correct IP addresses.

_editors note: you don't need a second dns server, but if you use one please use the same one 127.0.0.1_

## Windows 10

While probably not the right way to do it, for testing I simply followed these steps

1. Copy the files from the zip into c:\meh (if you are on a 64 bit version of Windows, use the 64 bit zip, otherwise use the 32 bit zip)
2. Press the windows key, type cmd, right click on the cmd.exe that shows up in the start menu, and click "Run as administrator"
3. Run dnspass install
4. Run dnspass start
