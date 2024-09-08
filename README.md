# Go Man Page Web Server
This Go application serves man pages over HTTP. It allows users to view, search, and browse man pages directly from a web browser. The server dynamically integrates the hostname into the pages and provides static file support (CSS and favicon).

## Table of Contents
- [Features](#features)
- [Dependencies](#dependencies)
- [Compiling The Binary](#compiling-the-binary)
- [Systemd Service](#using-as-systemd-service)
- [Accessing the Web Interface](#accessing-the-web-interface)
- [Notes and Warnings](#notes-and-warnings)

## Features
- Convert any man page into HTML, formatted for easy reading. With a dark theme.
- Search functionality to find specific man pages. Wildcards and regex are supported.
- Hyperlinked pages for easy navigation for any valid reference in the document, and the ability to open the page in a new tab.
- Will display all man pages in the manpath, including pages where `man2html` and `man -Thtml` fail.
- Filter by page function or section: 1=commands, 3=C/C++ Refs, 5/7=conf/formatting, 8=sudo commands, etc.
- Able to correctly interpret and display incorrectly formatted man pages, to a degree.
- Auto updates man pages when new packages are installed or removed using standard installation methods.

## Dependencies
`mandoc` is required for parsing. Ubuntu/Debian/apt installation:
```sh
sudo apt-get install mandoc -y
```
`go v1.21.5` or greater is required to build the source code. Official installation instructions at [go.dev](https://go.dev/doc/install).\
Quick and lazy script to install go 1.23.1 (Tested on Ubuntu 07-SEP-2024):
```sh
[ -d /usr/local/go ] && sudo rm -rf '/usr/local/go' ; # delete incompatible versions
which wget || sudo apt-get install wget -y && wget https://go.dev/dl/go1.23.1.linux-amd64.tar.gz ; # downlaod compatible version 
[ -f go1.23.1.linux-amd64.tar.gz ] && sudo tar -C /usr/local -xzf go1.23.1.linux-amd64.tar.gz ; # install into system
[ -f /usr/local/go/bin/go ] && sudo ln -s /usr/local/go/bin/go /bin/go ; # add to bin
go version ; # show version
```

## Compiling The Binary

 ```sh
# download the source code
git clone "https://sophuwu.site/manhttpd" && cd manhttpd
 
# build the binary with go
go build -ldflags="-s -w" -trimpath -o build/manhttpd

# install the binary into the system
sudo install ./build/manhttpd /usr/local/bin/manhttpd
```

## Using As Systemd Service:
The provided service file should work on most systems, but you may need to edit it to fit your needs.\
It will open a http server on port 8082 available through all network interfaces.\
You should change the `ListenAddr` variable to `127.0.0.1` and use a secure reverse proxy if you are on a public network.\
TLS and authentication are not implemented in this server.

### Variables in the service file:
Environment Variables:\
`HOSTNAME`: Used for http proxying.\
`ListenPort`: If unset, the server will default to 8082.\
`ListenAddr`: This should be changed if you are on a public network.\
`MANDOCPATH`: Path to the mandoc executable. If unset, the server will attempt to find it in the PATH.

### System Variables:
`User`: Reccomended to use your login user so the service can access your ~/.local man pages. But not required.\
`ExecStart`: The path to the manhttpd binary. If you installed it to /usr/local/bin, you can leave it as is.\
`WorkingDirectory`: This should be /tmp since the server doesn't need to write to disk.\

```sh
# to edit paths, users, and environment variables if needed
nano manhttpd.service 

# install the service file to systemd
sudo install manhttpd.service /etc/systemd/system/manhttpd.service
# reload the systemd daemon to load the new service
sudo systemctl daemon-reload

# to start the service and check its status
sudo systemctl start manhttpd.service
sudo systemctl status manhttpd.service

# if you would like the service to run when the system boots
sudo systemctl enable manhttpd.service

# to stop the service and disable it from starting automatically
sudo systemctl stop manhttpd.service
sudo systemctl disable manhttpd.service

# if you would like to change the variables in the service file
sudo systemctl stop manhttpd.service
sudo systemctl edit manhttpd.service
sudo systemctl daemon-reload 
sudo systemctl start manhttpd.service
 ```


## Accessing the Web Interface
Open your web browser and navigate to `http://localhost:8082` if you are running the server locally or the remote server's IP address or hostname.\
To search with regex, you can use the search bar at the top of the page with `-r` at the beginning of the search term.\
To look into a specific section, you can add `-sN` to the search term where N is the section number.\
If no section is specified, the server will display with the same priority as the defualt `man` command.\
Glob patterns are also supported in the search bar if regex not enabled.\

Example search terms:
- `ls*`: List all pages that begin with `ls`.
- `-r ^ls`: Same as above but with regex.
- `ls` or `ls -s1` or `ls.1`: Open the page for the `ls` user command.
- `-r ^ls -s1`: List all pages that begin with `ls` in section 1 (user/bin commands).
- `*config* -s8`: List all pages that contain `config` within the name and are in section 8 (sudo/sbin commands).  
- `vsftpd.5`: Open the manual page for vsftpd confuguration file if vsftpd is installed.
- `vsftpd.8`: Open the manual page for vsftpd executable if vsftpd is installed.

## Notes and Warnings
- packages installed with apt or dpkg will automatically be available
- manual pages may be unavailable if the package is not installed or the manuals are not included by the package
- all manual pages that are correctly installed and comply with the manpath will be searchable and viewable
- manuals that do not comply with the manpath will not be available for viewing may show up in search results

The application serves the following embedded static files:
- `index.html`: The main HTML template, which includes placeholders for the server's hostname.
- `dark_theme.css`: The CSS stylesheet used to style the web interface. No light theme is provided.
- `favicon.ico`: A favicon served with the site.
I do not recommend changing anything in index.html. CSS classes are hardcoded into the man pages and will not render if changed.\
However, changing the css rules in your browser's developer tools to see how it affects the page is a good way to test changes.\
You can then rebuild the binary with the changes you like to make the changes permanent.
