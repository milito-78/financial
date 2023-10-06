package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	"net/url"
	"regexp"
	"strings"
)

type Route struct {
	variables []string
	Handler   func(ctx *RequestContext) (Renderable, error)
	Pattern   string
	Path      string
	Name      string
}

type RequestContext struct {
	RouteParams       []string
	QueryParams       url.Values
	Received          tgbotapi.Update
	Route             *Route
	Message           string
	stateSaverChannel chan<- State
	lastState         *State
}

func (r *RequestContext) SetState(state State) {
	r.stateSaverChannel <- state
}

func (r *RequestContext) GetState() *State {
	return r.lastState
}

type Router struct {
	Routes    []*Route
	LastState func(lastStateUser string) (string, error)
}

func NewRouter() *Router {
	return &Router{}
}

func (router *Router) AddRoute(path string, handler func(ctx *RequestContext) (Renderable, error), name string) {
	pattern, vars := parsePattern(path)
	router.Routes = append(router.Routes, &Route{Path: path, Pattern: pattern, Handler: handler, variables: vars, Name: name})
}

func (router *Router) MatchRoute(path string) (*Route, *url.Values, []string, error) {
	found := router.foundHandler(path)
	if found != nil {
		target := regexp.MustCompile(found.Pattern)
		routeParams := extractDynamicVars(path, target)
		queryParams := extractQueryParams(path)
		return found, queryParams, routeParams, nil
	}
	return nil, &url.Values{}, nil, RouteNotFoundError{path}
}

func (router *Router) onlyTextMessage(message *tgbotapi.Message) bool {
	return message.Text != "" && message.Photo == nil && message.Document == nil && message.Video == nil
}

func (router *Router) foundHandler(pattern string) *Route {
	for _, r := range router.Routes {
		target := regexp.MustCompile(r.Pattern)
		if target.MatchString(pattern) {
			return r
		}
	}
	return nil
}

func extractDynamicVars(input string, target *regexp.Regexp) (vars []string) {
	matches := target.FindStringSubmatch(input)
	if len(matches) != 0 {
		vars = matches[1:]
	}

	return vars
}

func parsePattern(pattern string) (path string, vars []string) {
	dynamicSegmentRegex := regexp.MustCompile(`:[a-zA-Z0-9_]+`)
	replacedPattern := dynamicSegmentRegex.ReplaceAllString(pattern, `([^/]+)`)
	return replacedPattern, extractDynamicSegments(pattern, dynamicSegmentRegex)
}

func extractDynamicSegments(inputPattern string, dynamicSegmentRegex *regexp.Regexp) []string {
	matches := dynamicSegmentRegex.FindAllString(inputPattern, -1)
	result := make([]string, 0)

	for _, match := range matches {
		key := strings.TrimPrefix(match, ":")
		result = append(result, key)
	}

	return result
}

func extractQueryParams(path string) *url.Values {
	parsedURL, err := url.Parse(path)
	if err != nil {
		log.Error("Error parsing URL:", err)
		return nil
	}
	t := parsedURL.Query()
	return &t
}
