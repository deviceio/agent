# Deviceio Agent

The Deviceio Agent connects to a Deviceio Hub providing access to the device for real-time orchestration

# Quick Start with Docker

Running the agent within a docker container only provides access to the container environment and will not provide access to the host system. 

Create a v4 uuid to act as the unique ID for your agent instance.

```bash
UUID=$(uuidgen)
echo $UUID
```

Start an agent instance making sure to pass in your v4 UUID  (`--id`)  and ip or hostname of your hub instance  (`--host`). If you have not create a hub instance please review https://github.com/deviceio/hub

```bash
docker run -d --name deviceio-agent-1 deviceio/agent --id $UUID --host [ip_or_hostname_of_hub] --port 8975 --insecure
```

Follow the container logs to ensure we see `transport up` which indicates that the agent successfully connected to your hub instance

```bash
docker logs -f deviceio-agent-1
```

Next:

* Install and interact with your devices via the CLI integration https://github.com/deviceio/cli

