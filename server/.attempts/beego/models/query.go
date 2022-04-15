package models

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/beego/beego/v2/client/orm"
)

// Get or create an inflection point (but do not update)
func GetOrCreateInflectionPoint(point InflectionPointRequest) InflectionPoint {

	o := orm.NewOrm()

	// An object to save the result to
	var existing InflectionPoint

	// Create a queryset on the inflection point table
	err := o.QueryTable("inflection_point").Filter("commit__digest", point.Commit.Digest).Filter("package__name", point.Package.Name).Filter("package__version", point.Package.Version).One(&existing)

	// This should not happen
	if err == orm.ErrMultiRows {
		// Have multiple records
		fmt.Printf("Returned Multi Rows Not One")

		// No results, create new!
	} else if err == orm.ErrNoRows {

		commit := GetOrCreateCommit(point.Commit)
		pkg := GetOrCreatePackage(point.Package)

		newPoint := new(InflectionPoint)
		newPoint.Commit = &commit
		newPoint.Package = &pkg
		_, err = o.Update(newPoint)
		if err != nil {
			logs.Warn(err)
		}

		// Create and add each tag (assume unique)
		for _, tag := range point.Tags {
			newTag := GetOrCreateTag(*tag)
			inflectionTag := GetOrCreateInflectionTag(newTag, *newPoint)
			_, err = o.Update(inflectionTag)
			if err != nil {
				logs.Warn(err)
			}
		}

		return *newPoint
	}

	// Otherwise we retrieved it
	return existing

}

// Get or create a commit
func GetOrCreateCommit(commit CommitRequest) Commit {

	o := orm.NewOrm()

	// An object to save the result to
	var existing Commit
	err := o.QueryTable("commit").Filter("digest", commit.Digest).One(&existing)

	// This should not happen
	if err == orm.ErrMultiRows {
		// Have multiple records
		logs.Warn("Returned Multi Rows Not One")

		// No results, create new!
	} else if err == orm.ErrNoRows {

		newCommit := new(Commit)
		newCommit.Digest = commit.Digest
		newCommit.Timestamp = commit.Timestamp
		_, err = o.Update(newCommit)
		if err != nil {
			logs.Warn(err)
		}
		return *newCommit

	}

	// Otherwise we retrieved it
	return existing

}

// Get or create an inflection tag.
// I couldn't get the m2m relation to work between inflection points and tags
func GetOrCreateInflectionTag(tag Tag, point InflectionPoint) InflectionTag {

	o := orm.NewOrm()

	// An object to save the result to
	var existing InflectionTag
	err := o.QueryTable("inflection_tag").Filter("tag__Name", tag.Name).Filter("inflection_point__id", point.Id).One(&existing)

	// This should not happen
	if err == orm.ErrMultiRows {
		// Have multiple records
		logs.Warn("Returned Multi Rows Not One")

		// No results, create new!
	} else if err == orm.ErrNoRows {

		newTag := new(InflectionTag)
		newTag.Tag = &tag
		newTag.InflectionPoint = &point
		_, err = o.Update(newTag)
		if err != nil {
			logs.Warn(err)
		}
		return *newTag

	}

	// Otherwise we retrieved it
	return existing

}

// Get or create a package
func GetOrCreatePackage(pkg *PackageRequest) Package {

	o := orm.NewOrm()

	var existing Package
	err := o.QueryTable("package").Filter("name", pkg.Name).Filter("version", pkg.Version).One(&existing)

	// This should not happen
	if err == orm.ErrMultiRows {
		// Have multiple records
		logs.Warn("Returned Multi Rows Not One")

		// No results, create new!
	} else if err == orm.ErrNoRows {

		newPackage := new(Package)
		newPackage.Name = pkg.Name
		newPackage.Version = pkg.Version
		_, err = o.Update(newPackage)
		if err != nil {
			logs.Warn(err)
		}
		return *newPackage

	}

	// Otherwise we retrieved it
	return existing

}

// Get or create a tag
func GetOrCreateTag(tag TagRequest) Tag {

	o := orm.NewOrm()

	var existing Tag
	err := o.QueryTable("tag").Filter("name", tag.Name).One(&existing)

	// This should not happen
	if err == orm.ErrMultiRows {
		// Have multiple records
		logs.Warn("Returned Multi Rows Not One")

		// No results, create new!
	} else if err == orm.ErrNoRows {

		newTag := new(Tag)
		newTag.Name = tag.Name
		_, err = o.Update(newTag)
		if err != nil {
			logs.Warn(err)
		}
		return *newTag

	}

	// Otherwise we retrieved it
	return existing

}
