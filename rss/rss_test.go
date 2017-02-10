package rss

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

var rss_example = ` <?xml version="1.0" encoding="UTF-8" ?>
<rss version="2.0">

<channel>
  <title>ChannelTitle</title>
  <link>http://www.channel.com</link>
  <description>ChannelDescription</description>
  <item>
    <title>ItemTitle1</title>
    <link>http://www.ItemLink1.com</link>
    <description>ItemDescription1</description>
  </item>
  <item>
    <title>ItemTitle2</title>
    <link>http://www.ItemLink2.com</link>
    <description>ItemDescription2</description>
  </item>
</channel>

</rss>
`

func TestRss_Get(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, rss_example)
	}))
	defer ts.Close()

	channel, err := Get(ts.URL)

	if err != nil {
		t.Fatalf("%s", err)
	}

	expected := &Channel{
		Title:       "ChannelTitle",
		Description: "ChannelDescription",
		Link:        "http://www.channel.com",
		Items: []*Item{
			&Item{
				Title:       "ItemTitle1",
				Link:        "http://www.ItemLink1.com",
				Description: "ItemDescription1",
			},
			&Item{
				Title:       "ItemTitle2",
				Link:        "http://www.ItemLink2.com",
				Description: "ItemDescription2",
			},
		},
	}

	if !reflect.DeepEqual(channel, expected) {
		t.Errorf("expected \n%#v\n to be \n%#v\n", channel, expected)
	}
}
