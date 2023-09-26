package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
	RouteParams []string
	Received    tgbotapi.Update
	Route       *Route
	Message     string
}

type Router struct {
	Routes []*Route
	cache  string // todo cache driver
}

func NewRouter() *Router {
	return &Router{}
}

func (router *Router) AddRoute(path string, handler func(ctx *RequestContext) (Renderable, error), name string) {
	pattern, vars := parsePattern(path)
	router.Routes = append(router.Routes, &Route{Path: path, Pattern: pattern, Handler: handler, variables: vars, Name: name})
}

func (router *Router) MatchRoute(received tgbotapi.Update) (*Route, string, []string, error) {
	message := received.Message
	if router.onlyTextMessage(message) {
		found := router.foundHandler(message.Text)
		if found == nil {
			found = router.handleLastStep(received.Message.From)
			if found == nil {
				return nil, "", nil, RouteNotFoundError{message.Text}
			}
			target := regexp.MustCompile(found.Pattern)
			routeParams := extractDynamicVars(received.CallbackQuery.Data, target)
			return found, message.Text, routeParams, nil
		} else {
			return found, message.Text, nil, nil
		}
	}
	if received.CallbackQuery != nil {
		found := router.foundHandler(received.CallbackQuery.Data)
		if found == nil {
			return nil, "", nil, RouteNotFoundError{received.CallbackQuery.Data}
		}
		target := regexp.MustCompile(found.Pattern)
		routeParams := extractDynamicVars(received.CallbackQuery.Data, target)

		return found, received.CallbackQuery.Data, routeParams, nil
	}
	return nil, "", nil, RouteNotFoundError{message.Text}
}

func (router *Router) handleLastStep(user *tgbotapi.User) *Route {
	//TODO redis cache
	pattern := "last/step"
	return router.foundHandler(pattern)
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
