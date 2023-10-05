# SSH points

Earn points for connecting over SSH from unique IPs and networks

## Running


```bash
docker run -ti --rm -v redis:/data --publish 6379:6379 redis  # launch a redis docker container
ssh-keygen -f /tmp/id_rsa  # generate a server ssh key (path currently hardcoded)
sshpoints s --redis localhost:6379  # starts the ssh server on port 2222
sshpoints h --redis localhost:6379  # starts the scoreboard http server on port 3333
```

## Usage

Earn two new IPs by connecting over IPv4 and IPv6

```bash
ssh -4 -p 2222 <user>@<server>  # SSH into the server over IPv4 as a user
ssh -6 -p 2222 <user>@<server>  # SSH into the server over IPv6 as a user
```

The server will tell you how many IPs, countries, and ASN Numbers you've connected from

## Security & privacy

There is no authentication on the SSH server. You could ssh as anyone and score points for that user

We store the IP addresses that you connect from along with the associated ASN Numbers and countries. These are linked to your username, but do not include timestamps. We do not store your SSH key or any other information. Since there is no authentication, this data isn't strictly speaking associated with an individual.

We do not store any data if you visit the HTTP server, though proxies between us may do.
