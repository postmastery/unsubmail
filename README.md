# unsubmail

PowerMTA addon to support "mailto:" links in List-Unsubscribe headers.

For example:

    List-Unsubscribe: <mailto:229_3430_2346761_618e1a4097f3cf26a8b23f89b9848e5@unsubscribe.example.com>,
	  <https://www.example.com/unsubscribe?q=229&c=2346761&b=3430&p=160&hash=618e1a4097f3cf26a8b23f89b9848e5>

The unsubmail program is run by PowerMTA when it receives unsubscribe emails. It extracts parameters from the receiving address and calls a http link with these parameters. Thus the application only needs to support (one-click) unsubscribes via HTTP.

The List-Unsubscribe header must be added by the application. It is recommended to add a hash parameter which includes a secret code. The hash is checked when the http link is called to prevent manipulation of parameters. 

## Install

### Linux

- Download unsubmail from [https://postmastery.egnyte.com/dl/68U0V2AK3a](https://postmastery.egnyte.com/dl/68U0V2AK3a)
- Upload unsubmail to the PowerMTA into /opt/pmta
- Make sure it's executable for user pmta

### Windows

- Download unsubmail.exe from [https://postmastery.egnyte.com/dl/IHWT0nhSed](https://postmastery.egnyte.com/dl/IHWT0nhSed)
- Move it to \pmta\bin

## Usage

Unsubmail supports the following command line parameters:

	  -from string
	    	Envelope sender address (MAIL FROM)
	  -to string
	    	Envelope recipient address (RCPT TO)
	  -log string
	    	Log to specified file (default "stderr")
	  -pat string
	        Regex pattern to extract parameters from recipient address
	  -url string
	  	    URL template with parameter substitutions for one-click unsubscribe

Please ask Postmastery to help create a regex pattern and url template for your environment.

## Testing

Run unsubmail with sample address:

	unsubmail\
	  -to "ab7c5403-f10d-4a65-b88a-626f02a1fa05_24712345_4072@unsubscribe.example.com"\
	  -pat "([a-f0-9-]+)_(\d+)_(\d+)"\
	  -url "http://host.example.com/unsubscribe?u=\$1&c=\$2&l=\$3"
	Subject: test

	^D

## Configuration

The unsubscribe email address must use a domain which is unused for other email,
for example unsubscribe.example.com.

The MX record for unsubscribe.example.com must refer to the PowerMTA.

	unsubscribe.example.com MX 10 pmta1.example.com
	unsubscribe.example.com MX 10 pmta2.example.com

The domain unsubscribe.example.com must be configured in PowerMTA as follows:

	relay-domain unsubscribe.example.com

	<domain unsubscribe.example.com>
		type pipe
		command "/opt/pmta/unsubmail -from=$from -to=$to -pat=(\\d+)-(\\d+)-([a-f0-9]+)
		  -url=https://host.example.com/unsubscribe/$$1/$$2/$$3 -log=/var/log/pmta/unsubmail.log"
	</domain>

Note that the $from and $to parameters are set by PowerMTA to the envelope sender and
recipient and passed via the command line to the script.

Backslashes in command line parameters must be escaped as "\\". Dollar signs in command line parameters must be escaped as "$$".
