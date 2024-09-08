# Go Man Page Web Server

This Go application serves man pages over HTTP. It allows users to view, search, and browse man pages directly from a web browser. The server dynamically integrates the hostname into the pages and provides static file support (CSS and favicon).

## Features

- Convert any man page into HTML, formatted for easy reading. With a dark theme.
- Search functionality to find specific man pages. Wildcards and regex are supported.
- Hyperlinked pages for easy navigation for any valid reference in the document, and the ability to open the page in a new tab.
- Will display all man pages in the manpath, including pages where `man2html` and `man -Thtml` fail.
- Filter by page function or section: 1=commands, 3=C/C++ Refs, 5,7=config/format, 8=sudo commands, etc.
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

# install the service file to systemd and load it
sudo install manhttpd.service /etc/systemd/system/manhttpd.service
sudo systemctl daemon-reload

# start the service and check its status
sudo systemctl start manhttpd.service
sudo systemctl status manhttpd.service

# to keep the service running after a reboot
sudo systemctl enable manhttpd.service

# to stop the service and disable it from restarting
sudo systemctl stop manhttpd.service
sudo systemctl disable manhttpd.service

# to edit the server configuration after installation
sudo systemctl edit manhttpd.service
sudo systemctl daemon-reload 
sudo systemctl reload-or-restart manhttpd.service
sudo systemctl status manhttpd.service
```

## Accessing the Web Interface

Open your web browser and navigate to `http://localhost:8082` if you are running the server locally or the remote server's IP address or hostname.\
To search with regex, you can use the search bar at the top of the page with `-r` at the beginning of the search term.\
To look into a specific section, you can add `-sN` to the search term where N is the section number.\
If no section is specified, the server will display with the same priority as the defualt `man` command.\
Glob patterns are also supported in the search bar if regex not enabled.\

## Example Usage:

- `ls*`: List all pages that begin with `ls`, including `ls`, `lsblk`, `lsmod`, etc.
- `-r ^ls`: Same as above but with regex. Usually more useful for with multiple queries and logical operators. Like finding any C++ reference to `std::string` and `std::vector`.
- `ls` or `ls -s1` or `ls.1`: Open the page for the `ls` user command. This is orignal man behavior.
- `-r ^ls -s1`: List all pages that begin with `ls` in section 1 (user/bin commands). Useful for finding commands that list any information without requiring sudo.
- `*config* -s8`: List pages for sudo commands containing keyword `config`. this can will show you commands that edit critical system files.  
- `vsftpd.5`: Open the manual page for vsftpd confuguration file if vsftpd is installed. This will show you how to configure the ftp server.
- `vsftpd.8`: Open the manual page for vsftpd executable if vsftpd is installed. This will show how to call the ftp server from the command line.

## Notes and Warnings

- Regex syntax may vary when running on systems with differing core c/c++ libraries. 
- Manual pages may be unavailable if the package is not installed or the manuals are not included by the package.
- All manual pages that are correctly installed and comply with the manpath will be searchable and viewable.
- Manuals that do not comply with the manpath will not be available for viewing may show up in search results.
- If you would like to change the css styling, you should avoid changing the names, ids, or classes of the elements as they are dynamically generated. Instead, you can add properties to the existing classes or ids, or change the values of the properties.
- the static web files are embedded into the server binary so changes to css properties will require you to recompile the binary.

## Help and Support

I don't know how this git pull thing works. I will try if I see any issues. I've never collaborated on code before. If you have any suggestions, or questions about anything I've written, I would be happy to hear your thoughts.\
contant info: [skisiel.com](https://skisiel.com) or [sophuwu.site](https://sophuwu.site)


