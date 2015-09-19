---
layout: default
title: CLI
date: 2015-09-09 00:00:00
---

# Scalingo Command Line Tool

## Installation

Just copy/paste the command below in your terminal and execute it.

<div class='form-group install'>
  <div class='input-group'>
    <div class='input-group-addon cli-logo'>
      <i class='fa fa-terminal'></i>
    </div>
    <input class='form-control input-lg' readonly='readonly' type='text' value='curl -O https://cli-dl.scalingo.io/install &amp;&amp; bash install' style='background-color:white;font-size: 22px;'>
  </div>
</div>

## Document and features

Read everything in our [Documentation Center](http://doc.scalingo.com/app/command-line-tool.html).

## Supported operating systems

<ul class='list-inline h4' style='margin-left:-15px;'>
  <li>
  </li>
  <li>
    <i class='fa fa-linux'></i>
    Linux
  </li>
  <li>
    <i class='fa fa-apple'></i>
    MacOS
  </li>
  <li>
    <i class='fa fa-windows'></i>
    Windows
  </li>
  <li>
    <img src='/assets/images/freebsd.png' style='width:16px;'>
    FreeBSD
  </li>
  <li>
    <img src='/assets/images/openbsd.png' style='width:20px;'>
    OpenBSD
  </li>
</ul>

For each of these operating systems, you can download
[a precompiled binary](https://github.com/Scalingo/cli/releases).

The client is developed with Go. Therefore, there is no external dependency.

## Changelog

__1.1.0__

* [Feature - CLI] Setup Bash and ZSH completion thanks to codegangsta/cli helpers [#127](https://github.com/Scalingo/issues/127)
* [Feature - CLI] Add -r/--remote flag to specify a GIT remote instead of an app name [#89](https://github.com/Scalingo/issues/89)
* [Feature - CLI] Add -r/--remote flag to the `create` subcommand to specify an alternative git remote name (default is `scalingo`) [#129](https://github.com/Scalingo/issues/129)
* [Feature - Log] Add -F/--filter flag to filter log output by container types [#118](https://github.com/Scalingo/issues/118)
* [Bugfix - Run] Fix parsing of environment variables (flag -e) [#119](https://github.com/Scalingo/issues/119)
* [Bugfix - Mongo Console] Do not try to connect to the oplog user anymore (when enabled) [#117](https://github.com/Scalingo/issues/117)
* [Bugfix - Logs] Stream is cut with an 'invalid JSON' error, fixed by increasing the buffer size [#135](https://github.com/Scalingo/issues/135)
* [Bugfix - Tunnel] Error when the connection to the database failed, a panic could happen

__1.0.0__

* [Feature - Databases] Add helper to run interactive console for MySQL, PostgreSQL, MongoDB and Redis [#111](https://github.com/Scalingo/cli/issues/111)
* [Feature - Collaborators] Handle collaborators directly from the command line client [#113](https://github.com/Scalingo/cli/issues/113)
* [Feature - Proxy] Add support and documentation about how to use a HTTPS proxy [#104](https://github.com/Scalingo/cli/issues/104) [#110](https://github.com/Scalingo/cli/issues/110)
* [Refactoring - API calls] Completely refactor error management for Scalingo API calls [#112](https://github.com/Scalingo/cli/issues/112)
* [Improvement - SSL] Embed Scalingo new SSL certificate SHA-256 only [#109](https://github.com/Scalingo/cli/issues/109)
* [Bugfix - Macos] [#105](https://github.com/Scalingo/cli/issues/105) [#114](https://github.com/Scalingo/cli/issues/114)
* [Bugfix - Logs] No more weird error message when no log is available for an app [#108](https://github.com/Scalingo/cli/issues/108)
* [Bugfix - Logs] Use of websocket for log streaming [#86](https://github.com/Scalingo/cli/issues/86) [#115](https://github.com/Scalingo/cli/issues/115) [#116](https://github.com/Scalingo/cli/issues/116)
* [Bugfix - Windows] Babun shell compatibility [#106](https://github.com/Scalingo/cli/issues/106)


__1.0.0-rc1__

* [Feature] Modify size of containers with `scalingo scale` - [#102](https://github.com/Scalingo/cli/issues/102)
* [Bugfix] Fix ssh-agent error when no private key is available - Fixed [#100](https://github.com/Scalingo/cli/issues/100)
* [Bugfix] Fix domain-add issue. (error about domain.crt file) - Fixed [#98](https://github.com/Scalingo/cli/issues/98)
* [Bugfix] Fix addon plans description, no more HTML in them  - [#96](https://github.com/Scalingo/cli/issues/96)
* [Bugfix] Correctly handle db-tunnel when alias is given as argument - Fixed [#93](https://github.com/Scalingo/cli/issues/93)


__1.0.0-beta1__

* Windows, password: don't display password in clear
* Windows, db-tunnel: Correctly handle SSH key path with -i flag
* Send OS to one-off containers (for prompt handling, useful for Windows)
* Fix EOF error when writing the password
* Fix authentication request to adapt the API
* Correctly handle 402 errors (payment method required) [#90](https://github.com/Scalingo/cli/issues/90)
* Project is go gettable `go get github.com/Scalingo/cli`
* Fix GIT remote detection [#89](https://github.com/Scalingo/cli/issues/89)
* Correctly handle 404 Error, display clearer messages [#91](https://github.com/Scalingo/cli/issues/91)
* More documentation for the `run` command - Fixed [#79](https://github.com/Scalingo/cli/issues/79)
* Rewrite API client package, remove unsafe code - Fixed [#80](https://github.com/Scalingo/cli/issues/80)
* Allow environment variable name or value for `db-tunnel` as argument
* Extended help for `db-tunnel` - Fixed [#85](https://github.com/Scalingo/cli/issues/85)
* Ctrl^C doesn't kill an `run` command anymore - Fixed [#83](https://github.com/Scalingo/cli/issues/83)
* --app flag can be written everywhere in the command line - Fixed [#10](https://github.com/Scalingo/cli/issues/10)
* Use SSH agent if possible to get SSH credentials
* Correcty handle encrypted SSH keys (AES-CBC and DES-ECE2) - Fixed [#76](https://github.com/Scalingo/cli/issues/76), [#77](https://github.com/Scalingo/cli/issues/77)


__1.0.0-alpha4__

* Adapt to Scalingo API modifications
* Do not encode HTML entities anymore - command: logs
* New login command - command: login
* Allow to use encrypted SSH key (AES-128-CBC) - command: db-tunnel


__1.0.0-alpha3__

* Fix credential storage issue - fixed [#72](https://github.com/Scalingo/cli/issues/72), [#73](https://github.com/Scalingo/cli/issues/73)
* Fix wrong help for command 'db-tunnel' - fixed [#74](https://github.com/Scalingo/cli/issues/74)
* Fix logfile open operation on MacOS - fixed [#70](https://github.com/Scalingo/cli/issues/70)
* Build Windows version on Windows with CGO - fixed [#71](https://github.com/Scalingo/cli/issues/71)
* Build Mac OS verison on Mac OS with CGO - fixed [#71](https://github.com/Scalingo/cli/issues/71)


__1.0.0-alpha2__

* Move addons-related commands to toplevel
  * new-command: addons-add &lt;addon&gt; &lt;plan&gt;
  * new-command: addons-remove &lt;addon-id&gt;
  * new-command: addons-upgrade &lt;addon-id&gt; &lt;plan&gt;

__1.0.0-alpha1__

* First public draft
