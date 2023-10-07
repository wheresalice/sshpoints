package servers

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"io"
	"net/http"
)

const publicHost = "sshpoints.wheresalice.info"

var hctx = context.Background()

func HTTP(redisConnection string) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisConnection,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// @todo render this nicer
		w.Header().Add("Content-Type", "text/plain")

		io.WriteString(w, "# SSH Points\n\n")

		io.WriteString(w, fmt.Sprintf("Get to the top of the leaderboards by SSHing to %s on port 2222 from the most different places\n\n", publicHost))

		io.WriteString(w, "Privacy: We store your username, IP addresses, ASN Numbers, and countries forever when you SSH to this server. We don't store anything when you visit via a web browser\n\n")

		io.WriteString(w, "## Users by IPs\n\n")
		users, _ := rdb.ZRevRangeWithScores(hctx, "userIPLeaderboard", 0, 9).Result()
		for i := range users {
			io.WriteString(w, fmt.Sprintf("- %s %v\n", users[i].Member, users[i].Score))
		}

		io.WriteString(w, "\n## Users by Countries\n\n")
		users, _ = rdb.ZRevRangeWithScores(hctx, "userCountryLeaderboard", 0, 9).Result()
		for i := range users {
			io.WriteString(w, fmt.Sprintf("- %s %v\n", users[i].Member, users[i].Score))
		}

		io.WriteString(w, "\n## Users by ASN Numbers\n\n")
		users, _ = rdb.ZRevRangeWithScores(hctx, "userASNLeaderboard", 0, 9).Result()
		for i := range users {
			io.WriteString(w, fmt.Sprintf("- %s %v\n", users[i].Member, users[i].Score))
		}
	})
	// @todo make port configurable
	fmt.Println("Listening on port 3333")
	http.ListenAndServe(":3333", nil)
}
