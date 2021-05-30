package git_cmd

type GitApi struct {
	Dir string
	*option
}

type ItemBase struct {
	Dir string
}
