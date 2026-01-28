package data_access

type Create struct {
	Project Project `cmd:"project" help:"Create a project"`
}

type Get struct {
	Project Project `cmd:"project" help:"Get a project via it's ID'"`
}

type Update struct {
	Project Project `cmd:"project" help:"Create a project"`
}
