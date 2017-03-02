# Deviceio Agent

The Deviceio Agent is a small binary that runs on a target host and connects to a Deviceio Hub. The agent provides an http api over its hub connection that developers and system administrators can use to orchestrate the target host machine.

For more detailed information regarding installation, configuration, connectivity, and security please read the [](User Manual)

The remainder of this readme will document the Deviceio Agent API.

# Previous versions

Please see the various agent version branches in this repository for older agent api documentation

# Url Paths

The paths refected in this documentation are *relative* to the Deviceio Hub device resource path. For example, when the documentation describes a api path for an agent:

```
GET /rest/filesystem/{path}
```

in usage means calling the hub http api as follows:

```
GET https://<hub>/device/<id>/rest/filesystem/{path}
```

# Rest and Rpc endpoints

Deviceio Agents expose both REST and RPC style endpoints. Some orchestration capabilities may have a natural mapping to the REST style, while others may be best described in a RPC style.

Agents expose the different styles from root resource paths:

```
/rest/   <- all REST style apis are rooted here
/rpc/    <- all RPC style apis are rooted here
```

# REST endpoints
## /rest/filesystem/{path}
### GET
### PUT
### POST
### DELETE
### HEAD
### OPTIONS
# RPC endpoints


