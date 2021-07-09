from drift import DriftAnalysis

def test_dep_updated():
    da = DriftAnalysis("abyss")
    # Define "old" spec
    da.last["abyss"] = [{}]
    da.last["abyss"][0]["dependency"] = {}
    da.last["abyss"][0]["dependency"]["full_hash"] = "1"
    # Define "new" spec
    concrete_spec = [{}]
    concrete_spec[0]["dependency"] = {}
    concrete_spec[0]["dependency"]["full_hash"] = "2"

    tags = da.tag_deps(concrete_spec, "abyss", "abyss")
    assert tags == ["dep-updated"]

def test_dep_added():
    da = DriftAnalysis("abyss")
    # Define "old" spec
    da.last["abyss"] = [{}]
    da.last["abyss"][0]["dependency"] = {}
    da.last["abyss"][0]["dependency"]["full_hash"] = "1"
    # Define "new" spec
    concrete_spec = [{}]
    concrete_spec[0]["dependency"] = {}
    concrete_spec[0]["dependency2"] = {}
    concrete_spec[0]["dependency"]["full_hash"] = "1"
    concrete_spec[0]["dependency2"]["full_hash"] = "1"

    tags = da.tag_deps(concrete_spec, "abyss", "abyss")
    assert tags == ["dep-added"]

def test_dep_removed():
    da = DriftAnalysis("abyss")
    # Define "old" spec
    da.last["abyss"] = [{}]
    da.last["abyss"][0]["dependency"] = {}
    da.last["abyss"][0]["dependency2"] = {}
    da.last["abyss"][0]["dependency"]["full_hash"] = "1"
    da.last["abyss"][0]["dependency2"]["full_hash"] = "1"
    # Define "new" spec
    concrete_spec = [{}]
    concrete_spec[0]["dependency"] = {}
    concrete_spec[0]["dependency"]["full_hash"] = "1"

    tags = da.tag_deps(concrete_spec, "abyss", "abyss")
    assert tags == ["dep-deleted"]

def test_patch_modified():
    da = DriftAnalysis("abyss")
    # Define "old" spec
    da.last["abyss"] = [{}]
    da.last["abyss"][0]["abyss"] = {}
    da.last["abyss"][0]["abyss"]["parameters"] = {"patches":["a"]}
    # Define "new" spec
    concrete_spec = [{}]
    concrete_spec[0]["abyss"] = {}
    concrete_spec[0]["abyss"]["parameters"] = {"patches":["b"]}

    tags = da.tag_variants(concrete_spec, "abyss", "abyss")
    assert tags == ["patches-modified"]

def test_patch_added():
    da = DriftAnalysis("abyss")
    # Define "old" spec
    da.last["abyss"] = [{}]
    da.last["abyss"][0]["abyss"] = {}
    da.last["abyss"][0]["abyss"]["parameters"] = {}
    # Define "new" spec
    concrete_spec = [{}]
    concrete_spec[0]["abyss"] = {}
    concrete_spec[0]["abyss"]["parameters"] = {"patches":["b"]}

    tags = da.tag_variants(concrete_spec, "abyss", "abyss")
    assert tags == ["patches-added"]

def test_patch_delete():
    da = DriftAnalysis("abyss")
    # Define "old" spec
    da.last["abyss"] = [{}]
    da.last["abyss"][0]["abyss"] = {}
    da.last["abyss"][0]["abyss"]["parameters"] = {"patches":["a"]}
    # Define "new" spec
    concrete_spec = [{}]
    concrete_spec[0]["abyss"] = {}
    concrete_spec[0]["abyss"]["parameters"] = {}

    tags = da.tag_variants(concrete_spec, "abyss", "abyss")
    assert tags == ["patches-removed"]

def test_variant_added():
    da = DriftAnalysis("abyss")
    # Define "old" spec
    da.last["abyss"] = [{}]
    da.last["abyss"][0]["abyss"] = {}
    da.last["abyss"][0]["abyss"]["parameters"] = {}
    # Define "new" spec
    concrete_spec = [{}]
    concrete_spec[0]["abyss"] = {}
    concrete_spec[0]["abyss"]["parameters"] = {"variant":True}
    
    tags = da.tag_variants(concrete_spec, "abyss", "abyss")
    assert tags == ["variant-added"]

def test_patch_modified():
    da = DriftAnalysis("abyss")
    # Define "old" spec
    da.last["abyss"] = [{}]
    da.last["abyss"][0]["abyss"] = {}
    da.last["abyss"][0]["abyss"]["parameters"] = {"variant":True}
    # Define "new" spec
    concrete_spec = [{}]
    concrete_spec[0]["abyss"] = {}
    concrete_spec[0]["abyss"]["parameters"] = {"variant":False}

    tags = da.tag_variants(concrete_spec, "abyss", "abyss")
    assert tags == ["variant-modified"]

def test_variant_deleted():
    da = DriftAnalysis("abyss")
    # Define "old" spec
    da.last["abyss"] = [{}]
    da.last["abyss"][0]["abyss"] = {}
    da.last["abyss"][0]["abyss"]["parameters"] = {"variant":True}
    # Define "new" spec
    concrete_spec = [{}]
    concrete_spec[0]["abyss"] = {}
    concrete_spec[0]["abyss"]["parameters"] = {}

    tags = da.tag_variants(concrete_spec, "abyss", "abyss")
    assert tags == ["variant-removed"]