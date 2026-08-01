package onboarding

import (
	"github.com/hektyc/unraid-mcp-server/internal/config"
	"github.com/hektyc/unraid-mcp-server/internal/mcp"
)

func RegisterTools(s *mcp.Server, cfg *config.Config) {
	s.RegisterTool(&mcp.ToolDef{
		Name:        "onboarding_state",
		Description: "Get first-boot state. Read-only.",
		Query:       `query { customization { onboarding { status isPartnerBuild completed completedAtVersion shouldOpen onboardingState { registrationState isRegistered isFreshInstall hasActivationCode activationRequired } } } isFreshInstall }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "onboarding_complete",
		Description: "Complete onboarding.",
		Query:       `mutation { onboarding { completeOnboarding { status completed shouldOpen } } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "onboarding_reset",
		Description: "Reset onboarding.",
		Query:       `mutation { onboarding { resetOnboarding { status completed shouldOpen } } }`,
	})
}
