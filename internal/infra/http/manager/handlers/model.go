package handlers

import "github.com/jcleira/encinitas-collector-go/internal/app/manager/aggregates"

// httpProgramRequest represents the request to create a program.
type httpProgramCreateRequest struct {
	ProgramAddress string `json:"program_address"`
}

// ToAggregate converts the httpProgramCreateRequest to an aggregate.Program.
func (hpcr *httpProgramCreateRequest) ToAggregate() aggregates.Program {
	return aggregates.Program{
		ProgramAddress: hpcr.ProgramAddress,
	}
}

// httpProgramGetResponse represents the response to get programs.
type httpProgramGetResponse struct {
	Programs httpPrograms `json:"programs"`
}

// httpProgram represents a program in the HTTP response.
type httpProgram struct {
	ProgramAddress string `json:"program_address"`
	ProgramName    string `json:"program_name"`
}

func httpProgramFromAggregate(
	program aggregates.Program) httpProgram {
	return httpProgram{
		ProgramAddress: program.ProgramAddress,
		ProgramName:    program.ProgramName,
	}
}

type httpPrograms []httpProgram

func httpProgramsFromAggregates(
	programs []aggregates.Program) httpPrograms {
	httpPrograms := make(httpPrograms, len(programs))
	for i, program := range programs {
		httpPrograms[i] = httpProgramFromAggregate(program)
	}

	return httpPrograms
}
