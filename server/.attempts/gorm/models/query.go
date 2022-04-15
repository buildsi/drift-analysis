package models

import (
	"gorm.io/gorm/clause"
)

// Get or create an inflection point (but do not update)
func GetOrCreateInflectionPoint(point *InflectionPointRequest) InflectionPoint {

	// Start with the commit and package
	commit := GetOrCreateCommit(point.Commit)
	pkg := GetOrCreatePackage(point.Package)

	var tags []*Tag

	// Create and add each tag (assume unique)
	for _, tag := range point.Tags {
		newTag := GetOrCreateTag(*tag)

		// Don't add an empty tag
		if newTag.Name != "" {
			tags = append(tags, &newTag)
		}
	}

	// Whether or not it exists or not, we return the inflection point
	ip := InflectionPoint{Commit: &commit, Package: &pkg, Tags: tags}
	result := DB.Find(&ip)

	// If we don't find the entry, create it
	if result.RowsAffected == 0 {
		DB.Create(&ip)
	}
	return ip

}

// Get or create a package spec
func GetOrCreateSpec(specRequest *SpecRequest) Spec {

	// Start with the package
	pkg := GetOrCreatePackage(&specRequest.Package)

	// Whether or not it exists or not, we return the inflection point
	spec := Spec{Package: &pkg, Data: specRequest.Data}
	result := DB.Find(&spec)

	// If we don't find the entry, create it
	if result.RowsAffected == 0 {
		DB.Create(&spec)
	}
	return spec

}

// Get or create a build
func GetOrCreateBuild(buildRequest *BuildRequest) Build {

	// Get the spec and inflection point
	var ip InflectionPoint
	var spec Spec
	DB.Where("ID=?", buildRequest.InflectionPointID).First(&ip)
	DB.Where("ID=?", buildRequest.SpecID).First(&spec)

	// Create a new build
	build := Build{Spec: &spec, InflectionPoint: &ip, Status: buildRequest.Status}
	result := DB.Find(&build)

	// If we don't find the entry, create it
	if result.RowsAffected == 0 {
		DB.Create(&build)
	}
	return build

}

// Get or create a commit
func GetOrCreateCommit(commitRequest CommitRequest) Commit {

	// Map the commit request to a commit object
	commit := Commit{
		Digest:    commitRequest.Digest,
		Timestamp: commitRequest.Timestamp,
	}

	// Get or create the commit
	DB.Where("digest=?", commit.Digest).FirstOrCreate(&commit)

	// Whether or not it exists or not, we return the commit
	return commit

}

// Get list of all inflection points
func ListInflectionPoints(pageID int, name string, version string) (points []InflectionPoint) {
	// Default page size 100
	pageSize := 100
	offset := (pageID - 1) * pageSize
	// Initialize points as an empty list
	points = make([]InflectionPoint, 0)

	// Do a query for all points and fill in associated fields
	if len(name) > 0 {
		pkg := Package{}
		DB.Where(&Package{Name: name, Version: version}).First(&pkg)
		if pkg.Name != "" {
			DB.Preload("Tags").Joins("Package").Joins("Commit").Where(
				&InflectionPoint{
					Package: &pkg,
				},
			).Find(&points).Offset(offset).Limit(pageSize)
		}
	} else {
		DB.Preload("Tags").Joins("Package").Joins("Commit").Find(&points).Offset(offset).Limit(pageSize)
	}
	return points
}

// Get list of all commits
func ListCommits(pageID int) *[]Commit {

	// Default page size 100
	pageSize := 100
	offset := (pageID - 1) * pageSize

	commits := []Commit{}
	DB.Find(&commits).Offset(offset).Limit(pageSize)
	return &commits
}

// Get list of all tags
func ListTags(pageID int) *[]Tag {

	pageSize := 100
	offset := (pageID - 1) * pageSize

	tags := []Tag{}

	DB.Find(&tags).Offset(offset).Limit(pageSize)
	return &tags
}

// Get list of all packages
func ListPackages(pageID int) *[]Package {

	// Default page size 100
	pageSize := 100
	offset := (pageID - 1) * pageSize

	packages := []Package{}
	DB.Find(&packages).Offset(offset).Limit(pageSize)
	return &packages
}

// Get list of all specs
func ListSpecs(pageID int) *[]Spec {

	pageSize := 100
	offset := (pageID - 1) * pageSize

	specs := []Spec{}

	DB.Preload("Package").Find(&specs).Offset(offset).Limit(pageSize)
	return &specs
}

// Get list of all builds
func ListBuilds(pageID int) *[]Build {

	pageSize := 100
	offset := (pageID - 1) * pageSize

	builds := []Build{}

	// Includes most nested fields
	DB.Joins("InflectionPoint").Preload("Spec.Package").Preload("InflectionPoint.Commit").
		Preload("InflectionPoint.Tags").Preload("InflectionPoint.Package").
		Preload(clause.Associations).Find(&builds).Offset(offset).Limit(pageSize)
	return &builds
}

// Get or create a package
func GetOrCreatePackage(pkgRequest *PackageRequest) Package {

	// An object to save the result to
	pkg := Package{
		Name:    pkgRequest.Name,
		Version: pkgRequest.Version,
	}

	DB.Where("name=? AND version=?", pkg.Name, pkg.Version).FirstOrCreate(&pkg)

	// Whether or not it exists or not, we return the package
	return pkg

}

// Get or create a tag
func GetOrCreateTag(tagRequest TagRequest) Tag {

	if tagRequest.Name == "" {
		return Tag{}
	}

	// An object to save the result to
	tag := Tag{
		Name: tagRequest.Name,
	}

	DB.Where("name=?", tag.Name).FirstOrCreate(&tag)

	// Whether or not it exists or not, we return the package
	return tag
}
