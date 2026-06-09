package utils

import (
	"errors"
	"go.mongodb.org/mongo-driver/v2/mongo"
)


func IsIndexAlreadyExistsError(err error) bool {
    // MongoDB returns code 85 (IndexOptionsConflict) or 86 (IndexKeySpecsConflict)
    // when the index already exists with the same or different options.
    var cmdErr mongo.CommandError
    if errors.As(err, &cmdErr) {
        return cmdErr.Code == 85 || cmdErr.Code == 86
    }
    return false
}