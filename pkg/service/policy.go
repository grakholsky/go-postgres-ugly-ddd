package service

import "github.com/casbin/casbin"

type (
	PolicyConfig struct {
		CasbinModelPath  string
		CasbinPolicyPath string
	}

	Policy struct {
		config         *PolicyConfig
		casbinEnforcer *casbin.Enforcer
	}
)

func NewPolicy(config *PolicyConfig) *Policy {
	p := new(Policy)
	p.casbinEnforcer = casbin.NewEnforcer(config.CasbinModelPath, config.CasbinPolicyPath)
	return p
}

func (s *Policy) Check(role, url, method string) bool {
	return s.casbinEnforcer.Enforce(role, url, method)
}
