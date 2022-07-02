package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/mail"
	"os"
	"regexp"
	"strings"
	"time"
)

const defaultLog = "stderr"
const defaultPat = "([a-f0-9-]+)_(\\d+)_(\\d+)"
const defaultURL = "https://example.com/unsubscribe/?u=$1&c=$2&l=$3"

func main() {
	var (
		from     = flag.String("from", "", "Envelope sender address (MAIL FROM)")
		to       = flag.String("to", "", "Envelope recipient address (RCPT TO)")
		logFile  = flag.String("log", defaultLog, "Log to specified file")
		addrPat  = flag.String("pat", defaultPat, "Regex pattern to extract parameters from recipient address")
		template = flag.String("url", defaultURL, "URL template with parameter substitutions for one-click unsubscribe")
	)
	// return exit code 64 (EX_USAGE) instead of 2 in case of Error
	flag.CommandLine.Init(os.Args[0], flag.ContinueOnError)
	err := flag.CommandLine.Parse(os.Args[1:])
	if err != nil {
		os.Exit(64) // EX_USAGE
	}

	if *to == "" {
		log.Printf("Parameter -to not specified")
		os.Exit(64) // EX_USAGE
	}

	re, err := regexp.Compile(*addrPat)
	if err != nil {
		log.Printf("Parameter -pat not valid: %v", err)
		os.Exit(64) // EX_USAGE
	}

	sender := *from
	recipient := *to

	if *logFile != defaultLog {
		// NOTE: logrotate does not work if not run as deamon
		writer, err := os.OpenFile(*logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY|os.O_SYNC, 0644)
		if err != nil {
			log.Printf("Error opening log file %v", err)
			os.Exit(73) // EX_CANTCREAT
		}
		defer writer.Close()
		log.SetOutput(writer)
	}

	message, err := mail.ReadMessage(os.Stdin)
	if err != nil {
		log.Printf("Unable to parse mail: %v", err)
		os.Exit(65) // EX_DATAERR, will cause 5.5.2 (syntax error)
	}
	subject := message.Header.Get("Subject")
	io.Copy(ioutil.Discard, message.Body)

	log.Printf("Received mail from <%s> to <%s> with subject %q", sender, recipient, subject)

	url, err := replace(re, *template, recipient)
	if err != nil {
		log.Printf("Recipient address %s does not match %s", recipient, *addrPat)
		os.Exit(67) // EX_NOUSER, will cause 5.1.1 bounce
	}

	http.DefaultClient.Timeout = time.Minute

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error calling %s: %v", url, err)
		os.Exit(75) // EX_TEMPFAIL, will cause 4.0.0 bounce
	}
	resp.Body.Close()
	if !strings.HasPrefix(resp.Status, "2") {
		log.Printf("Error calling %s: %s", url, http.StatusText(resp.StatusCode))
		os.Exit(75) // EX_TEMPFAIL, will cause 4.0.0 bounce
	}
	log.Printf("Succes calling %s", url)
}
