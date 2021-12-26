package command

type CreateAccountCommand struct {
	firstName string
	lastName  string
}

func NewCreateAccountCommand(firstName string, lastName string) *CreateAccountCommand {
	return &CreateAccountCommand{
		firstName: firstName,
		lastName:  lastName,
	}
}

func (command *CreateAccountCommand) GetFirstName() string {
	return command.firstName
}

func (command *CreateAccountCommand) GetLastName() string {
	return command.lastName
}
