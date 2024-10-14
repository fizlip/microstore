package services

import (
    "io"
    "os"

    log "github.com/sirupsen/logrus"
)

func CreateTempFileTest(content io.Reader, id string) (error) {
    return nil
}

func CreateTempFile(content io.Reader, id string) (error) {
    tempFile, err := os.CreateTemp("/tmp", "docker-image-*.tar")
    log.Println("Create /tmp file with id ", id)
    if err != nil {
        return err
    }
    defer tempFile.Close()

    _, err = io.Copy(tempFile, content)
    log.Println(content)
    if err != nil {
        return err
    }

    return nil
}

func cleanupTempFile(filePath string) {
    //os.Remove(filePath)
}
