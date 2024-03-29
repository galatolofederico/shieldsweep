# shieldsweep

> 📢 shieldsweep is currently in its early stage of development. Breaking changes may occur!

**Shieldsweep** is a **security analysis tool** written in **Go** designed to **fortify** your systems by **integrating and running** a **suite** of **well-known security utilities**.

## Screenshots

|                        |                        |                        |
|:----------------------:|:----------------------:|:----------------------:|
| ![](README.md.d/1.png) | ![](README.md.d/2.png) |![](README.md.d/3.png)  | 


## Supported Tools

The following table lists the tools currently supported by Shieldsweep.

| Tool      | Supported |
| --------- | :-------: |
| rkhunter  |    ✅     |
| chkrootkit|    ✅     |
| lynis     |    ✅     |


## Features

The following table outlines the current and planned features for Shieldsweep

| Feature                                       | Status                |
| --------------------------------------------- | :-------------------: |
| Basic daemon and scanning functionalities     |  :white_check_mark:   |
| CLI interface                                 |  :white_check_mark:   |
| Web interface                                 |  :white_check_mark:   |
| Notifications                                 |  :white_check_mark:   |
| Log history                                   |  :white_check_mark:   |
| Telegram bot                                  |  :construction:   |

## Installation

To install shieldsweep clone this repository

```
git clone https://github.com/galatolofederico/shieldsweep.git
cd shieldsweep
```

Build and install the project

```
make
sudo make install
```

The `shsw-daemon` should now be up and running you can dobule check it with

```
systemctl status shsw
```

## Usage

You can use the CLI tool `shsw` to interact with the daemon 

| Command                     | Description                                              |
| --------------------------- | -------------------------------------------------------- |
| `shsw status`               | Check the current state of the shieldsweep daemon.       |
| `shsw`                      | Run a scan using the integrated suite of security tools. |
| `shsw list <tool>`          | List the logs for a specified tool.                      |
| `shsw log <tool> <logid>`   | Read a specific log for a tool using the log ID.         |

Or you can use the web-app:

```
shsw-web
```

You will find the dashboard at `http://localhost:3000/`

## Configuration

You can edit the configuration file `/etc/shsw/shsw.json` to enable/disable tools, specify settings for each tool, adjust the level of parallelism, and set up custom notification commands.

```json
{
    "parallelism": 2,
    "notifications": [
        {
            "type": "command",
            "config": {
                "command": [
                    "/bin/sh",
                    "-c",
                    "wall \"New logs available in Shieldsweep\""
                ]
            }
        }
    ],
    "tools": [
        {
            "name": "rkhunter",
            "enabled": true
        },
        {
            "name": "chkrootkit",
            "enabled": true
        },
        {
            "name": "lynis",
            "enabled": true
        }
    ]
}
```

## Development

Current development backlog:
- [ ] Switch from fiber to native go 1.22 HTTP server
- [ ] Daemon interaction refactor: create a common (non-internal) package to interact with the daemon (used by cli, web, telegram, etc...) to abstract the HTTP over unix socket interface.
- [ ] Create a CI/CD pipeline to build and release the software
- [ ] Write a simple script just to handle Telegram notifications (to be used while the actual Telegram client is in development)

## License 

shieldsweep is released under the GNU General Public License v3.0 (GPLv3).