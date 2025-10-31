package types

type Pipeline struct {
    Name string `yaml:"name"`
    Jobs map[string]*Job `yaml:"jobs"`
}

type Job struct {
    Name          string `yaml:"name"`
    RunsOn        string `yaml:"runs-on"`
    Steps         []Step `yaml:"steps"`
    Environment   map[string]string `yaml:"env,omitempty"`
}

type Step struct {
    Name    string `yaml:"name"`
    Run     string `yaml:"run,omitempty"`
    Uses    string `yaml:"uses,omitempty"`
}
