package gitremote

import "errors"

var AlreadyExistErr = errors.New("remote already exists")

func (config *Config) Add(r *Remote) error {
	return config.add(r)
}

func (config *Config) AddOrUpdate(r *Remote) error {
	err := config.add(r)
	if err == AlreadyExistErr {
		return config.Update(r)
	}
	return nil
}

func (config *Config) add(r *Remote) error {
	remotes, err := config.List()
	if err != nil {
		return err
	}

	for _, remote := range remotes {
		if remote.Name == r.Name {
			return AlreadyExistErr
		}
	}

	content, err := config.newContent(r)
	if err != nil {
		return err
	}

	err = config.Write(content)
	if err != nil {
		return err
	}

	return nil
}

func (config *Config) newContent(r *Remote) (string, error) {
	content, err := config.Content()
	if err != nil {
		return "", err
	}
	remoteStr, err := r.ToConfig()
	if err != nil {
		return "", err
	}
	return content + "\n" + remoteStr, nil
}
