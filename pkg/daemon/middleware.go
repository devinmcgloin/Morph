package daemon

import (
	"github.com/devinmcgloin/fokal/pkg/handler"
	"github.com/devinmcgloin/fokal/pkg/security"
	"github.com/devinmcgloin/fokal/pkg/security/permissions"
	"github.com/justinas/alice"
)

func permission(chain alice.Chain, t permissions.Permission, target uint32) alice.Chain {
	return chain.Append(handler.Middleware{State: &AppState, M: security.Authenticate}.Handler,
		permissions.Middleware{State: &AppState,
			T:          t,
			TargetType: target,
			M:          permissions.PermissionMiddle}.Handler)
}

func auth(chain alice.Chain) alice.Chain {
	return chain.Append(handler.Middleware{State: &AppState, M: security.Authenticate}.Handler)
}
