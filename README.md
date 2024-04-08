# HttpInfoServer
Minecraft HTTP API server that runs locally on the player's computer and hosts in-game player's info. It consists of `Minecraft fabric mod (agent)` and `HTTP API server`. The code of the HTTP API server is located in this repo and the Minecraft fabric mod code is located by this link: [\*link\*](https://github.com/dadencukillia/httpInfoServer-mod). Read the full manual to get everything set up correctly.

# Manual
The full manual you can read here: [read](https://gist.github.com/dadencukillia/006473a91295191963596ab21f9d3b8b)

# Run locally
You have two ways:
- Download the repo code using git (`git clone https://github.com/dadencukillia/httpInfoServer`)
- Run the `go run .` command (don't forget to check if you have already installed Golang version `1.22.x`)

Or:
- Follow the link: [\*link\*](https://github.com/dadencukillia/httpInfoServer/releases)
- Choose and download the newest (highest) release (match the file name with the name of your OS)
- Run the file (terminal may be required)

# Contribution
I will be glad if you can help in the development of this project. What you can do:
- Spelling correction. My English is not good, so you can correct the comments in the code or the text in the files (like [README.md](https://github.com/dadencukillia/httpInfoServer/blob/main/README.md)).
- Code cleanup
- Make the code safer and faster, using pull requests.
- Add new features or improve old ones

Code structure:
- `main.go` — file that contains the "main" function (also known as entrance function). The project uses the `net/http` package to create HTTP servers, so there are two functions that create two routers in the file and these routers run in the "main" function.
- `handlers.go` — file that contains handler functions for these two routes (servers). The project uses the [gorilla/websocket](https://github.com/gorilla/websocket) package and a function handler for the WebSocket server uses it to make a WebSocket connection.