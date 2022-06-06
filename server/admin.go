// Copyright (C) 2019-2022 Chrystian Huot <chrystian.huot@saubeo.solutions>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"
)

type Admin struct {
	Attempts         AdminLoginAttempts
	AttemptsMax      uint
	AttemptsMaxDelay time.Duration
	Broadcast        chan *[]byte
	Conns            map[*websocket.Conn]bool
	Controller       *Controller
	Register         chan *websocket.Conn
	Tokens           []string
	Unregister       chan *websocket.Conn
	mutex            sync.Mutex
	running          bool
}

type AdminLoginAttempt struct {
	Count uint
	Date  time.Time
}

type AdminLoginAttempts map[string]*AdminLoginAttempt

func NewAdmin(controller *Controller) *Admin {
	return &Admin{
		Attempts:         AdminLoginAttempts{},
		AttemptsMax:      uint(3),
		AttemptsMaxDelay: time.Duration(time.Duration.Minutes(10)),
		Broadcast:        make(chan *[]byte),
		Conns:            make(map[*websocket.Conn]bool),
		Controller:       controller,
		Register:         make(chan *websocket.Conn),
		Tokens:           []string{},
		Unregister:       make(chan *websocket.Conn),
		mutex:            sync.Mutex{},
	}
}

func (admin *Admin) BroadcastConfig() {
	if b, err := json.Marshal(admin.GetConfig()); err == nil {
		for conn := range admin.Conns {
			conn.WriteMessage(websocket.TextMessage, b)
		}
	}
}

func (admin *Admin) ChangePassword(currentPassword interface{}, newPassword string) error {
	var (
		err  error
		hash []byte
	)

	if len(newPassword) == 0 {
		return errors.New("newPassword is empty")
	}

	switch v := currentPassword.(type) {
	case string:
		if err = bcrypt.CompareHashAndPassword([]byte(admin.Controller.Options.adminPassword), []byte(v)); err != nil {
			return err
		}
	}

	if hash, err = bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost); err != nil {
		return err
	}

	admin.Controller.Options.adminPassword = string(hash)
	admin.Controller.Options.adminPasswordNeedChange = newPassword == defaults.adminPassword

	if err := admin.Controller.Options.Write(admin.Controller.Database); err != nil {
		return err
	}

	if err := admin.Controller.Options.Read(admin.Controller.Database); err != nil {
		return err
	}

	admin.Controller.Logs.LogEvent(LogLevelWarn, "admin password changed.")

	return nil
}

