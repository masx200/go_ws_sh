package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"

	)

type Session struct {
	Name string `json:"name"`

	Cmd       string    `json:"cmd"`
	Args      []string  `json:"args"`
	Dir       string    `json:"dir"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SessionStore struct {
	gorm.Model
	Name string `json:"name" gorm:"index;unique;not null"`
	Cmd  string `json:"cmd" gorm:"index;not null"`
	Args string `json:"args" gorm:"index;not null"`
	Dir  string `json:"dir" gorm:"index;not null"`
}
func (SessionStore) TableName() string {
	return strings.ToLower("SessionStore")
}

type StringSlice []string

func (s StringSlice) Value() ([]byte, error) {
	return json.Marshal(s)
}

func (s *StringSlice) Scan(value []byte) error {
	if value == nil {
		return nil
	}
	var bytes []byte = value
	return json.Unmarshal(bytes, s)
}


func CreateSessionHandler(credentialdb *gorm.DB, tokendb *gorm.DB, sessiondb *gorm.DB) func(c context.Context, r *app.RequestContext, next HertzNext) {
	return func(c context.Context, r *app.RequestContext, next HertzNext) {
		var body struct {
			Session struct {
				Name string   `json:"name"`
				Cmd  string   `json:"cmd"`
				Args []string `json:"args"`
				Dir  string   `json:"dir"`
			} `json:"session"`
			Authorization CredentialsClient `json:"authorization"`
		}

		if err := r.BindJSON(&body); err != nil {
			r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
			return
		}

		if body.Session.Name == "" || body.Session.Cmd == "" || body.Session.Dir == "" {
			r.AbortWithMsg("Error: Name is empty or  Cmd or Dir is empty ", consts.StatusBadRequest)
			return
		}

		var err error
		username := body.Authorization.Username
		if username == "" {
			username, err = GetUsernameByTokenIdentifier(tokendb, body.Authorization.Identifier)
			if err != nil {
				log.Println("Error:", err)
				r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
				return
			}
			log.Println("Username:", username)
		}

		var existingSession SessionStore
		if err := sessiondb.Where(&SessionStore{Name: body.Session.Name}).First(&existingSession).Error; err == nil {
			r.JSON(consts.StatusConflict, map[string]any{
				"message":  "Error: Session already exists",
				"username": username,
				"session": map[string]string{
					"name": body.Session.Name,
				},
			})
			return
		}

		if sessiondb.Unscoped().Where("name = ?", body.Session.Name).Delete(&SessionStore{}).Error != nil {
			log.Println("Error: Failed to delete session")
			r.AbortWithMsg("Error: Failed to delete session", consts.StatusInternalServerError)
			return
		}

		argsstringarray := StringSlice(body.Session.Args)
		var argsstring string
		argsbytes, err := argsstringarray.Value()
		if err != nil {
			r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
			return
		}
		argsstring = string(argsbytes)

		newSession := SessionStore{
			Name: body.Session.Name,
			Cmd:  body.Session.Cmd,
			Args: argsstring,
			Dir:  body.Session.Dir,
		}

		if err := sessiondb.Create(&newSession).Error; err != nil {
			r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
			return
		}

		var args []string
		if err := json.Unmarshal([]byte(newSession.Args), &args); err != nil {
			r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
			return
		}

		r.JSON(consts.StatusOK, map[string]any{
			"message":  "Session created successfully",
			"username": username,
			"session": map[string]interface{}{
				"name":     body.Session.Name,
				"cmd":      body.Session.Cmd,
				"args":     args,
				"dir":      body.Session.Dir,
				"username": username,
			},
		})
	}
}

func UpdateSessionHandler(credentialdb *gorm.DB, tokendb *gorm.DB, sessiondb *gorm.DB) func(c context.Context, r *app.RequestContext, next HertzNext) {
	return func(c context.Context, r *app.RequestContext, next HertzNext) {
		var req struct {
			Session struct {
				Name string   `json:"name"`
				Cmd  string   `json:"cmd"`
				Args []string `json:"args"`
				Dir  string   `json:"dir"`
			} `json:"session"`
			Authorization CredentialsClient `json:"authorization"`
		}

		if err := r.BindJSON(&req); err != nil {
			r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
			return
		}

		if req.Session.Name == "" {
			r.AbortWithMsg("Error: Name is empty", consts.StatusBadRequest)
			return
		}

		var err error
		username := req.Authorization.Username
		if username == "" {
			username, err = GetUsernameByTokenIdentifier(tokendb, req.Authorization.Identifier)
			if err != nil {
				log.Println("Error:", err)
				r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
				return
			}
			log.Println("Username:", username)
		}

		var session SessionStore
		if err := sessiondb.Where(&SessionStore{Name: req.Session.Name}).First(&session).Error; err != nil {
			r.JSON(consts.StatusNotFound, map[string]any{
				"message":  "Error: Session not found",
				"username": username,
				"session": map[string]string{
					"name": req.Session.Name,
				},
			})
			return
		}

		session.Cmd = req.Session.Cmd
		argsstringarray := StringSlice(req.Session.Args)
		var argsstring string

		argsbytes, err := argsstringarray.Value()
		if err != nil {
			r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
			return
		}
		argsstring = string(argsbytes)
		session.Args = argsstring
		session.Dir = req.Session.Dir
		if err := sessiondb.Save(&session).Error; err != nil {
			r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
			return
		}

		r.JSON(consts.StatusOK, map[string]any{
			"message":  "Session updated successfully",
			"username": username,
			"session": map[string]interface{}{
				"name":     req.Session.Name,
				"cmd":      req.Session.Cmd,
				"args":     req.Session.Args,
				"dir":      req.Session.Dir,
				"username": username,
			},
		})
	}
}

func DeleteSessionHandler(credentialdb *gorm.DB, tokendb *gorm.DB, sessiondb *gorm.DB, initial_sessions []Session) func(c context.Context, r *app.RequestContext, next HertzNext) {
	return func(c context.Context, r *app.RequestContext, next HertzNext) {
		var body struct {
			Session struct {
				Name string `json:"name"`
			} `json:"session"`
			Authorization CredentialsClient `json:"authorization"`
		}

		if err := r.BindJSON(&body); err != nil {
			r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
			return
		}

		if body.Session.Name == "" {
			r.AbortWithMsg("Error: Name is empty", consts.StatusBadRequest)
			return
		}

		var err error
		username := body.Authorization.Username
		if username == "" {
			username, err = GetUsernameByTokenIdentifier(tokendb, body.Authorization.Identifier)
			if err != nil {
				log.Println("Error:", err)
				r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
				return
			}
		}

		if err := sessiondb.Where("name = ?", body.Session.Name).Delete(&SessionStore{}).Error; err != nil {
			r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
			return
		}

		r.JSON(consts.StatusOK, map[string]any{
			"message":  "Session deleted successfully",
			"username": username,
		})
	}
}

func GetSessionsHandler(credentialdb *gorm.DB, tokendb *gorm.DB, sessiondb *gorm.DB) func(c context.Context, r *app.RequestContext, next HertzNext) {
	return func(c context.Context, r *app.RequestContext, next HertzNext) {
		sessiondb = sessiondb.Debug()
		var body struct {
			Authorization CredentialsClient `json:"authorization"`
			Session       struct {
				Name string `json:"name"`
			} `json:"session"`
		}

		err := r.BindJSON(&body)
		if err != nil {
			r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
			return
		}
		log.Println(body)

		username := body.Authorization.Username
		if username == "" {
			username, err = GetUsernameByTokenIdentifier(tokendb, body.Authorization.Identifier)
			if err != nil {
				log.Println("Error:", err)
				r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
				return
			}
			log.Println("Username:", username)
		}

		var sessions []Session
		if body.Session.Name != "" {
			log.Println("查询Name:", body.Session.Name)
			sessions, err = ReadAllSessionsWithName(sessiondb, body.Session.Name)
			if err != nil {
				r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
				return
			}
		} else {
			sessions, err = ReadAllSessions(sessiondb)
			if err != nil {
				r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
				return
			}
		}

		r.JSON(
			consts.StatusOK,
			map[string]interface{}{
				"message":  "List of Sessions ok",
				"sessions": SessionsToMapSlice(sessions),
				"username": username,
			},
		)
	}
}

func GenerateSessionRoutes(credentialdb *gorm.DB, tokendb *gorm.DB, sessiondb *gorm.DB, initial_sessions []Session) []RouteConfig {
	return []RouteConfig{
		{
			Headers: map[string]string{"x-HTTP-method-override": "POST"},
			Path:    "/sessions",
			Method:  "POST",
			MiddleWare: CreateSessionHandler(credentialdb, tokendb, sessiondb),
		},
		{
			Path:   "/sessions",
			Method: "PUT",
			MiddleWare: UpdateSessionHandler(credentialdb, tokendb, sessiondb),
		},
		{
			Path:   "/sessions",
			Method: "DELETE",
			MiddleWare: DeleteSessionHandler(credentialdb, tokendb, sessiondb, initial_sessions),
		},
		{
			Headers: map[string]string{"x-HTTP-method-override": "GET"},
			Path:    "/sessions",
			Method:  "POST",
			MiddleWare: GetSessionsHandler(credentialdb, tokendb, sessiondb),
		},
	}
}

func ReadAllSessions(sessiondb *gorm.DB) ([]Session, error) {
	var sessionStores []SessionStore
	if err := sessiondb.Find(&sessionStores).Error; err != nil {
		return nil, err
	}

	sessions := make([]Session, 0, len(sessionStores))
	for _, store := range sessionStores {
		var args []string
		if err := json.Unmarshal([]byte(store.Args), &args); err != nil {
			return nil, err
		}
		session := Session{
			Name:      store.Name,
			Cmd:       store.Cmd,
			Args:      args,
			Dir:       store.Dir,
			CreatedAt: store.CreatedAt,
			UpdatedAt: store.UpdatedAt,
		}
		sessions = append(sessions, session)
	}
	return sessions, nil
}

func ReadAllSessionsWithName(sessiondb *gorm.DB, name string) ([]Session, error) {
	var sessionStores []SessionStore
	if err := sessiondb.Where(
		"name = ?",
		name).Find(&sessionStores).Error; err != nil {
		return nil, err
	}

	sessions := make([]Session, 0, len(sessionStores))
	for _, store := range sessionStores {
		var args []string
		if err := json.Unmarshal([]byte(store.Args), &args); err != nil {
			return nil, err
		}
		session := Session{
			Name:      store.Name,
			Cmd:       store.Cmd,
			Args:      args,
			Dir:       store.Dir,
			CreatedAt: store.CreatedAt,
			UpdatedAt: store.UpdatedAt,
		}
		sessions = append(sessions, session)
	}
	return sessions, nil
}

func SessionsToMapSlice(sessions []Session) []map[string]any {
	result := make([]map[string]any, len(sessions))
	for i, session := range sessions {
		result[i] = map[string]any{
			"name":       session.Name,
			"cmd":        session.Cmd,
			"args":       session.Args,
			"dir":        session.Dir,
			"created_at": FormatTimeWithCarbon(carbon.CreateFromStdTime(session.CreatedAt)),
			"updated_at": FormatTimeWithCarbon(carbon.CreateFromStdTime(session.UpdatedAt)),
		}
	}
	return result
}

func CopyMiddleware(credentialdb *gorm.DB, tokendb *gorm.DB, sessiondb *gorm.DB, c context.Context, r *app.RequestContext, next HertzNext) {
	var body struct {
		Session struct {
			Name string `json:"name"`
		} `json:"session"`
		Authorization CredentialsClient `json:"authorization"`
		Destination   struct {
			Name string `json:"name"`
		} `json:"destination"`
	}

	if err := r.BindJSON(&body); err != nil {
		r.AbortWithMsg("Error: "+err.Error(), consts.StatusBadRequest)
		return
	}

	var err error
	username := body.Authorization.Username
	if username == "" {
		username, err = GetUsernameByTokenIdentifier(tokendb, body.Authorization.Identifier)
		if err != nil {
			log.Println("Error:", err)
			r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
			return
		}
		log.Println("Username:", username)
	}

	var existingSession SessionStore
	if err := sessiondb.Where(&SessionStore{Name: body.Destination.Name}).First(&existingSession).Error; err == nil {
		r.JSON(consts.StatusConflict, map[string]any{
			"message":  "Error: Session already exists",
			"username": username,
			"session": map[string]string{
				"name": body.Session.Name,
			},
		})
		return
	}

	var newSession *SessionStore
	if newSession, err = CopySession(sessiondb, body.Session.Name, body.Destination.Name); err != nil {
		log.Printf("Failed to copy session: %v", err)
		r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
		return
	}

	var args []string
	if err := json.Unmarshal([]byte(newSession.Args), &args); err != nil {
		r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
		return
	}

	r.JSON(consts.StatusOK, map[string]any{
		"message":  "Session copied successfully",
		"username": username,
		"session": map[string]interface{}{
			"name":     newSession.Name,
			"cmd":      newSession.Cmd,
			"args":     args,
			"dir":      newSession.Dir,
			"username": username,
		},
	})
}

func MoveMiddleware(initial_sessions []Session, credentialdb *gorm.DB, tokendb *gorm.DB, sessiondb *gorm.DB, c context.Context, r *app.RequestContext, next HertzNext) {
	var body struct {
		Session struct {
			Name string `json:"name"`
		} `json:"session"`
		Authorization CredentialsClient `json:"authorization"`
		Destination   struct {
			Name string `json:"name"`
		} `json:"destination"`
	}

	if err := r.BindJSON(&body); err != nil {
		r.AbortWithMsg("Error: "+err.Error(), consts.StatusBadRequest)
		return
	}

	for _, session := range initial_sessions {
		if session.Name == body.Session.Name {
			r.AbortWithMsg("Error: Session is initial session,不允许删除", consts.StatusBadRequest)
			return
		}
	}

	var err error
	username := body.Authorization.Username
	if username == "" {
		username, err = GetUsernameByTokenIdentifier(tokendb, body.Authorization.Identifier)
		if err != nil {
			log.Println("Error:", err)
			r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
			return
		}
		log.Println("Username:", username)
	}

	var existingSession SessionStore
	if err := sessiondb.Where(&SessionStore{Name: body.Destination.Name}).First(&existingSession).Error; err == nil {
		r.JSON(consts.StatusConflict, map[string]any{
			"message":  "Error: Session already exists",
			"username": username,
			"session": map[string]string{
				"name": body.Session.Name,
			},
		})
		return
	}

	var newSession *SessionStore
	if newSession, err = MoveSession(sessiondb, body.Session.Name, body.Destination.Name); err != nil {
		log.Printf("Failed to move session: %v", err)
		r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
		return
	}

	var args []string
	if err := json.Unmarshal([]byte(newSession.Args), &args); err != nil {
		r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
		return
	}

	r.JSON(consts.StatusOK, map[string]any{
		"message":  "Session moved successfully",
		"username": username,
		"session": map[string]interface{}{
			"name":     newSession.Name,
			"cmd":      newSession.Cmd,
			"args":     args,
			"dir":      newSession.Dir,
			"username": username,
		},
	})
}

func CopySession(sessiondb *gorm.DB, sourceName, destName string) (*SessionStore, error) {
	var sourceSession SessionStore
	if err := sessiondb.Where("name = ?", sourceName).First(&sourceSession).Error; err != nil {
		return nil, fmt.Errorf("source session not found: %v", err)
	}

	newSession := SessionStore{
		Name: destName,
		Cmd:  sourceSession.Cmd,
		Args: sourceSession.Args,
		Dir:  sourceSession.Dir,
	}

	if err := sessiondb.Create(&newSession).Error; err != nil {
		return nil, fmt.Errorf("failed to create session: %v", err)
	}

	return &newSession, nil
}

func MoveSession(sessiondb *gorm.DB, sourceName, destName string) (*SessionStore, error) {
	var sourceSession SessionStore
	if err := sessiondb.Where("name = ?", sourceName).First(&sourceSession).Error; err != nil {
		return nil, fmt.Errorf("source session not found: %v", err)
	}

	sourceSession.Name = destName

	if err := sessiondb.Save(&sourceSession).Error; err != nil {
		return nil, fmt.Errorf("failed to update session: %v", err)
	}

	return &sourceSession, nil
}