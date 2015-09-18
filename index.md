---
layout: default
title: CLI
date: 2015-09-09 00:00:00
---

<h1>Scalingo Command Line Tool</h1>

<h2>Installation</h2>

<p>Just copy/paste the command below in your terminal and execute it.</p>

<div class='form-group install'>
  <div class='input-group'>
    <div class='input-group-addon cli-logo'>
      <i class='fa fa-terminal'></i>
    </div>
    <input class='form-control input-lg' readonly='readonly' type='text' value='curl -O https://cli-dl.scalingo.io/install &amp;&amp; bash install' style='background-color:white;font-size: 22px;'>
  </div>
</div>

<h2>Document and features</h2>

<p>
  Read everything in our <a href='http://doc.scalingo.com/app/command-line-tool.html'>Documentation Center</a>.
</p>

<h2>Supported operating systems</h2>

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
<p>
  For each of these operating systems, you can download
  <a href='https://github.com/Scalingo/cli/releases'>a precompiled binary</a>.
</p>
<p>The client is developed with Go. Therefore, there is no external dependency.</p>

<h2>Changelog</h2>

<strong>1.0.0</strong>
<ul>
<li>[Feature - Databases] Add helper to run interactive console for MySQL, PostgreSQL, MongoDB and Redis <a href="https://github.com/Scalingo/cli/issues/111">#111</a></li>
<li>[Feature - Collaborators] Handle collaborators directly from the command line client <a href="https://github.com/Scalingo/cli/issues/113">#113</a></li>
<li>[Feature - Proxy] Add support and documentation about how to use a HTTPS proxy <a href="https://github.com/Scalingo/cli/issues/104">#104</a> <a href="https://github.com/Scalingo/cli/issues/110">#110</a></li>
<li>[Refactoring - API calls] Completely refactor error management for Scalingo API calls <a href="https://github.com/Scalingo/cli/issues/112">#112</a></li>
<li>[Improvement - SSL] Embed Scalingo new SSL certificate SHA-256 only <a href="https://github.com/Scalingo/cli/issues/109">#109</a></li>
<li>[Bugfix - Macos] <a href="https://github.com/Scalingo/cli/issues/105">#105</a> <a href="https://github.com/Scalingo/cli/issues/114">#114</a></li>
<li>[Bugfix - Logs] No more weird error message when no log is available for an app <a href="https://github.com/Scalingo/cli/issues/108">#108</a></li>
<li>[Bugfix - Logs] Use of websocket for log streaming <a href="https://github.com/Scalingo/cli/issues/86">#86</a> <a href="https://github.com/Scalingo/cli/issues/115">#115</a> <a href="https://github.com/Scalingo/cli/issues/116">#116</a></li>
<li>[Bugfix - Windows] Babun shell compatibility <a href="https://github.com/Scalingo/cli/issues/106">#106</a></li>
</ul>

<strong>1.0.0-rc1</strong>
<ul>
<li>[Feature] Modify size of containers with `scalingo scale` - <a href="https://github.com/Scalingo/cli/issues/102">#102</a></li>
<li>[Bugfix] Fix ssh-agent error when no private key is available - Fixed <a href="https://github.com/Scalingo/cli/issues/100">#100</a></li>
<li>[Bugfix] Fix domain-add issue. (error about domain.crt file) - Fixed <a href="https://github.com/Scalingo/cli/issues/98">#98</a></li>
<li>[Bugfix] Fix addon plans description, no more HTML in them  - <a href="https://github.com/Scalingo/cli/issues/96">#96</a></li>
<li>[Bugfix] Correctly handle db-tunnel when alias is given as argument - Fixed <a href="https://github.com/Scalingo/cli/issues/93">#93</a></li>
</ul>

<strong>1.0.0-beta1</strong>
<ul>
<li>Windows, password: don't display password in clear</li>
<li>Windows, db-tunnel: Correctly handle SSH key path with -i flag</li>
<li>Send OS to one-off containers (for prompt handling, useful for Windows)</li>
<li>Fix EOF error when writing the password</li>
<li>Fix authentication request to adapt the API</li>
<li>Correctly handle 402 errors (payment method required) <a href="https://github.com/Scalingo/cli/issues/90">#90</a></li>
<li>Project is go gettable `go get github.com/Scalingo/cli`</li>
<li>Fix GIT remote detection <a href="https://github.com/Scalingo/cli/issues/89">#89</a></li>
<li>Correctly handle 404 Error, display clearer messages <a href="https://github.com/Scalingo/cli/issues/91">#91</a></li>
<li>More documentation for the `run` command - Fixed <a href="https://github.com/Scalingo/cli/issues/79">#79</a></li>
<li>Rewrite API client package, remove unsafe code - Fixed <a href="https://github.com/Scalingo/cli/issues/80">#80</a></li>
<li>Allow environment variable name or value for `db-tunnel` as argument</li>
<li>Extended help for `db-tunnel` - Fixed <a href="https://github.com/Scalingo/cli/issues/85">#85</a></li>
<li>Ctrl^C doesn't kill an `run` command anymore - Fixed <a href="https://github.com/Scalingo/cli/issues/83">#83</a></li>
<li>--app flag can be written everywhere in the command line - Fixed <a href="https://github.com/Scalingo/cli/issues/10">#10</a></li>
<li>Use SSH agent if possible to get SSH credentials</li>
<li>Correcty handle encrypted SSH keys (AES-CBC and DES-ECE2) - Fixed <a href="https://github.com/Scalingo/cli/issues/76">#76</a>, <a href="https://github.com/Scalingo/cli/issues/77">#77</a></li>
</ul>

<strong>1.0.0-alpha4</strong>
<ul>
<li>Adapt to Scalingo API modifications</li>
<li>Do not encode HTML entities anymore - command: logs</li>
<li>New login command - command: login</li>
<li>Allow to use encrypted SSH key (AES-128-CBC) - command: db-tunnel</li>
</ul>

<strong>1.0.0-alpha3</strong>
<ul>
<li>Fix credential storage issue - fixed <a href="https://github.com/Scalingo/cli/issues/72">#72</a>, <a href="https://github.com/Scalingo/cli/issues/73">#73</a></li>
<li>Fix wrong help for command 'db-tunnel' - fixed <a href="https://github.com/Scalingo/cli/issues/74">#74</a></li>
<li>Fix logfile open operation on MacOS - fixed <a href="https://github.com/Scalingo/cli/issues/70">#70</a></li>
<li>Build Windows version on Windows with CGO - fixed <a href="https://github.com/Scalingo/cli/issues/71"><a href="https://github.com/Scalingo/cli/issues/71">#71</a></a></li>
<li>Build Mac OS verison on Mac OS with CGO - fixed <a href="https://github.com/Scalingo/cli/issues/71"><a href="https://github.com/Scalingo/cli/issues/71">#71</a></a></li>
</ul>

<strong>1.0.0-alpha2</strong>
<ul>
<li>
Move addons-related commands to toplevel
<ul>
<li>new-command: addons-add &lt;addon&gt; &lt;plan&gt;</li>
<li>new-command: addons-remove &lt;addon-id&gt;</li>
<li>new-command: addons-upgrade &lt;addon-id&gt; &lt;plan&gt;</li>
</ul>
</li>
</ul>

<strong>1.0.0-alpha1</strong>
<ul>
<li>First public draft</li>
</ul>
