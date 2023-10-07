package servers

import (
	"context"
	"fmt"
	"github.com/gliderlabs/ssh"
	"log"
	"net"

	"github.com/jamesog/iptoasn"

	"github.com/redis/go-redis/v9"
)

var sctx = context.Background()

func SSH(redisConnection string) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisConnection,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	ssh.Handle(func(s ssh.Session) {
		//authorizedKey := gossh.MarshalAuthorizedKey(s.PublicKey())
		//io.WriteString(s, fmt.Sprintf("public key used by %s:\n", s.User()))
		//s.Write(authorizedKey)

		host, _, _ := net.SplitHostPort(s.RemoteAddr().String())
		s.Write([]byte(fmt.Sprintf("hello %s@%s\n", s.User(), host)))

		// register IP
		knownIP, err := rdb.SIsMember(sctx, fmt.Sprintf("%s:ips", s.User()), host).Result()
		if knownIP {
			s.Write([]byte("sorry, this is not a new IP\n"))
		} else {
			//	@todo make transactional
			rdb.SAdd(sctx, fmt.Sprintf("%s:ips", s.User()), host)
			rdb.ZIncrBy(sctx, "userIPLeaderboard", 1, s.User())
		}
		visitedIPs, err := rdb.SCard(sctx, fmt.Sprintf("%s:ips", s.User())).Result()
		s.Write([]byte(fmt.Sprintf("you have visted from %d IPs\n", visitedIPs)))

		// lookup IP
		ip, err := iptoasn.LookupIP(host)
		if err != nil {
			log.Println(err)
			s.Write([]byte("err: failed to lookup ASN data for this IP\n"))
			return
		}

		s.Write([]byte(fmt.Sprintf("you are visiting from %s\n", ip.Country)))
		// register countries
		knownCountry, err := rdb.SIsMember(sctx, fmt.Sprintf("%s:countries", s.User()), ip.Country).Result()
		if knownCountry {
			s.Write([]byte("sorry, this is not a new Country\n"))
		} else {
			//	@todo make transactional
			rdb.SAdd(sctx, fmt.Sprintf("%s:countries", s.User()), ip.Country)
			rdb.ZIncrBy(sctx, "userCountryLeaderboard", 1, s.User())
		}
		visitedCountries, err := rdb.SCard(sctx, fmt.Sprintf("%s:countries", s.User())).Result()
		s.Write([]byte(fmt.Sprintf("you have visted %d countries\n", visitedCountries)))

		s.Write([]byte(fmt.Sprintf("you are visiting from %d %s\n", ip.ASNum, ip.ASName)))
		// register countries
		knownASN, err := rdb.SIsMember(sctx, fmt.Sprintf("%s:asn", s.User()), ip.ASNum).Result()
		if knownASN {
			s.Write([]byte("sorry, this is not a new AS Number\n"))
		} else {
			//	@todo make transactional
			rdb.SAdd(sctx, fmt.Sprintf("%s:asns", s.User()), ip.ASNum)
			rdb.ZIncrBy(sctx, "userASNLeaderboard", 1, s.User())
		}
		visitedASNs, err := rdb.SCard(sctx, fmt.Sprintf("%s:asns", s.User())).Result()
		s.Write([]byte(fmt.Sprintf("you have visted %d AS Numbers\n", visitedASNs)))

	})

	publicKeyOption := ssh.PublicKeyAuth(func(sctx ssh.Context, key ssh.PublicKey) bool {
		return true // allow all keys, or use ssh.KeysEqual() to compare against known keys
	})

	// @todo make port and key configurable
	log.Println("starting ssh server on port 2222...")
	log.Fatal(ssh.ListenAndServe(":2222", nil, publicKeyOption, ssh.HostKeyFile("/etc/ssh/ssh_host_rsa_key")))
}
