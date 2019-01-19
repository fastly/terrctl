[![Build Status](https://travis-ci.com/fastly/terrctl.svg?branch=master)](https://travis-ci.com/fastly/terrctl?branch=master)

# ![terrctl](https://github.com/fastly/terrctl/raw/master/logo.png)

`terrctl` uploads source code leveraging the [Fastly Labs Terrarium](https://wasm.fastlylabs.com) [API](https://wasm.fastlylabs.com/docs) directly to the Terrarium sandbox.

## Download

Pre-built binaries for most platforms are available [in the release section](https://github.com/fastly/terrctl/releases/latest).

## Usage

```sh
Usage: terrctl [options] <source code path>

  -deploy-timeout uint
    	Timeout for deployment (seconds) (default 90)
  -health-timeout uint
    	Timeout for health checks (seconds) (default 30)
  -http-timeout uint
    	Timeout for HTTP client queries (seconds) (default 30)
  -language string
    	language (auto|c|rust|assemblyscript|wasm) (default "auto")
  -logfile string
    	Write logs to file
  -loglevel value
    	Log level (0-6) (default 1)
  -max-deploy-attempts uint
    	Maximum number of attempts for deployment (default 10)
  -syslog
    	Send logs to the local system logger (Eventlog on Windows, syslog on Unix)
```

## Demo

```text
./terrctl /tmp/src/image_example

[2019-01-19 00:30:36] [INFO] Preparing upload of directory [/tmp/src/image_example]
[2019-01-19 00:30:36] [INFO] Guessed programming language: c
[2019-01-19 00:30:36] [NOTICE] Upload in progress...
[2019-01-19 00:30:42] [NOTICE] Upload done, compilation in progress...
[2019-01-19 00:30:43] [INFO] Upload complete, waiting for build...
[2019-01-19 00:30:44] [INFO] Building...
[2019-01-19 00:30:51] [INFO] Generating machine code...
[2019-01-19 00:31:00] [INFO] Deploying...
[2019-01-19 00:31:02] [INFO] Deploy complete: https://capital-telephone-electricity-since.fastly-terrarium.com/
[2019-01-19 00:31:02] [INFO] Instance is deployed
[2019-01-19 00:31:23] [NOTICE] Instance is running and reachable over HTTPS
[2019-01-19 00:31:23] [NOTICE] New instance deployed at [https://capital-telephone-electricity-since.fastly-terrarium.com]
```

## Contact

`labs`@`fastly`.`com`
