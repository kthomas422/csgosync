## Will synchronize your Counter Strike Go maps with a server

Not really tied to csgo but a client/server application that synchronizes
a client directory with a server directory.

The application works by the client taking a hash of the files in the directory
and sends those hashes along with the file names to the server. The server will
compare each filename and hash to its own list of filenames and hashes. Any files
missing from the server's list or where with a mismatched hash will be sent back
to the client. The client will then request a download of each of those files.

#### Server:
Start the server by running the `csgosyncd` (the "d" at the end). It will
load the settings in `csgosyncd.yaml`. Notable settings to change is the *PASSWORD*
setting and your *MAP_PATH* setting. If you want to log to a different filename or
to standard error/standard out you may set those as well.

The settings can also be set with environment variables instead.

#### Client:
The client can be started from a terminal or by "double clicking". It will need the
matching password for the server obviously and the url for the server. The map path
is already set for typical CSGO installations.

For the *URI* it must contain the dns/ip address of the server and the port number (:8080)
default. It can optionally include the `http://`.
