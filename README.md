# Sleep On Lan (SOL)

## Principle

Wake-on-LAN is a standard low-level protocol implemented in various hardware. At this time, there is not standard to make the opposite and send a computer in sleep mode.

This project allows a windows or linux box to be put into sleep from any other device. 

It works with the exact same magic packet than for Wake-On-LAN, the only difference is that the MAC address has to be written in reverse order.

Technically, you have to run a little daemon on your computer that will listen the same Wake-On-LAN port and send the computer in sleep mode when the reversed MAC address received matches a local address. 

Written in `go`, the code may run on linux and windows platforms.

## Usage

Grab the windows + linux release : https://github.com/SR-G/sleep-on-lan/releases/download/1.0.0-SNAPSHOT/SleepOnLAN-1.0.0-SNAPSHOT.zip

### Sleep through UDP

Just send a regular wake-on-lan command but with a reversed MAC address. Thus, the same wake-on-lan tools may be used for both wake and sleep operations (python wake-on-lan script, OpenHab WoL plugin, Android applications, and so on).

Provided you are using a wake-on-lan script like this one [wake-on-lan python script](https://github.com/jpoliv/wakeonlan) (available as a debian package for example), you may use :

<pre>wakeonlan c4:d9:87:7a:78:35 192.168.255.255 // regular mac address, will wake an asleep computer
wakeonlan 35:78:7a:87:d9:c4 192.168.255.255 // reversed mac address, will trigger the UDP listener of the sleep-on-lan process and will thus remotely sleep the computer
</pre>

### Sleep through REST service

If this HTTP listener is activated, the SleepOnLan process then exposes a few REST services, for example :

<pre>http://127.0.0.1:8009/                               // index page, just shows local IP / mac
http://127.0.0.1:8009/sleep                          // remotely sleep this computer through this URL
http://127.0.0.1:8009/wol/c4:d9:87:7a:78:35          // sends a wake-on-lan magic packet on the network to the provided mac address
</pre>

## Configuration

An optional configuration file may be used.

Taken automatically if named "sol.json" and located in the same folder than the SleepOnLan binary.

Content is as follow :

<pre>{
  "Listeners" : ["UDP:9", "HTTP:8009" ],
  "LogLevel" : "INFO",
  "BroadcastIP" : "255.255.255.255",
  "Commands" : []
}
</pre>

**Listeners** defines which mechanism will be activated

- UDP         : will listen on the default port (= 9)
- UDP:&lt;port&gt;  : will listen on the provided port 
- HTTP        : will listen on the default port (= 8009)
- HTTP:&lt;port&gt; : will listen on the provided port

Several listeners may be defined (e.g., "UDP:7", "UDP:9", "HTTP:8009")

If no configuration file is provided, UDP:9 and HTTP:8009 are assumed by default.

The REST services are exposed on 0.0.0.0 and are thus accessibles from http://localhost/, http://127.0.0.1/, http://192.168.1.x/ and so on.

**LogLevel** defines the log level to use. Available values are NONE|OFF, DEBUG, INFO, WARN|WARNING, ERROR. Logs are just written to the stderr/stdout outputs.

**BroadcastIP** defines the broadcast IP used by the /wol service. By default the IP used is 192.168.255.255 (local network range).

**Commands** defines the available commands.

By default, on both windows and linux, only one command is defined : sleep command (through "pm-suspend" on linux and a DLL API call on windows).

You may customize / override this behavior, or add new commands (that will then be available under `http://<IP>:<HTTP PORT>/<operation>` if a HTTP listener is defined), if needed.

Each command has 4 attributes :
- "Operation" : the name of the operation (for the HTTP url)
- "Type" : the type of the operation, either "external" (by default, for remote execution) or "internal-dll" (on windows, to trigger a sleep  through a DLL API call)
- "Default" : true or false. Default command will be executed when UDP magic packets are received. If only one command is defined, it will automatically be the default one
- "Command" : for external commands, the exact command that has to be executed (see examples below). May have to contain full path on windows.

Example 1 : only one (default) operation that will shutdown the system on windows. Through HTTP, the operation will be triggerable with `http://<IP>:<PORT_HTTP>/halt/`.

<pre>
  "Commands" : [ 
    {
        "Operation" : "halt",
        "Command" : "C:\\Windows\\System32\\Shutdown.exe -s -t 0"
    }]
</pre>

Example 2 : force sleep on windows through the rundll32.exe trick (and not through the default API call)

<pre>
  "Commands" : [ 
    {
        "Operation" : "sleep",
        "Command" : "C:\\Windows\\System32\\rundll32.exe powrprof.dll,SetSuspendState 0,1,1"
    }]
</pre>

Example 3 : default operation will put the computer to sleep on linux and a second operation will be published to shutdown the computer through HTTP.

<pre>
  "Commands" : [ 
    {
        "Operation" : "halt",
        "Command" : "pm-halt",
		"Default" : "false"
    },
    {
        "Operation" : "sleep",
        "Command" : "pm-sleep",
		"Default" : "true"
    }]
</pre>

## Installation

### Under windows

The SleepOnLan process may be run manually or, for convenience, installed as a service. The easiest way to install the SleepOnLan service is probably to use [NSSM](https://nssm.cc/) (the Non-Sucking Service Manager).

Usage :

<pre>nssm install &lt;service name&gt; &lt;full path to binary&gt;
</pre>

Installation example :

<pre>c:\Tools\nssm\2.24\win64\nssm.exe install SleepOnLan c:\Tools\SleepOnLan\sol.exe
</pre>

Removal example : 

<pre>c:\Tools\nssm\2.24\win64\nssm.exe remove SleepOnLan confirm
</pre>

Reference : [nssm](https://nssm.cc/usage)

### Under Linux

The SleepOnLan process must use (usually) port 9 (see configuration section if you need another port or if you need to listen to several UDP ports).

Thus the process has either to be ran as root, either has to have the authorization to start on ports < 1024.

The following example allows the process to run on ports &lt; 1024 on recent Linux kernels (for example on ubuntu) :

<pre>sudo setcap 'cap_net_bind_service=+ep' /path/to/sol_binary
nohup /path/to/sol_binary &gt; /var/log/sleep-on-lan.log 2&gt;&1 &
</pre>

You may of course daemonize the process or launch it through an external monitor (like [monit](http://mmonit.com/monit/) or [supervisor](http://supervisord.org/introduction.html)).

## Miscellaneous

### Standalone sleep on lan under windows

Another way to sleep a windows computer remotely :

<pre>net rpc SHUTDOWN -f -I xxx.xxx.xxx.xxx -U uname%psswd
</pre>

### Other similar projects

- [Sleep On Lan](https://github.com/philipnrmn/sleeponlan) (pure java implementation, magic anti-packet starts with 6 * 0x00 instead of 6 * 0xFF)

### OpenHab configuration

Example of configuration under OpenHab.

![OpenHab](sleep-on-lan-openhab.png)

This is a very standard configuration : MAC addresses have just to be reversed.

<pre>
Switch  Network_WoL_Solaris   	"Wake PC (solaris)"   <wake>		(WoL, Status, Network)   { wol="192.168.8.255#14:da:e9:01:98:19" }
Switch  Network_WoL_Jupiter   	"Wake PC (jupiter)"   <wake>		(WoL, Status, Network)   { wol="192.168.8.255#bc:5f:f4:2b:df:26" }
Switch  Network_WoL_Laptop   	"Wake PC (laptop)"    <wake>		(WoL, Status, Network)   { wol="192.168.8.255#C4:D9:87:7A:78:35" }

Switch  Network_SoL_Solaris   	"Sleep PC (solaris)"  <sleep>		(WoL, Status, Network)   { wol="192.168.8.255#19:98:01:e9:da:14" }
Switch  Network_SoL_Laptop   	"Sleep PC (laptop)"   <sleep>		(WoL, Status, Network)   { wol="192.168.8.255#35:78:7A:87:D9:C4" }
</pre>

## Developement

Compile from docker :

```bash
docker run --rm -it -v $(pwd):/go golang /bin/bash
go get -d .../.
go install -ldflags "-d -s -w -X tensin.org/watchthatpage/core.Build=`git rev-parse HEAD`" -a -tags netgo -installsuffix netgo tensin.org/watchthatpage 
GOARCH=amd64 GOOS=windows go install ...

cd /go/src
GOARCH=amd64 GOOS=windows go install main/go/sol/

```