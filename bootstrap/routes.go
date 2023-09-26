package bootstrap

const (
	StartCmd = "/start"

	GroupsCmd    = "/groups_menu"
	GroupsCreate = "groups_create"
	GroupsStore  = "groups_store"
	GroupsDelete = "groups_delete"
	GroupsList   = "groups_list"
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
