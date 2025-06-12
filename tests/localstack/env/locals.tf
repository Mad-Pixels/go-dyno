locals {
  project     = "testing"
  provisioner = "infra"

  tags = {
    "TF"   = "true",
    "Test" = "true"
  }

  # Exclude schemas with "invalid-" prefix
  all_schema_files = fileset("${path.root}/../../fixtures", "*.json")
  schema_files = [
    for file in local.all_schema_files : file
    if !startswith(file, "invalid-")
  ]

  schemas = {
    for file in local.schema_files :
    trimsuffix(file, ".json") => jsondecode(file("${path.root}/../../fixtures/${file}"))
  }
}