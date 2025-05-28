locals {
  project     = "applingo"
  provisioner = "infra"

  tags = {
    "TF"   = "true",
    "Test" = "true"
  }

  schema_files = fileset("${path.root}/../../../examples", "*.json")
  schemas = {
    for file in local.schema_files :
    trimsuffix(file, ".json") => jsondecode(file("${path.root}/../../../examples/${file}"))
  }
}