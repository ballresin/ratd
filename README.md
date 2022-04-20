# ratd
Remote Access Toolkit Daemon

IT Asset Remote Management and Monitoring (RMM) platform

# Release Plan
- Release 0
	- Create release plan to tease features
	- Begin Documentation-driven development with rough spec (this)
	- Make release plan public

- Release 1
	- Server
        - [x] Stand up DB
        - [x] Listen on a port
        - [x] Receive checkins
        - [x] Super rough web output
        - [ ] Web portal for agent management
        - [ ] Check agent version, notify agent of new version if applicable
        - [ ] Issue commands from queue if applicable
	- Agent
		- [ ] Checkin regularly
		- [ ] If agent is old, self-update
		- [ ] Receive commands, execute, and relay result to server

- Release 2
	- Per-user remote agent as user daemon, with remote control access

- Release 3
	- Swarm tech to reduce incoming connections to server

- Release 4
	- Promises, or state-detection-and-pursuit
