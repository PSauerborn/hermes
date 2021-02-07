package utils

import (

)

// function used to determine if a slice of strings
// contains a particular item/string
func SliceContains(slice []string, val string) bool {
    for _, value := range(slice) {
        if value == val {
            return true
        }
    }
    return false
}