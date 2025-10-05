package requests

import (
	"github.com/mauriciorobertodev/whappy-go/internal/app/input"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/instance"
	"github.com/mauriciorobertodev/whappy-go/internal/presentation/http"
)

type CreateInstanceRequest struct {
	Name *string `json:"name"`
}

func (r *CreateInstanceRequest) Validate() *http.ErrorBag {
	var bag = http.NewErrorBag()

	if r.Name != nil && len(*r.Name) > instance.MaxNameLength {
		bag.Add("name", "name is too long")
	}

	return bag
}

func (r *CreateInstanceRequest) ToInput() input.CreateInstance {
	name := ""
	if r.Name != nil {
		name = *r.Name
	}

	return input.CreateInstance{
		Name: name,
	}
}

type JoinGroupRequest struct {
	Invite string `json:"invite"`
}
