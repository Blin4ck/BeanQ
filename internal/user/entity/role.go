package entity

const (
	RoleAdmin   = "admin"   //будет свой вход
	RoleManager = "manager" //будет свой вход
	RoleClient  = "client"  //заходит как обычный клиент
)

// Константы
var PredefinedPermissions = []Permission{
	{Operation: "create", Resource: "product"},
	{Operation: "update", Resource: "product"},
	{Operation: "delete", Resource: "product"},
	{Operation: "read", Resource: "product"},
	{Operation: "update", Resource: "order"},
	{Operation: "read", Resource: "order"},
}
