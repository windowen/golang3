package agent

import (
	"context"
	"serverApi/pkg/protobuf/agent"
)

func (s *agentSvr) AgentInfo(ctx context.Context, req *agent.AgentInfoReq) (*agent.AgentInfoResp, error) {
	resp := &agent.AgentInfoResp{
		List: []*agent.Agent{{
			Id: 123,
		}},
		Total: 1,
	}
	return resp, nil
}
