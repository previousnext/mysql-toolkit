package backups

// Backup performed by the Acquia platform.
type Backup struct {
	ID        int64  `json:"id,string"`
	Completed int64  `json:"completed,string"`
	Link      string `json:"link"`
}
