package supervisor

import (
	"github.com/coreos/go-systemd/sdjournal"
	"github.com/m7shapan/ratelimit"
	"github.com/sirupsen/logrus"
	"github.com/vjeantet/grok"
	"math"
	"os"
	"time"
)

var (
	patterns *grok.Grok
)

func WatchSshd() {
	logrus.Info("Creating ssh login rate limit...")
	sshRateLimit := ratelimit.CreateLimit("5r/h")

	logrus.Info("Loading ssh log patterns...")
	LoadPatterns()
	//go func() {
	//	for {
	//		time.Sleep(time.Second)
	//		log.Printf("%+15s", rl.Rates)
	//	}
	//}()

	jr, err := sdjournal.NewJournalReader(sdjournal.JournalReaderConfig{
		Since: 1 * time.Nanosecond,
		Formatter: func(entry *sdjournal.JournalEntry) (string, error) {
			// parse timestamp from log
			//ts, err := strconv.ParseInt(entry.Fields["_SOURCE_REALTIME_TIMESTAMP"], 10, 64)
			//if err != nil {
			//	panic(err)
			//}
			//log.Printf("%v", time.Unix(0, i*1000).String())

			// parse log entry
			parsedEntry := ParseLog(entry.Fields["MESSAGE"])
			// skip log entry if we don't have IP
			if len(parsedEntry.ClientIp) < 1 {
				return "", nil
			}

			// count rate for parsed IP
			err := sshRateLimit.Hit(parsedEntry.ClientIp)
			if err != nil {
				// too much requests, report this
				ReportIp(parsedEntry.ClientIp, entry.Fields["MESSAGE"])
			}

			// we don't need real output so skip it
			return "", nil
		},
		Matches: []sdjournal.Match{
			{
				Field: sdjournal.SD_JOURNAL_FIELD_SYSTEMD_UNIT,
				Value: "ssh.service",
			},
		},
		//Matches: []sdjournal.Match{
		//	{"SYSLOG_IDENTIFIER", "sshd.service"},
		//},
	})

	if err != nil {
		logrus.Error(err)
	}

	if jr == nil {
		logrus.Error("nil journal reader")
	}

	//b := make([]byte, 64*1<<(10)) // 64KB.
	//for {
	//	c, err := jr.Read(b)
	//	if err != nil {
	//		if err == io.EOF {
	//			break
	//		}
	//		panic(err)
	//	}
	//	logrus.Info(string(b[:c]))
	//}

	logrus.Info("Starting ssh log parsing...")
	defer jr.Close()
	err = jr.Follow(time.After(time.Duration(math.MaxInt64)), os.Stdout)
	if err != nil {
		logrus.Error(err)
	}

	logrus.Info("ssh log parsing stopped")
}

type ParsedSshLog struct {
	ClientIp string
}

func LoadPatterns() {
	patterns, _ = grok.NewWithConfig(&grok.Config{NamedCapturesOnly: true})
	// Invalid user systemz from 192.168.121.1 port 40060
	patterns.AddPattern("SSHD_BAD_USERNAME_PASSWORD", `Invalid user %{USER:client_username} from %{IP:client_ip} port %{INT:client_port}`)
	// Failed password for vagrant from 192.168.121.1 port 42270 ssh2
	patterns.AddPattern("SSHD_BAD_PASSWORD", `Failed password for %{USER:client_username} from %{IP:client_ip} port %{INT:client_port} %{WORD:ssh_version}`)
	// Invalid user gqh from 192.168.121.1 port 35340
	patterns.AddPattern("SSHD_BAD_USERNAME", `Invalid user %{USER:client_username} from %{IP:client_ip} port %{INT:client_port}`)
}

func ParseLog(str string) (res ParsedSshLog) {
	patternList := []string{
		"%{SSHD_BAD_USERNAME_PASSWORD}",
		"%{SSHD_BAD_PASSWORD}",
		"%{SSHD_BAD_USERNAME}",
	}
	for _, v := range patternList {
		parsedValues, err := patterns.Parse(v, str)
		//for k, v := range parsedValues {
		//	fmt.Printf("%+15s: %s\n", k, v)
		//}
		if len(parsedValues["client_ip"]) < 1 || err != nil {
			continue
		}
		res = ParsedSshLog{
			ClientIp: parsedValues["client_ip"],
		}
		// we already have what we want
		break
	}
	return
}
