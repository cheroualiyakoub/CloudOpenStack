
package terraform_utilis


import (
    "io"
    "os"
    "path/filepath"
)


// CopyFile copies a file from src to dst. If dst does not exist, it will be created.
func CopyFile(src, dst string) error {
    // Open the source file
    in, err := os.Open(src)
    if err != nil {
        return err
    }
    defer in.Close()

    // Ensure the destination directory exists
    if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
        return err
    }

    // Create the destination file
    out, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer out.Close()

    // Copy the contents
    _, err = io.Copy(out, in)
    if err != nil {
        return err
    }

    // Flush to disk
    return out.Sync()
}
