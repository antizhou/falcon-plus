package reader

var (
	// Directory is the directory to store persistence file
	Directory = "/var/www/logs/agent"
)

// Load can load data from persistence directory
func Load(name string) ([]byte, error) {
	path := Directory + "/" + name

	return ReadFile(path)
}

// Save can save data into persistence directory
func Save(name string, data []byte) error {
	exists, _ := DirExists(Directory)

	if !exists {
		err := MkdirAll(Directory, 0755)

		if err != nil {
			return err
		}
	}

	path := Directory + "/" + name
	return WriteFile(path, data, 0644)
}
