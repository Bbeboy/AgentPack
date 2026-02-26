package fsutil

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func CopyDirContents(src, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return fmt.Errorf("no se pudo leer el directorio origen '%s': %w", src, err)
	}

	for _, entry := range entries {
		sourcePath := filepath.Join(src, entry.Name())
		targetPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := CopyDir(sourcePath, targetPath); err != nil {
				return err
			}
			continue
		}

		if err := copyFile(sourcePath, targetPath); err != nil {
			return err
		}
	}

	return nil
}

func CopyDir(src, dst string) error {
	info, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("no se pudo leer el directorio '%s': %w", src, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("'%s' no es un directorio", src)
	}

	if err := os.MkdirAll(dst, info.Mode().Perm()); err != nil {
		return fmt.Errorf("no se pudo crear el directorio '%s': %w", dst, err)
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return fmt.Errorf("no se pudo leer el contenido de '%s': %w", src, err)
	}

	for _, entry := range entries {
		sourcePath := filepath.Join(src, entry.Name())
		targetPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := CopyDir(sourcePath, targetPath); err != nil {
				return err
			}
			continue
		}

		if err := copyFile(sourcePath, targetPath); err != nil {
			return err
		}
	}

	return nil
}

func MergeDir(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("no se pudo leer el origen '%s': %w", src, err)
	}
	if !srcInfo.IsDir() {
		return fmt.Errorf("'%s' no es un directorio", src)
	}

	if err := os.MkdirAll(dst, 0o755); err != nil {
		return fmt.Errorf("no se pudo crear el destino '%s': %w", dst, err)
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return fmt.Errorf("no se pudo leer el contenido de '%s': %w", src, err)
	}

	for _, entry := range entries {
		sourcePath := filepath.Join(src, entry.Name())
		targetPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			targetInfo, err := os.Stat(targetPath)
			if err != nil {
				if os.IsNotExist(err) {
					if err := CopyDir(sourcePath, targetPath); err != nil {
						return err
					}
					continue
				}
				return fmt.Errorf("no se pudo revisar '%s': %w", targetPath, err)
			}

			if !targetInfo.IsDir() {
				if err := os.Remove(targetPath); err != nil {
					return fmt.Errorf("no se pudo reemplazar archivo por directorio en '%s': %w", targetPath, err)
				}
				if err := CopyDir(sourcePath, targetPath); err != nil {
					return err
				}
				continue
			}

			if err := MergeDir(sourcePath, targetPath); err != nil {
				return err
			}
			continue
		}

		if err := copyFile(sourcePath, targetPath); err != nil {
			return err
		}
	}

	return nil
}

func copyFile(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("no se pudo leer el archivo '%s': %w", src, err)
	}

	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return fmt.Errorf("no se pudo crear directorio de destino para '%s': %w", dst, err)
	}

	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("no se pudo abrir '%s': %w", src, err)
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, srcInfo.Mode().Perm())
	if err != nil {
		return fmt.Errorf("no se pudo crear '%s': %w", dst, err)
	}

	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return fmt.Errorf("no se pudo copiar '%s' en '%s': %w", src, dst, err)
	}

	if err := out.Close(); err != nil {
		return fmt.Errorf("no se pudo cerrar '%s': %w", dst, err)
	}

	return nil
}
