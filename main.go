package main

import (
	"bufio"
	"bytes"
	"flag"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	logFn := flag.String("log", "chat.log", "filename for logging; destructively created each run")
	uname := flag.String("uname", "default_guy_123", "username")
	host := flag.String("host", "localhost", "hostname of stomp server")
	port := flag.Int("port", 32801, "port number for stomp server")

	flag.Parse()

	lf, err := os.Create(*logFn)
	if err != nil {
		panic(err)
	}

	log.SetOutput(lf)

	c := NewClient("main", *uname, *host, *port)
	c.Start()
}

type client struct {
	username string
	input    *tview.InputField
	view     *tview.TextView
	app      *tview.Application
	layout   *tview.Flex
	updates  chan string
	buffer   []byte
	conn     net.Conn
}

func NewClient(initial, username, host string, port int) *client {
	conn, err := net.Dial("tcp", host+":"+strconv.Itoa(port))
	if err != nil {
		panic(err)
	}
	return &client{
		username: username,
		view:     tview.NewTextView().SetText("Lorem ipsum"),
		input:    tview.NewInputField(),
		layout:   tview.NewFlex().SetDirection(tview.FlexRow),
		updates:  make(chan string),
		buffer:   make([]byte, 300),
		conn:     conn,
	}
}

func (c *client) Start() {
	c.app = tview.NewApplication()
	c.app.SetBeforeDrawFunc(func(screen tcell.Screen) bool {
		screen.Clear()
		return false
	})

	subh := make(map[string]string)
	subh["destination"] = "/channel/main"
	subh["id"] = "1"

	subFr := Frame{
		Command: SUBSCRIBE,
		Headers: subh,
		Body:    "",
	}

	subFrS := UnmarshalFrame(subFr)
	_, err := c.conn.Write([]byte(subFrS))
	if err != nil {
		panic(err)
	}

	c.input.SetLabel(c.username + " >")

	c.layout.AddItem(c.view, 0, 99, false).
		AddItem(c.input, 0, 1, true)

	go func() {
		log.Println("starting msg receive gofunc")
		for m := range c.updates {
			log.Println("msg received")
			c.app.QueueUpdateDraw(func() {
				mBytes := []byte(m)
				c.buffer = append(c.buffer, mBytes...)

				c.view.Clear()
				_, _ = c.view.Write(c.buffer)

				log.Printf("msg written: %s", m)

				c.view.ScrollToEnd()
			})
		}
		log.Println("closing msg receive gofunc")
	}()

	go c.Read()

	c.input.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			c.Send(c.username, c.input.GetText()+"\n")
			c.input.SetText("")
		}
	})

	if err := c.app.SetRoot(c.layout, true).Run(); err != nil {
		panic(err)
	}
}

func (c *client) MessageReceiver(m string) {
	c.updates <- m
}

func (c *client) Send(u, m string) {
	ts := string(time.Now().Format(time.StampMilli))
	h := make(map[string]string)
	h["destination"] = "/channel/main"
	fr := Frame{
		Command: SEND,
		Headers: h,
		Body:    ts + " " + u + "> " + m,
	}

	s := UnmarshalFrame(fr)
	_, err := c.conn.Write([]byte(s))
	if err != nil {
		panic(err)
	}

	log.Printf("sent message: %s > %s", u, m)
}

func (c *client) Read() {
	scanner := bufio.NewScanner(c.conn)
	scanner.Split(ScanNullTerm)
	for {
		if ok := scanner.Scan(); !ok {
			break
		}
		txt := scanner.Text()
		if txt != "\n" {
			//if it's not just a heartbeat
			fr, err := ParseFrame(txt + "\000")
			if err != nil {
				log.Printf("malformed frame err %s: %v\n", err, []byte(txt))
			}
			if fr.Command == MESSAGE {
				c.updates <- fr.Body
			}
		}
	}
}

// Custom scanner to split incoming stream on \000
func ScanNullTerm(data []byte, atEOF bool) (int, []byte, error) {
	// if we're at EOF, we're done for now
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	// if we find a '\000', return the data up to and including that index
	if i := bytes.IndexByte(data, '\000'); i >= 0 {
		// there is a null-terminated frame
		return i + 1, data[0:i], nil
	}

	// if we did not find a null-terminated frame, still need to check for \n
	if len(data) > 0 && data[0] == '\n' {
		return 1, []byte{data[0]}, nil
	}

	// if we are at EOF and we have data, return it so we can see what's going on
	if atEOF {
		return len(data), data, nil
	}
	return 0, nil, nil
}
