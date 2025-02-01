package utils

// converts a prefix and image name like "dining" and "toni_morrison" to "static/dining/toni_morrison.jpg"
func ImageNameToPath(prefix, name string) string {
	fileNameFull := name + ".jpg"
	filePath := "static/" + prefix + "/" + fileNameFull
	return filePath
}
