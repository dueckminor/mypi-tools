package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	texttemplate "text/template"

	"github.com/dueckminor/mypi-tools/go/gotty/utils"
	"github.com/dueckminor/mypi-tools/go/gotty/webtty"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

type InitMessage struct {
	Arguments string `json:"Arguments,omitempty"`
	AuthToken string `json:"AuthToken,omitempty"`
}

type Server struct {
	Factory       Factory
	Opt           *Options
	titleTemplate *texttemplate.Template
}

func NewServer(f Factory) (s *Server, err error) {
	s = &Server{
		Opt:     &Options{},
		Factory: f,
	}

	err = utils.ApplyDefaultValues(s.Opt)
	if err != nil {
		return nil, err
	}
	s.Opt.PermitWrite = true
	s.titleTemplate, err = texttemplate.New("title").Parse(s.Opt.TitleFormat)
	return s, err
}

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	Subprotocols:    webtty.Protocols,
}

func Handler(ctx context.Context, w http.ResponseWriter, r *http.Request, f Factory) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to set websocket upgrade: %V", err)
		return
	}

	srv, err := NewServer(f)
	if err != nil {
		fmt.Println("Failed to set websocket upgrade: %V", err)
		return
	}
	err = srv.ProcessWSConn(ctx, conn)
	if err != nil {
		fmt.Println("Failed to process websocket: %V", err)
	}
}

func (server *Server) titleVariables(order []string, varUnits map[string]map[string]interface{}) map[string]interface{} {
	titleVars := map[string]interface{}{}

	for _, name := range order {
		vars, ok := varUnits[name]
		if !ok {
			panic("title variable name error")
		}
		for key, val := range vars {
			titleVars[key] = val
		}
	}

	// safe net for conflicted keys
	for _, name := range order {
		titleVars[name] = varUnits[name]
	}

	return titleVars
}

func (server *Server) ProcessWSConn(ctx context.Context, conn *websocket.Conn) error {
	typ, initLine, err := conn.ReadMessage()
	if err != nil {
		return errors.Wrapf(err, "failed to authenticate websocket connection: invalid/missing init message")
	}
	if typ != websocket.TextMessage {
		return errors.New("failed to authenticate websocket connection: invalid message type")
	}

	var init InitMessage
	err = json.Unmarshal(initLine, &init)
	if err != nil {
		return errors.Wrapf(err, "failed to authenticate websocket connection")
	}
	if init.AuthToken != server.Opt.Credential {
		return errors.New("failed to authenticate websocket connection")
	}

	var slave Slave
	slave, err = server.Factory.New(nil /*params*/)
	if err != nil {
		return errors.Wrapf(err, "failed to create backend")
	}
	defer slave.Close()

	titleVars := server.titleVariables(
		[]string{"server", "master", "slave"},
		map[string]map[string]interface{}{
			"server": server.Opt.TitleVariables,
			"master": map[string]interface{}{
				"remote_addr": conn.RemoteAddr(),
			},
			"slave": slave.WindowTitleVariables(),
		},
	)

	titleBuf := new(bytes.Buffer)
	err = server.titleTemplate.Execute(titleBuf, titleVars)
	if err != nil {
		return errors.Wrapf(err, "failed to fill window title template")
	}

	opts := []webtty.Option{
		webtty.WithWindowTitle(titleBuf.Bytes()),
	}
	if server.Opt.PermitWrite {
		opts = append(opts, webtty.WithPermitWrite())
	}
	if server.Opt.EnableReconnect {
		opts = append(opts, webtty.WithReconnect(server.Opt.ReconnectTime))
	}
	if server.Opt.Width > 0 {
		opts = append(opts, webtty.WithFixedColumns(server.Opt.Width))
	}
	if server.Opt.Height > 0 {
		opts = append(opts, webtty.WithFixedRows(server.Opt.Height))
	}
	if server.Opt.Preferences != nil {
		opts = append(opts, webtty.WithMasterPreferences(server.Opt.Preferences))
	}

	tty, err := webtty.New(&wsWrapper{conn}, slave, opts...)
	if err != nil {
		return errors.Wrapf(err, "failed to create webtty")
	}

	err = tty.Run(ctx)

	return err
}
