package replacepackscommand

type ReplacePacksCommand struct {
	Sizes []int `json:"sizes"`
}

func (ReplacePacksCommand) CommandName() string { return "ReplacePacksCommand" }
