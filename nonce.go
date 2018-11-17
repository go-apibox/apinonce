package apinonce

import (
	"net/http"
	"time"

	"github.com/go-apibox/api"
	"github.com/go-apibox/cache"
	"github.com/go-apibox/utils"
)

type Nonce struct {
	app      *api.App
	disabled bool
	inited   bool

	cache         *cache.Cache
	length        int
	expireTime    int
	maxCacheCount int
	actionMatcher *utils.Matcher
}

func NewNonce(app *api.App) *Nonce {
	app.Error.RegisterGroupErrors("nonce", ErrorDefines)

	nonce := new(Nonce)
	nonce.app = app

	cfg := app.Config
	disabled := cfg.GetDefaultBool("apinonce.disabled", false)
	nonce.disabled = disabled
	if disabled {
		return nonce
	}

	nonce.init()
	return nonce
}

func (n *Nonce) init() {
	if n.inited {
		return
	}

	app := n.app
	cfg := app.Config
	nonceLen := cfg.GetDefaultInt("apinonce.length", 16)
	expireTime := cfg.GetDefaultInt("apinonce.expire_time", 1000)
	maxcacheCount := cfg.GetDefaultInt("apinonce.max_cache_count", 100000)
	actionWhitelist := cfg.GetDefaultStringArray("apinonce.actions.whitelist", []string{"*"})
	actionBlacklist := cfg.GetDefaultStringArray("apinonce.actions.blacklist", []string{})

	matcher := utils.NewMatcher()
	matcher.SetWhiteList(actionWhitelist)
	matcher.SetBlackList(actionBlacklist)

	cache := cache.NewCache(time.Duration(expireTime) * time.Second)

	n.cache = cache
	n.length = nonceLen
	n.expireTime = expireTime
	n.maxCacheCount = maxcacheCount
	n.actionMatcher = matcher
	n.inited = true
}

func (n *Nonce) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if n.disabled {
		next(w, r)
		return
	}

	c, err := api.NewContext(n.app, w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// check if action not required nonce check
	action := c.Input.GetAction()
	if !n.actionMatcher.Match(action) {
		next(w, r)
		return
	}

	nonce := c.Input.Get("api_nonce")
	if nonce == "" {
		api.WriteResponse(c, n.app.Error.NewGroupError("nonce", errorMissingNonce))
		return
	}
	if len(nonce) != n.length {
		api.WriteResponse(c, n.app.Error.NewGroupError("nonce", errorInvalidNonce))
		return
	}
	_, exists := n.cache.Get(nonce)
	if exists {
		api.WriteResponse(c, n.app.Error.NewGroupError("nonce", errorNonceExist))
		return
	}

	// 检查数量是否超出限制
	if n.cache.Count() >= n.maxCacheCount {
		api.WriteResponse(c, n.app.Error.NewGroupError("nonce", errorNonceCountExceed))
		return
	}

	n.cache.Set(nonce, true)

	// next middleware
	next(w, r)
}

// Enable enable the middle ware.
func (n *Nonce) Enable() {
	n.disabled = false
	n.init()
}

// Disable disable the middle ware.
func (n *Nonce) Disable() {
	n.disabled = true
}
