package backups

// Latest backup.
func (c Client) Latest() (Backup, error) {
	var backup Backup

	list, err := c.List()
	if err != nil {
		return backup, nil
	}

	for _, item := range list {
		if item.Completed > backup.Completed {
			backup = item
		}
	}

	return backup, nil
}
