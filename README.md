# netsa-stuff
Example configs and tools for NetSA toolset

* pipeline.conf             - Example config for pipeline which detects common attack traffic
* rwsender.sh               - Helper script for managing large regexps within rwsender (play lists of probes in /etc/silk/rwsender/{receivername})
* silk-v7.patch             - Adds support for devices exporting Netflow v7 (although ideally configure them to use v5)
* silk-flowcapformat.patch  - Modifies flowcap so it uses FT_RWIPV6ROUTING format (which has wider time duration) for Netflow v5
