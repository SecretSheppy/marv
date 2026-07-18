package mutant

type Exception struct{}

type ProcessStatus struct {
	ExitStatus int `json:"exitstatus"`
}

type Timeout struct{}

type Value struct {
	Passed bool `json:"passed"`
}

type IsolationResult struct {
	Exception     *Exception     `json:"exception"`
	ProcessStatus *ProcessStatus `json:"process_status"`
	Timeout       *Timeout       `json:"timeout"`
	Value         *Value         `json:"value"`
}

type MutationResult struct {
	IsolationResult        *IsolationResult `json:"isolation_result"`
	MutationSource         string           `json:"mutation_source"`
	MutationIdentification string           `json:"mutation_identification"`
}

type SubjectResult struct {
	CoverageResults []*MutationResult `json:"coverage_results"`
	Identification  string            `json:"identification"`
	Source          string            `json:"source"`
	SourcePath      string            `json:"source_path"`
}

type Results struct {
	SubjectResults []*SubjectResult `json:"subject_results"`
}

// TODO: Diff MutationResult.MutationSource and SubjectResult.Source to get the actual mutation and its lines etc...

// TODO: schema very useful:
//  https://github.com/mbj/mutant/blob/main/docs/session-json-schema.yml

// TODO: mutation_type can be either evil (not test killed it) or neutral (killed by a test) or noop

// TODO: will have to extract replacement lines beginning with + and deleted lines with -
// TODO: (contd) will have to diff the original and replacement to produce descriptions as well as actual replacements

// TODO: operators will have to be defined by marv and then determined based off of this list: (? maybe, this could be very difficult)
//  https://github.com/mbj/mutant/blob/59517844547eef3d67b71a3c736f05bb3c2376da/ruby/lib/mutant/mutation/operators.rb

// TODO: exit_status != 0 || exception == STATUS CRASHED

// TODO: if value:passed == STATUS SURVIVED else STATUS KILLED
