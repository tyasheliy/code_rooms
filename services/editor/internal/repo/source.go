package repo

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/tyasheliy/code_rooms/services/editor/internal/entity"
	"os"
	"path/filepath"
)

type FileSourceRepo struct {
	dirPath string
	files   map[uuid.UUID]*fileSourceRepoFile
}

type fileSourceRepoFile struct {
	path   string
	entity *entity.Source
}

func NewFileSourceRepo(dirPath string) *FileSourceRepo {
	return &FileSourceRepo{
		dirPath: dirPath,
		files:   make(map[uuid.UUID]*fileSourceRepoFile),
	}
}

func (r *FileSourceRepo) GetById(ctx context.Context, id uuid.UUID) (*entity.Source, error) {
	file, ok := r.files[id]
	if !ok {
		return nil, errors.New("source not found")
	}

	data, err := os.ReadFile(file.path)
	if err != nil {
		return nil, err
	}

	source := *file.entity
	source.Data = data

	return &source, nil
}

func (r *FileSourceRepo) Create(ctx context.Context, source entity.Source) error {
	exists := false

	for k := range r.files {
		if k == source.Id {
			exists = true
		}
	}

	if exists {
		return errors.New("source exists")
	}

	path, err := r.getPath()
	if err != nil {
		return err
	}

	file, err := os.CreateTemp(path, source.Id.String())
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(source.Data)
	if err != nil {
		return err
	}

	source.Data = make([]byte, 0)

	r.files[source.Id] = &fileSourceRepoFile{
		path:   file.Name(),
		entity: &source,
	}

	return nil
}

func (r *FileSourceRepo) Delete(ctx context.Context, id uuid.UUID) error {
	file, ok := r.files[id]
	if !ok {
		return errors.New("source not found")
	}

	err := os.Remove(file.path)
	if err != nil {
		return err
	}

	delete(r.files, id)

	return nil
}

func (r *FileSourceRepo) Update(ctx context.Context, source entity.Source) error {
	file, ok := r.files[source.Id]
	if !ok {
		return errors.New("source not found")
	}

	f, err := os.OpenFile(file.path, os.O_WRONLY, 0750)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(source.Data)
	if err != nil {
		return err
	}

	file.entity.Session = source.Session
	file.entity.Name = source.Name

	return nil
}

func (r *FileSourceRepo) GetBySession(ctx context.Context, sessionId uuid.UUID) ([]*entity.Source, error) {
	// TODO: implement algorithm
	entities := make([]*entity.Source, 0, len(r.files))

	for _, file := range r.files {
		if file.entity.Session == sessionId {
			source := *file.entity

			data, err := os.ReadFile(file.path)
			if err != nil {
				return nil, err
			}

			source.Data = data

			entities = append(entities, &source)
		}
	}

	return entities, nil
}

func (r *FileSourceRepo) getPath(args ...string) (string, error) {
	path := filepath.Join(args...)
	return filepath.Abs(filepath.Join(r.dirPath, path))
}
