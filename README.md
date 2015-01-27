# Sleep On Lan (SOL)

## Principe

## Usage

### Sleep through UDP

Just send a regular wake-on-lan command but with a reversed MAC address. Thus, the same wake-on-lan tools may be used for both wake and sleep operations (python wake-on-lan script, OpenHab WoL plugin, Android applications, and so on).

Provided you are using a wake-on-lan script like this one [wake-on-lan python script](https://github.com/jpoliv/wakeonlan) (available as a debian package for example), you may use :

> wakeonlan c4:d9:87:7a:78:35 192.168.255.255 // regular mac address, will wake an asleep computer
> wakeonlan 35:78:7a:87:d9:c4 192.168.255.255 // reversed mac address, will trigger the UDP listener of the sleep-on-lan process and will thus remotely sleep the computer

### Sleep through REST service

If this listener is activated, the SleepOnLan process then exposes a few REST services, example :

> http://127.0.0.1:8009/                               // index page, just shows local IP / mac
> http://127.0.0.1:8009/sleep                          // remotely sleep this computer
> http://127.0.0.1:8009/wol/c4:d9:87:7a:78:35          // sends a wake-on-lan magic packet on the network to the provided mac address

## configuration

An optional configuration file may be used.

Taken automatically if named "sol.json" and located in the same folder than the SleepOnLan binary.

Content is as follow :

> {
>   "Listeners" : ["UDP:9", "HTTP:8009" ],
>   "LogLevel" : "INFO",
>	"BroadcastIP", "255.255.255.255"
> }

**Listeners** defines which mechanism will be activated
UDP         : will listen on the default port (= 9)
UDP:<port>  : will listen on the provided port 
HTTP        : will listen on the default port (= 8009)
HTTP:<port> : will listen on the provided port

Several listeners may be defined (e.g., "UDP:7", "UDP:9", "HTTP:8009")

If no configuration file is provided, UDP:9 and HTTP:8009 are assumed by default.

The REST services are exposed on 0.0.0.0 and are thus accessibles from http://localhost/, http://127.0.0.1/, http://192.168.1.x/ and so on.

**LogLevel** defines the log level to use. Available values are NONE|OFF, DEBUG, INFO, WARN|WARNING, ERROR. Logs are just written to the stderr/stdout outputs.

**BroadcastIP** defines the broadcast IP used by the /wol service. By default the IP used is 192.168.255.255 (local network range).

## Installation

### Under windows

The SleepOnLan process may be run manually or, for convenience, installed as a service. The easiest way to install the SleepOnLan service is probably to use [NSSM](https://nssm.cc/) (the Non-Sucking Service Manager).

Usage :

> nssm install <service name> <full path to binary>

Installation example :

> c:\Tools\nssm\2.24\win64\nssm.exe install SleepOnLan c:\Tools\SleepOnLan\sol.exe

Removal example : 

> c:\Tools\nssm\2.24\win64\nssm.exe remove SleepOnLan confirm

Reference : [nssm](https://nssm.cc/usage)

### Under Linux

The SleepOnLan process must use (usually) port 9 (see configuration section if you need another port or if you need to listen to several UDP ports).

Thus the process has either to be ran as root, either has to have the authorization to start on ports < 1024.

The following example allows the process to run on ports < 1024 on recent Linux kernels (for example on ubuntu) :

> sudo setcap 'cap_net_bind_service=+ep' /path/to/sol_binary
> nohup /path/to/sol_binary > /var/log/sleep-on-lan.log 2>&1 &

You may of course daemonize the process (see []()) or launch it through an external monitor (like [monit](http://mmonit.com/monit/) or [supervisor](http://supervisord.org/introduction.html))

## Misc

Other way to sleep a windows computer remotely :

> net rpc SHUTDOWN -f -I xxx.xxx.xxx.xxx -U uname%psswd

Other similar projects :

- https://github.com/philipnrmn/sleeponlan (pure java implementation, magic anti-packet starts with 6 * 0x00 instead of 6 * 0xFF)
- 
