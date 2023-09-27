package telegram

import (
	"financial/application"
	"github.com/golobby/container/v3"
)

const (
	StartCmd   = "start"
	BackToMenu = "Back to menu"

	GroupsCmd    = "Groups menu"
	GroupsCreate = "groups_create"
	GroupsStore  = "groups_store"
	GroupsDelete = "groups_delete"
	GroupsList   = "My Groups"
	GroupsShow   = "groups_show/:id"
	GroupsInvite = "groups_invite_link/:id"

	MyReceiptsGroupsCmd = "/my_receipts"
	ReceiptsPayForm     = "receipts_receipts_form/:id"
	ReceiptsPay         = "receipts_pay/:id"

	ReceiptsGroupsCmd = "/receipts_menu"
	ReceiptsCreate    = "receipts_create_form"
	ReceiptsStore     = "receipts_store"
	ReceiptsList      = "receipts_list"
	ReceiptsShow      = "receipts_show"
	ReceiptsEdit      = "receipts_edit"
	ReceiptsUpdate    = "receipts_update"
	ReceiptsDelete    = "receipts_delete"
)

func registerRoutes(router *Router) {
	var userService application.IUserService
	_ = container.Resolve(&userService)
	startCmd := NewStartCommand(userService)
	router.AddRoute(StartCmd, startCmd.Handle, "Start Command")
	router.AddRoute(BackToMenu, startCmd.Handle, "Groups menu Command")

	groupMenuCmd := NewGroupMenuCommand(userService)
	router.AddRoute(GroupsCmd, groupMenuCmd.Handle, "Groups menu Command")

	var groupService application.IGroupService
	_ = container.Resolve(&groupService)
	groupCmd := NewGroupListCommand(userService, groupService)
	router.AddRoute(GroupsList, groupCmd.Handle, "Groups list")

}
