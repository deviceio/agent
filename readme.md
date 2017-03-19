<aside class="notice">
This is alpha grade software and is not suitable for production use. Breaking changes are frequent as we march towards 1.0 and our <a href="">1.0 Compatability Promise</a>
</aside>

# 1. Deviceio Agent

The Deviceio Agent is a small binary that runs on a target host and connects to a Deviceio Hub. The agent provides an http api over its hub connection that developers and system administrators can use to orchestrate the target host machine.

For more detailed information regarding installation, configuration, connectivity, and security please read the [](User Manual)

The remainder of this readme will document the Deviceio Agent API.

<!-- TOC -->

- [1. Deviceio Agent](#1-deviceio-agent)
- [2. Url Paths](#2-url-paths)
- [3. Endpoints](#3-endpoints)
- [4. RPC endpoints](#4-rpc-endpoints)
    - [4.1. POST /rpc/filesystem/read](#41-post-rpcfilesystemread)
        - [4.1.1. Argument Headers](#411-argument-headers)
        - [4.1.2. Returns](#412-returns)
        - [4.1.3. Trailers](#413-trailers)

<!-- /TOC -->

# 2. Url Paths

The paths refected in this documentation are *relative* to the Deviceio Hub device resource path. For example, when the documentation describes a api path for an agent:

```
GET /rest/filesystem/{path}
```

in usage means calling the hub http api as follows:

```
GET https://{hub_url}/device/{id_or_hostname}/rest/filesystem/{path}
```

# 3. Endpoints

Agents exposes different style endpoints from root paths:

```
/rest/   <- all REST style apis are rooted here
/rpc/    <- all RPC style apis are rooted here
/dsc/    <- all DSC style apis are rooted here
```

# 4. RPC endpoints

RPC endpoints are useful for operations that do not fit within a resource based RESTful paradigm or operate exclusivly on streamed data. URLs follow pattern /rpc/{ref}/{method} where `ref` is a reference to a static module or instanced object and `method` is a reference to a method on the parent reference. 

The majority of work will be conducted on static module references and their exported methods. However, some modules may expose methods that create instanced objects that offer more advanced orchestration, but MAY require user cleanup after the object instances are no longer needed. Object instances are referenced by GUIDs in the form `/rpc/{ref:GUID}/{method}` 

Static module methods that return object instances will return a `Location: /rpc/{ref:GUID}` HTTP header that can be used to locate the object instance after the static module method call is completed. After the instanced object is no longer needed users MUST make an Http call in the form of `DELETE /rpc/{ref:GUID}` to dispose of the object instance, unless the object documentation explicitly states that it is self disposing under some condition.

All method calls on references MUST be issued as a `POST` request. Arguments to methods are supplied through HTTP Headers (as to not disrupt request/response streaming). Methods MAY support request and/or response streaming in various ways as documented.

## 4.1. POST /rpc/filesystem/read

Reads a file from the local filesystem of the device and returns its contents. The response body ALWAYs returns with HTTP `Transfer-Encoding: chunked` to accomodate large files for clients that can support chunked reading. Due to chunked response, errors are supplied in a trailing `Error: {error}` header in the response. If either side disconnects from the http connection, no error will be provided as a disconnect is not an error in this endpoint. To validate the length of the returned content you must obtain the total content length through some other agent api endpoint.

### 4.1.1. Argument Headers

* `X-Path <string>`: The path to the local file on the device
* `X-Offset <int>`: The byte offset to start the file read from. `Default: 0`. If the byte offset exceeds the length of the file an error will be thrown
* `X-OffsetAt <int>`: Where to place the `X-Offset`. `Default: 0`. 0 means offset starts at the origin (beginning) of the file. 1 means offset starts at the end of the file moving backwards.
* `X-Count <int>`: The total number of bytes to read from `X-Offset`. `Default: -1`. If the count exceeds the total length starting from the offset, only the total number of bytes will be read and no error thrown. 

### 4.1.2. Returns

Byte content using `Transfer-Encoding: chunked` response streaming 

### 4.1.3. Trailers

* `Error`: Any error that is observed during the read operation