package main

import (
	"regexp"
	"testing"
)

type testCase struct {
	to  string
	pat string
	url string
	exp string
}

var testCases = []testCase{
	{
		"ab7c5403-f10d-4a65-b88a-626f02a1fa05_24712345_4072@unsubscribe.example.com",
		"([a-f0-9-]+)_(\\d+)_(\\d+)",
		"http://weblaunch.example.com/listener3/unsubscribe?u=$1&c=$2&l=$3",
		"http://weblaunch.example.com/listener3/unsubscribe?u=ab7c5403-f10d-4a65-b88a-626f02a1fa05&c=24712345&l=4072",
	},
	{
		"229_3430_2346761_618e1a4097f3cf26a8b23f89b9848e5@unsubscribe.example.com",
		"(\\d+)_(\\d+)_(\\d+)_([a-f0-9]+)",
		"https://example.com/page/unsubscribe?q=$1&b=$2&c=$3&hash=$4",
		"https://example.com/page/unsubscribe?q=229&b=3430&c=2346761&hash=618e1a4097f3cf26a8b23f89b9848e5",
	},
	{
		"1252485_0_ns_74_cgfvbglsbg8uz2l1c2vwcgvaywxpy2uuaxq=@unsubscribe.example.net",
		"(\\d+)_(\\d+)_([a-z]+)_(\\d+)_([a-z0-9=]+)@(.*)",
		"http://$6/unsub.php?userid=$1&campaignid=$2&campaignname=$3&siteid=$4&hash=$5",
		"http://unsubscribe.example.net/unsub.php?userid=1252485&campaignid=0&campaignname=ns&siteid=74&hash=cgfvbglsbg8uz2l1c2vwcgvaywxpy2uuaxq=",
	},
	{
		"123_all_456_78_90_1234abcdefg@unsub.excite.example.mx",
		"(\\d+)_([a-z]+)_(\\d+)_(\\d+)_(\\d+)_([a-z0-9]+)",
		"http://unsub.excite.example.mx/list_unsubscribe.aspx?ID=$1&Type=$2&Mem=$3&CStr=$4&DStr=$5&hash=$6",
		"http://unsub.excite.example.mx/list_unsubscribe.aspx?ID=123&Type=all&Mem=456&CStr=78&DStr=90&hash=1234abcdefg",
	},
	{
		"93990-149551-4086f8f2060cc2a1f9c98f0e@unsub.example.com",
		"(\\d+)-(\\d+)-([a-f0-9]+)",
		"https://app.example.com/hooks/unsubscribe/$1/$2/$3",
		"https://app.example.com/hooks/unsubscribe/93990/149551/4086f8f2060cc2a1f9c98f0e",
	},
}

func TestPatterns(t *testing.T) {

	for _, tc := range testCases {
		re, err := regexp.Compile(tc.pat)
		if err != nil {
			t.Fatal(err)
		}

		url, err := replace(re, tc.url, tc.to)
		if err != nil {
			t.Fatal(err)
		}

		if url != tc.exp {
			t.Errorf("unexpected url: %q", url)
		}
	}
}
