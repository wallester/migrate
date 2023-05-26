package migrator

type Args struct {
	Path           string
	URL            string
	Steps          int
	TimeoutSeconds int
	Up             bool
	NoVerify       bool
}
