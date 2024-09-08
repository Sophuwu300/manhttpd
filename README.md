# Go Man Page Web Server

This Go application serves man pages over HTTP. It allows users to view, search, and browse man pages directly from a web browser. The server dynamically integrates the hostname into the pages and provides static file support (CSS and favicon).

I hope I find the dedication to translate mandoc into Go one day.
But for now, you need to have `mandoc` installed on and set
the `MANDOCPATH` environment variable to the path of the `mandoc` executable if the http user doesn't have it in its path.


## Features

- Serve UNIX man pages over a web interface.
- Search functionality to find specific man pages.
- Embedded HTML, CSS, and favicon for a simple, customizable UI.
- Dynamic insertion of the server's hostname into the web interface.

## Prerequisites

- `go` to compile the server binary. Exmaple installation for amd64 linux:
   ```sh
   wget https://go.dev/dl/go1.23.1.linux-amd64.tar.gz && \
  sudo tar -C /usr/local -xzf go1.23.1.linux-amd64.tar.gz
   ```
- `mandoc` package for parsing different man page formats.
- `git` for cloning the repository.
    ```sh
    sudo apt-get install mandoc git
    ```

## Build Steps
 ```sh
# download the source code
git clone "https://sophuwu.site/manhttpd" && cd manhttpd
 
# build the binary with go
go build -ldflags="-s -w" -trimpath -o build/manhttpd

# install the binary into the system
sudo install ./build/manhttpd /usr/local/bin/manhttpd
```

## To Use Systemd Service:
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
- `ls` or `ls -s1` or `ls.1`: Open the requested page.
- `-r ^ls -s1`: List all pages that begin with `ls` in section 1 (user/bin commands).
- `*config* -s8`: List all pages that contain `config` within the name and are in section 8 (sudo/sbin commands).  
- `vsftpd.5`: Open the manual page for vsftpd confuguration file if vsftpd is installed.
- `vsftpd.8`: Open the manual page for vsftpd executable if vsftpd is installed.

## Details
- packages installed with apt or dpkg will automatically be available
- manual pages may be unavailable if the package is not installed or the manuals are not included by the package
- all manual pages that are correctly installed and comply with the manpath will be searchable and viewable
- manuals that do not comply with the manpath will not be available for viewing may show up in search results

## Embedded Static Files
The application serves the following embedded static files:
- `index.html`: The main HTML template, which includes placeholders for the server's hostname.
- `dark_theme.css`: The CSS stylesheet used to style the web interface. No light theme is provided.
- `favicon.ico`: A favicon served with the site.
I do not recommend changing anything in index.html. CSS classes are hardcoded into the man pages and will not render if changed.\
However, changing the css rules in your browser's developer tools to see how it affects the page is a good way to test changes.\
You can then rebuild the binary with the changes you like to make the changes permanent.