func (admin *Admin) ConfigHandler(w http.ResponseWriter, r *http.Request) {
	if strings.EqualFold(r.Header.Get("upgrade"), "websocket") {
		upgrader := websocket.Upgrader{}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		admin.Register <- conn

		go func() {
			conn.SetReadDeadline(time.Time{})

			for {
				_, b, err := conn.ReadMessage()
				if err != nil {
					break
				}

				if !admin.ValidateToken(string(b)) {
					break
				}
			}

			admin.Unregister <- conn

			conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(1000, ""))
		}()

	} else {
		logError := func(err error) {
			admin.Controller.Logs.LogEvent(LogLevelError, fmt.Sprintf("admin.confighandler.put: %s", err.Error()))
		}

		t := admin.GetAuthorization(r)
		if !admin.ValidateToken(t) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		switch r.Method {
		case http.MethodGet:
			admin.SendConfig(w)

		case http.MethodPut:
			m := map[string]interface{}{}
			err := json.NewDecoder(r.Body).Decode(&m)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			admin.Controller.IngestLock()
			admin.mutex.Lock()

			admin.Controller.Dirwatches.Stop()

			switch v := m["access"].(type) {
			case []interface{}:
				admin.Controller.Accesses.FromMap(v)
				err := admin.Controller.Accesses.Write(admin.Controller.Database)
				if err != nil {
					logError(err)
				} else {
					err = admin.Controller.Accesses.Read(admin.Controller.Database)
					if err != nil {
						logError(err)
					}
				}
			}

			switch v := m["apiKeys"].(type) {
			case []interface{}:
				admin.Controller.Apikeys.FromMap(v)
				err = admin.Controller.Apikeys.Write(admin.Controller.Database)
				if err != nil {
					logError(err)
				} else {
					err = admin.Controller.Apikeys.Read(admin.Controller.Database)
					if err != nil {
						logError(err)
					}
				}
			}

			switch v := m["dirWatch"].(type) {
			case []interface{}:
				admin.Controller.Dirwatches.FromMap(v)
				err = admin.Controller.Dirwatches.Write(admin.Controller.Database)
				if err != nil {
					logError(err)
				} else {
					err = admin.Controller.Dirwatches.Read(admin.Controller.Database)
					if err != nil {
						logError(err)
					}
				}
			}

			switch v := m["downstreams"].(type) {
			case []interface{}:
				admin.Controller.Downstreams.FromMap(v)
				err = admin.Controller.Downstreams.Write(admin.Controller.Database)
				if err != nil {
					logError(err)
				} else {
					err = admin.Controller.Downstreams.Read(admin.Controller.Database)
					if err != nil {
						logError(err)
					}
				}
			}

			switch v := m["groups"].(type) {
			case []interface{}:
				admin.Controller.Groups.FromMap(v)
				err = admin.Controller.Groups.Write(admin.Controller.Database)
				if err != nil {
					logError(err)
				} else {
					err = admin.Controller.Groups.Read(admin.Controller.Database)
					if err != nil {
						logError(err)
					}
				}
			}

			switch v := m["options"].(type) {
			case map[string]interface{}:
				admin.Controller.Options.FromMap(v)
				err = admin.Controller.Options.Write(admin.Controller.Database)
				if err != nil {
					logError(err)
				}
			}

			switch v := m["systems"].(type) {
			case []interface{}:
				admin.Controller.Systems.FromMap(v)
				err = admin.Controller.Systems.Write(admin.Controller.Database)
				if err != nil {
					logError(err)
				} else {
					err = admin.Controller.Systems.Read(admin.Controller.Database)
					if err != nil {
						logError(err)
					}
				}
			}

			switch v := m["tags"].(type) {
			case []interface{}:
				admin.Controller.Tags.FromMap(v)
				err = admin.Controller.Tags.Write(admin.Controller.Database)
				if err != nil {
					logError(err)
				} else {
					err = admin.Controller.Tags.Read(admin.Controller.Database)
					if err != nil {
						logError(err)
					}
				}
			}

			admin.mutex.Unlock()
			admin.Controller.IngestUnlock()

			admin.Controller.EmitConfig()
			admin.Controller.Dirwatches.Start(admin.Controller)

			admin.SendConfig(w)

			admin.Controller.Logs.LogEvent(LogLevelWarn, "configuration changed")

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func (admin *Admin) GetAuthorization(r *http.Request) string {
	return r.Header.Get("Authorization")
}

func (admin *Admin) GetConfig() map[string]interface{} {
	systems := []map[string]interface{}{}
	for _, system := range admin.Controller.Systems.List {
		systems = append(systems, map[string]interface{}{
			"_id":          system.RowId,
			"autoPopulate": system.AutoPopulate,
			"blacklists":   system.Blacklists,
			"id":           system.Id,
			"label":        system.Label,
			"led":          system.Led,
			"order":        system.Order,
			"talkgroups":   system.Talkgroups.List,
			"units":        system.Units.List,
		})
	}

	return map[string]interface{}{
		"access":      admin.Controller.Accesses.List,
		"apiKeys":     admin.Controller.Apikeys.List,
		"dirWatch":    admin.Controller.Dirwatches.List,
		"downstreams": admin.Controller.Downstreams.List,
		"groups":      admin.Controller.Groups.List,
		"options":     admin.Controller.Options,
		"systems":     systems,
		"tags":        admin.Controller.Tags.List,
	}
}

func (admin *Admin) LogsHandler(w http.ResponseWriter, r *http.Request) {
	t := admin.GetAuthorization(r)
	if !admin.ValidateToken(t) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	switch r.Method {
	case http.MethodPost:
		m := map[string]interface{}{}
		err := json.NewDecoder(r.Body).Decode(&m)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		logOptions := LogsSearchOptions{}
		err = logOptions.FromMap(m)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		r, err := admin.Controller.Logs.Search(&logOptions, admin.Controller.Database)
		if err != nil {
			admin.Controller.Logs.LogEvent(LogLevelError, err.Error())
			w.WriteHeader(http.StatusExpectationFailed)
			return
		}

		b, err := json.Marshal(r)
		if err != nil {
			admin.Controller.Logs.LogEvent(LogLevelError, err.Error())
			w.WriteHeader(http.StatusExpectationFailed)
			return
		}

		w.Write(b)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (admin *Admin) LoginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		m := map[string]interface{}{}

		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		remoteAddr := GetRemoteAddr(r)

		attempt := admin.Attempts[remoteAddr]

		if attempt == nil {
			admin.Attempts[remoteAddr] = &AdminLoginAttempt{
				Count: 1,
				Date:  time.Now(),
			}
			attempt = admin.Attempts[remoteAddr]
		} else {
			attempt.Count++
			attempt.Date = time.Now()
		}

		if attempt.Count > admin.AttemptsMax || time.Since(attempt.Date) < admin.AttemptsMaxDelay {
			if attempt.Count == admin.AttemptsMax+1 {
				admin.Controller.Logs.LogEvent(LogLevelWarn, fmt.Sprintf("too many login attempts for ip=\"%v\"", remoteAddr))
			}

			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ok := false

		switch v := m["password"].(type) {
		case string:
			if len(v) > 0 {
				if err := bcrypt.CompareHashAndPassword([]byte(admin.Controller.Options.adminPassword), []byte(v)); err == nil {
					ok = true
				}
			}
		}

		if !ok {
			admin.Controller.Logs.LogEvent(LogLevelWarn, fmt.Sprintf("invalid login attempt for ip=%v", remoteAddr))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		id, err := uuid.NewRandom()

		if err != nil {
			w.WriteHeader(http.StatusExpectationFailed)
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{ID: id.String()})
		sToken, err := token.SignedString([]byte(admin.Controller.Options.secret))

		if err != nil {
			w.WriteHeader(http.StatusExpectationFailed)
			return
		}

		if len(admin.Tokens) < 5 {
			admin.Tokens = append(admin.Tokens, sToken)
		} else {
			admin.Tokens = append(admin.Tokens[1:], sToken)
		}

		b, err := json.Marshal(map[string]interface{}{
			"passwordNeedChange": true,
			"token":              sToken,
		})
		if err != nil {
			w.WriteHeader(http.StatusExpectationFailed)
			return
		}

		for k, v := range admin.Attempts {
			if time.Since(v.Date) > admin.AttemptsMaxDelay {
				delete(admin.Attempts, k)
			}
		}

		w.Write(b)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (admin *Admin) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		t := admin.GetAuthorization(r)
		if !admin.ValidateToken(t) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		for k, v := range admin.Tokens {
			if v == t {
				admin.Tokens = append(admin.Tokens[:k], admin.Tokens[k+1:]...)
			}
		}
		w.WriteHeader(http.StatusOK)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (admin *Admin) PasswordHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var (
			b               []byte
			currentPassword interface{}
			newPassword     string
		)

		logError := func(err error) {
			admin.Controller.Logs.LogEvent(LogLevelError, fmt.Sprintf("admin.passwordhandler.post: %s", err.Error()))
		}

		t := admin.GetAuthorization(r)
		if !admin.ValidateToken(t) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		m := map[string]interface{}{}
		err := json.NewDecoder(r.Body).Decode(&m)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		switch v := m["currentPassword"].(type) {
		case string:
			currentPassword = v
		}

		switch v := m["newPassword"].(type) {
		case string:
			newPassword = v
		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err = admin.ChangePassword(currentPassword, newPassword); err != nil {
			logError(errors.New("unable to change admin password, current password is invalid"))
			w.WriteHeader(http.StatusExpectationFailed)
			return
		}

		if b, err = json.Marshal(map[string]interface{}{"passwordNeedChange": admin.Controller.Options.adminPasswordNeedChange}); err == nil {
			w.Write(b)
		} else {
			w.WriteHeader(http.StatusExpectationFailed)
		}

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (admin *Admin) SendConfig(w http.ResponseWriter) {
	var m map[string]interface{}
	_, docker := os.LookupEnv("DOCKER")
	if docker {
		m = map[string]interface{}{
			"config":             admin.GetConfig(),
			"docker":             docker,
			"passwordNeedChange": admin.Controller.Options.adminPasswordNeedChange,
		}
	} else {
		m = map[string]interface{}{
			"config":             admin.GetConfig(),
			"passwordNeedChange": admin.Controller.Options.adminPasswordNeedChange,
		}
	}
	if b, err := json.Marshal(m); err == nil {
		w.Write(b)
	} else {
		w.WriteHeader(http.StatusExpectationFailed)
	}
}

func (admin *Admin) Start() error {
	if admin.running {
		return errors.New("admin already running")
	} else {
		admin.running = true
	}

	go func() {
		for {
			select {
			case data, ok := <-admin.Broadcast:
				if !ok {
					return
				}

				for conn := range admin.Conns {
					err := conn.WriteMessage(websocket.TextMessage, *data)
					if err != nil {
						admin.Unregister <- conn
					}
				}

			case conn := <-admin.Register:
				admin.Conns[conn] = true

			case conn := <-admin.Unregister:
				if _, ok := admin.Conns[conn]; ok {
					delete(admin.Conns, conn)
					conn.Close()
				}
			}
		}
	}()

	return nil
}

func (admin *Admin) UserAddHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		logError := func(err error) {
			admin.Controller.Logs.LogEvent(LogLevelError, fmt.Sprintf("admin.useraddhandler.post: %s", err.Error()))
		}

		t := admin.GetAuthorization(r)
		if !admin.ValidateToken(t) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		m := map[string]interface{}{}
		err := json.NewDecoder(r.Body).Decode(&m)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		admin.Controller.Accesses.Add(NewAccess().FromMap(m))

		if err := admin.Controller.Accesses.Write(admin.Controller.Database); err == nil {
			if err := admin.Controller.Accesses.Read(admin.Controller.Database); err == nil {
				admin.BroadcastConfig()
				w.WriteHeader(http.StatusOK)
			} else {
				logError(err)
				w.WriteHeader(http.StatusExpectationFailed)
			}
		} else {
			logError(err)
			w.WriteHeader(http.StatusExpectationFailed)
		}

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (admin *Admin) UserRemoveHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		logError := func(err error) {
			admin.Controller.Logs.LogEvent(LogLevelError, fmt.Sprintf("admin.userremovehandler.post: %s", err.Error()))
		}

		t := admin.GetAuthorization(r)
		if !admin.ValidateToken(t) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		m := map[string]interface{}{}
		err := json.NewDecoder(r.Body).Decode(&m)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if _, ok := admin.Controller.Accesses.Remove(NewAccess().FromMap(m)); ok {
			if err := admin.Controller.Accesses.Write(admin.Controller.Database); err == nil {
				if err := admin.Controller.Accesses.Read(admin.Controller.Database); err == nil {
					admin.BroadcastConfig()
					w.WriteHeader(http.StatusOK)
				} else {
					logError(err)
					w.WriteHeader(http.StatusExpectationFailed)
				}
			} else {
				logError(err)
				w.WriteHeader(http.StatusExpectationFailed)
			}
		}

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (admin *Admin) ValidateToken(sToken string) bool {
	found := false
	for _, t := range admin.Tokens {
		if t == sToken {
			found = true
			break
		}
	}
	if !found {
		return false
	}

	token, err := jwt.Parse(sToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(admin.Controller.Options.secret), nil
	})
	if err != nil {
		return false
	}

	return token.Valid
}
